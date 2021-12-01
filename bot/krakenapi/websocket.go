package krakenapi

import (
	"bot/domain"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

const ReconnectAttempts = 10

var (
	ErrMaxReconnects = errors.New("the maximum number of reconnection attempts has been reached")
)

func (k *KrakenAPI) Connect(trialNum int) error {
	if trialNum < 0 {
		return ErrMaxReconnects
	}
	if c, _, err := websocket.DefaultDialer.Dial(WebSocketURL, nil); err == nil {
		k.conn = c
		k.closed = false
		return nil
	}
	return k.Connect(trialNum - 1)
}

func (k *KrakenAPI) Subscribe(productIDs ...string) (<-chan domain.Ticker, error) {
	out := make(chan domain.Ticker)
	ConnectAndSendMsg := func(msg domain.Message) error {
		if err := k.Connect(ReconnectAttempts); err != nil {
			return err
		}
		if err := k.conn.WriteJSON(msg); err != nil {
			return err
		}
		return nil
	}
	msg := domain.Message{
		Event:      "subscribe",
		Feed:       "candles_trade_1m",
		ProductIDs: productIDs,
	}
	if err := ConnectAndSendMsg(msg); err != nil {
		return out, err
	}
	go func() {
		defer close(out)
		for {
			_, message, err := k.conn.ReadMessage()
			if err != nil {
				if k.closed {
					log.Println("Normal Termination")
					return
				} else {
					log.Println("Websocket: retry to connect")
					if err := ConnectAndSendMsg(msg); err != nil {
						return
					}
				}
			}
			ticker := domain.Ticker{}
			json.Unmarshal(message, &ticker)
			out <- ticker
		}
	}()
	return out, nil
}

func (k *KrakenAPI) Unsubscribe() error {
	k.closed = true
	err := k.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}
