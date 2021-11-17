package krakenapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	WebSocketURL    string = "wss://futures.kraken.com/ws/v1"
	SendOrderURL    string = "/api/v3/sendorder"
	CancelOrdersURL string = "/api/v3/cancelallorders"
)

type Message struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIDs []string `json:"product_ids"`
}

//type ExhangeAPI interface {
//	Connect(url string) error
//	Subscribe(instruments []string) (<-chan string, error)
//	Close() error
//}

func NewKrakenService(publicKey string, privateKey string, timeout time.Duration) *KrakenAPI {
	return &KrakenAPI{
		publicKey,
		privateKey,
		&http.Client{Timeout: timeout},
		nil,
	}
}

type KrakenAPI struct {
	publicKey  string
	privateKey string
	client     *http.Client
	conn       *websocket.Conn
}

// Websocket

func (k *KrakenAPI) Connect(ctx context.Context, url string) error {
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	k.conn = c
	return err
}

func (k *KrakenAPI) Subscribe(productIDs []string) (<-chan string, error) {
	out := make(chan string)
	msg := Message{
		Event:      "subscribe",
		Feed:       "ticker",
		ProductIDs: productIDs,
	}
	go func() {
		for {
			_, message, err := k.conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			out <- string(message)
		}
	}()
	err := k.conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("write:", err)
	}
	return out, err
}

func (k *KrakenAPI) Close() error {
	err := k.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}

// REST

//func (k *KrakenAPI) SendOrder() error {
//	req, err := http.NewRequest(http.MethodPost, SendOrderURL+"?", nil)
//	resp, err := k.client.Do()
//}

//func (k *KrakenAPI) CancelOrder() error {
//	req, err := http.NewRequest(http.MethodPost, CancelOrdersURL+"?", nil)
//	//resp, err := k.client.Do()
//}

func (k *KrakenAPI) createRequest(reqURL string, data url.Values) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, reqURL, strings.NewReader(data.Encode()))
}

func (k *KrakenAPI) privateRequest(req *http.Request) (*http.Request, error) {
	// Get sign
	body, _ := req.GetBody()
	postData, _ := io.ReadAll(body)
	endpointPath := []byte(req.URL.Path)
	sign := generateSign(postData, endpointPath, k.privateKey)
	// Add headers
	req.Header.Add("APIKey", k.publicKey)
	req.Header.Add("Authent", sign)
	return req, nil
}

//func (k *KrakenAPI) sendRequest(req *http.Request) (interface{}, error) {
//	resp, err := k.client.Do(req)
//	return resp, err
//}
//

func generateSign(postData []byte, endpointPath []byte, privateKey string) string {
	// Concatenate postData + nonce + endpointPath
	src := append(postData, endpointPath...)
	//Hash the result of step 1 with the SHA-256 algorithm
	sha := sha256.New()
	sha.Write(src)
	//Base64-decode your api_secret
	apiSecret, _ := base64.StdEncoding.DecodeString(privateKey)
	//Use the result of step 3 to hash the result of the step 2 with the HMAC-SHA-512 algorithm
	h := hmac.New(sha512.New, apiSecret)
	h.Write(sha.Sum(nil))
	//Base64-encode the result of step 4
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}
