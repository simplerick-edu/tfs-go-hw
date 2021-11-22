package krakenapi

import (
	"bot/domain"
	"fmt"
	"github.com/gorilla/websocket"
)

const ReconnectAttempts = 10

func (k *KrakenAPI) Connect(trialNum int) error {
	if c, _, err := websocket.DefaultDialer.Dial(WebSocketURL, nil); err == nil {
		k.conn = c
		k.closed = false
		return nil
	}
	return k.Connect(trialNum - 1)
}

func (k *KrakenAPI) Subscribe(productIDs ...string) (<-chan string, error) {
	out := make(chan string)
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
		Feed:       "ticker",
		ProductIDs: productIDs,
	}
	if err := ConnectAndSendMsg(msg); err != nil {
		return out, nil
	}
	go func() {
		defer close(out)
		for {
			_, message, err := k.conn.ReadMessage()
			if err != nil {
				if k.closed {
					fmt.Println("Normal Termination")
					return
				} else {
					fmt.Println("Retry")
					if err := ConnectAndSendMsg(msg); err != nil {
						return
					}
				}
			}
			out <- string(message)
		}
	}()
	return out, nil
}

func (k *KrakenAPI) Unsubscribe() error {
	k.closed = true
	err := k.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}
