package krakenapi

import (
	"encoding/json"
	"errors"
	"fmt"
	// not-std
	"bot/domain"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const ReconnectAttempts = 10

var ErrMaxReconnects = errors.New("the maximum number of reconnection attempts has been reached")

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
			return fmt.Errorf("can't write to websocket: %w", err)
		}
		return nil
	}
	msg := domain.Message{
		Event:      "subscribe",
		Feed:       "ticker",
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
					log.Info("normal termination")
					return
				} else {
					log.Warningf("websocket: retry to connect")
					if err := ConnectAndSendMsg(msg); err != nil {
						log.Error(err)
						return
					}
				}
			}
			ticker := domain.Ticker{}
			err = json.Unmarshal(message, &ticker)
			if err != nil {
				log.Info("websocket: ", err)
				return
			}
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
