package exchange


import (
    "github.com/gorilla/websocket"
    "fmt"
    // "context"
)


type Message struct {
    Event string `json:"event"`
    Feed string `json:"feed"`
    ProductIDs []string `json:"product_ids"`
}





type ExhangeService interface {
    Connect(url string) (error)
    Subscribe(productIDs []string) (<-chan string, error)
    Close() (error)
}


type KrakenService struct {
    conn *websocket.Conn
}


func (k *KrakenService) Connect(url string) (error) {
    c, _, err := websocket.DefaultDialer.Dial(url, nil)
    k.conn = c
	return err
}



func (k *KrakenService) Subscribe(productIDs []string) (<-chan string, error) {
    out := make(chan string)
    msg := Message{
        Event : "subscribe",
        Feed : "ticker",
        ProductIDs : productIDs,
    }
    go func() {
		for {
			_, message, err := k.conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
			}
			out <-string(message)
		}
	}()
    err := k.conn.WriteJSON(msg)
    if err != nil {
        fmt.Println("write:", err)
    }
    return out, err
}


func (k *KrakenService) Close() (error) {
    err := k.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
    return err
}
