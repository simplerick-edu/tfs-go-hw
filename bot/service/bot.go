package service

import (
	"bot/domain"
	"bytes"
	"sync"
	"text/template"
)

type ExchangeAPI interface {
	SendOrder(symbol string, side domain.Action, orderType domain.OrderType, price float64, size int64) (*domain.SendOrderResponse, error)
	CancelOrders() (*domain.CancelOrdersResponse, error)
	Subscribe(...string) (<-chan domain.Ticker, error)
	Unsubscribe() error
}

type Predictor interface {
	Predict(domain.Ticker) (float64, error)
}

type Database interface {
	Store(domain.OrderEvent) error
}

type Notifier interface {
	Start() error
	Notify(text string)
	Stop()
}

var NotificationTemplate = template.Must(template.ParseFiles("notification.tmpl"))

type Bot struct {
	// api
	exchangeAPI ExchangeAPI
	notifier    Notifier
	database    Database
	model       Predictor
	// parameters
	instrument        string
	maxPositionSize   int64
	OrderSize         int64
	decisionThreshold float64
	priceSlipPercent  int64
	// internal variables
	mu            sync.Mutex
	openPositions map[string]int64
}

//func (b *Bot) Start() error {
//	tickers, err := b.exchangeAPI.Subscribe(b.instrument)
//	go func() {
//		for ticker := range(tickers) {
//			value, err := b.calculateIndicator(ticker)
//			price := ticker.Ask * (1 + float64(b.priceSlipPercent)/100)
//			price := ticker.Bid * (1 - float64(b.priceSlipPercent)/100)
//			if err != nil {
//				return
//			}
//			err = b.executeTrade(value)
//		}
//	}()
//	return err
//}

func (b *Bot) GetOpenPositions() {

}

func (b *Bot) ChangePosition(side domain.Action, size int64, price float64) error {
	var sign int64
	switch side {
	case domain.None:
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
		resp, err := b.exchangeAPI.SendOrder(b.instrument, side, domain.LmtType, price, size)
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
		message = *resp.Error
	} else {
		var buff bytes.Buffer
		NotificationTemplate.Execute(&buff, resp.SendStatus)
		message = buff.String()
		// get filled amount, store to db
		if resp.SendStatus.Status == "placed" {
			event := resp.SendStatus.OrderEvents[0]
			b.database.Store(event)
			amount = event.Amount
		}
	}
	// send notification
	b.notifier.Notify(message)
	return amount
}

func (b *Bot) makeDecision(ticker domain.Ticker) (domain.Action, error) {
	// receive predicted value in range (0,1)
	value, err := b.model.Predict(ticker)
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

//func (b *Bot) ClosePositions() error {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	for symbol, position := range b.openPositions {
//		if position > 0
//	}
//}
