package service

import (
	"bot/domain"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"text/template"
)

type ExchangeAPI interface {
	GetPositions() (*domain.OpenPositionsResponse, error)
	SendOrder(symbol string, side domain.Action, orderType domain.OrderType, price float64, size int64) (*domain.SendOrderResponse, error)
	CancelOrders() (*domain.CancelOrdersResponse, error)
	Subscribe(...string) (<-chan domain.Ticker, error)
	Unsubscribe() error
}

type Predictor interface {
	Predict(tickers ...domain.Ticker) (float64, error)
}

type Storage interface {
	StoreEvent(ctx context.Context, event domain.OrderEvent) error
}

type Notifier interface {
	Start() error
	Notify(text string)
	Stop()
}

var NotificationTemplate = template.Must(template.ParseFiles("service/notification.tmpl"))

type Bot struct {
	// api
	exchangeAPI ExchangeAPI
	notifier    Notifier
	storage     Storage
	model       Predictor
	// parameters
	instrument        string
	maxPositionSize   int64
	orderSize         int64
	decisionThreshold float64
	sequenceLength    int
	priceSlipPercent  int64
	// internal variables
	mu              sync.Mutex
	openPositions   map[string]int64
	shutdownChannel chan interface{}
}

func New(exchangeAPI ExchangeAPI,
	notifier Notifier,
	storage Storage,
	model Predictor,
	instrument string,
	maxPositionSize int64,
	orderSize int64,
	decisionThreshold float64,
	sequenceLength int,
	priceSlipPercent int64) *Bot {
	return &Bot{
		exchangeAPI,
		notifier,
		storage,
		model,
		instrument,
		maxPositionSize,
		orderSize,
		decisionThreshold,
		sequenceLength,
		priceSlipPercent,
		sync.Mutex{},
		make(map[string]int64),
		make(chan interface{}),
	}
}

func (b *Bot) Start() error {
	if err := b.FetchOpenPositions(); err != nil {
		return err
	}
	tickers, err := b.exchangeAPI.Subscribe(b.instrument)
	if err != nil {
		return err
	}
	if err := b.notifier.Start(); err != nil {
		return err
	}
	// collect tickers
	tickerSequences := make(chan []domain.Ticker)
	go func() {
		defer func() {
			log.Println("shutdown")
			close(tickerSequences)
			b.exchangeAPI.Unsubscribe()
			b.notifier.Stop()
		}()
		seq := make([]domain.Ticker, 0, b.sequenceLength)
		for ticker := range tickers {
			select {
			case <-b.shutdownChannel:
				return
			default:
				if len(seq) == b.sequenceLength {
					tickerSequences <- seq
					seq = make([]domain.Ticker, 0, b.sequenceLength)
				}
				seq = append(seq, ticker)
			}
		}
	}()
	// process sequences
	go func() {
		for tickers := range tickerSequences {
			action, err := b.makeDecision(tickers)
			if err != nil {
				log.Fatal(err)
				return
			}
			// calculate price
			var price float64
			if action == domain.Buy {
				price = tickers[len(tickers)-1].Ask * (1 + float64(b.priceSlipPercent)/100)
			}
			if action == domain.Sell {
				price = tickers[len(tickers)-1].Bid * (1 - float64(b.priceSlipPercent)/100)
			}
			err = b.ChangePosition(action, b.orderSize, price)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}()
	return err
}

func (b *Bot) FetchOpenPositions() error {
	resp, err := b.exchangeAPI.GetPositions()
	if err != nil {
		return fmt.Errorf("failed to fetch open positions: %w", err)
	}
	if resp.Result != domain.Success {
		return errors.New(*resp.Error)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.openPositions = make(map[string]int64)
	for _, pos := range resp.OpenPositions {
		if pos.Side == "long" {
			b.openPositions[pos.Symbol] = pos.Size
		} else {
			b.openPositions[pos.Symbol] = -pos.Size
		}
	}
	log.Println("current open positions: ", b.openPositions)
	return nil
}

func (b *Bot) ChangePosition(side domain.Action, size int64, price float64) error {
	var sign int64
	switch side {
	case domain.None:
		log.Println("action was not specified, position unchanged")
		return nil
	case domain.Buy:
		sign = 1
	case domain.Sell:
		sign = -1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	currentPos := b.openPositions[b.instrument]
	// keep position size within limits (-maxPositionSize, +maxPositionSize)
	size = domain.Min(size, b.maxPositionSize-sign*currentPos)
	if size != 0 {
		resp, err := b.exchangeAPI.SendOrder(b.instrument, side, domain.IocType, price, size)
		if err != nil {
			return err
		}
		actualAmount := b.processResponse(resp)
		b.openPositions[b.instrument] += sign * actualAmount
		if b.openPositions[b.instrument] == 0 {
			delete(b.openPositions, b.instrument)
		}
	}
	return nil
}

func (b *Bot) processResponse(resp *domain.SendOrderResponse) int64 {
	var message string
	var amount int64
	if resp.Result == domain.Error {
		message = "sending order failed: " + *resp.Error
	} else {
		var buff bytes.Buffer
		NotificationTemplate.Execute(&buff, resp.SendStatus)
		message = buff.String()
		// get filled amount, store to db
		if resp.SendStatus.Status == "placed" {
			event := resp.SendStatus.OrderEvents[0]
			err := b.storage.StoreEvent(context.Background(), event)
			if err != nil {
				log.Println(err)
			}
			amount = event.Amount
		}
	}
	// send notification
	b.notifier.Notify(message)
	log.Println(message)
	return amount
}

func (b *Bot) makeDecision(tickers []domain.Ticker) (domain.Action, error) {
	// receive predicted value in range (0,1)
	value, err := b.model.Predict(tickers...)
	if err != nil {
		return domain.None, err
	}
	if value > b.decisionThreshold {
		return domain.Buy, nil
	}
	if value < 1-b.decisionThreshold {
		return domain.Sell, nil
	}
	return domain.None, nil
}

func (b *Bot) Stop() {
	close(b.shutdownChannel)
}
