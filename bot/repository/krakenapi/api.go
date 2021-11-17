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
	"strconv"
	"strings"
	"time"
)

type OpSide string

const (
	Buy  OpSide = "buy"
	Sell OpSide = "sell"

	WebSocketURL    string = "wss://futures.kraken.com/ws/v1"
	SendOrderURL    string = "https://demo-futures.kraken.com/derivatives/api/v3/sendorder"
	CancelOrdersURL string = "https://demo-futures.kraken.com/derivatives/api/v3/cancelallorders"
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

func New(publicKey string, privateKey string, timeout time.Duration) *KrakenAPI {
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

// REST

func (k *KrakenAPI) SendOrder(pair string, size int, side OpSide) (*http.Response, error) {
	data := map[string]interface{}{
		"pair":      pair,
		"size":      size,
		"side":      side,
		"ordertype": "mkt",
	}
	req, err := k.createRequest(SendOrderURL, data)
	if err != nil {
		return nil, err
	}
	req = k.privateRequest(req)
	return k.sendRequest(req)
}

func (k *KrakenAPI) CancelOrders() (*http.Response, error) {
	req, err := k.createRequest(CancelOrdersURL, nil)
	if err != nil {
		return nil, err
	}
	req = k.privateRequest(req)
	return k.sendRequest(req)
}

func (k *KrakenAPI) createRequest(reqURL string, data map[string]interface{}) (*http.Request, error) {
	values := url.Values{}
	for key, value := range data {
		switch v := value.(type) {
		case string:
			values.Set(key, v)
		case int64:
			values.Set(key, strconv.FormatInt(v, 10))
		case float64:
			values.Set(key, strconv.FormatFloat(v, 'f', 8, 64))
		case bool:
			values.Set(key, strconv.FormatBool(v))
		default:
			return nil, fmt.Errorf("unknown value type %v for key %s", value, key)
		}
	}
	return http.NewRequest(http.MethodPost, reqURL, strings.NewReader(values.Encode()))
}

func (k *KrakenAPI) privateRequest(req *http.Request) *http.Request {
	// Get sign
	body, _ := req.GetBody()
	postData, _ := io.ReadAll(body)
	endpointPath := []byte(req.URL.Path[12:])
	sign := generateSign(postData, endpointPath, k.privateKey)
	// Add headers
	req.Header.Add("APIKey", k.publicKey)
	req.Header.Add("Authent", sign)
	return req
}

func (k *KrakenAPI) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := k.client.Do(req)
	return resp, err
}

func (k *KrakenAPI) Close() error {
	err := k.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}

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
