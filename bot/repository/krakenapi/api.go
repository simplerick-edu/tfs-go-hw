package krakenapi

import (
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"time"
)

const (
	Buy             = "buy"
	Sell            = "sell"
	WebSocketURL    = "wss://demo-futures.kraken.com/ws/v1"
	SendOrderURL    = "https://demo-futures.kraken.com/derivatives/api/v3/sendorder?"
	CancelOrdersURL = "https://demo-futures.kraken.com/derivatives/api/v3/cancelallorders"
	BalanceURL      = "https://demo-futures.kraken.com/derivatives/api/v3/balance"
)

type KrakenAPI struct {
	publicKey  string
	privateKey string
	client     *http.Client
	conn       *websocket.Conn
	closed     bool
}

func New(publicKey string, privateKey string, timeout time.Duration) *KrakenAPI {
	return &KrakenAPI{
		publicKey,
		privateKey,
		&http.Client{Timeout: timeout},
		nil,
		false,
	}
}

func NewFromConfig(filePath string) *KrakenAPI {
	var c map[string]string
	data, _ := os.ReadFile(filePath)
	_ = yaml.Unmarshal(data, &c)
	timeout, _ := time.ParseDuration(c["timeout"])
	return New(c["public_key"], c["private_key"], timeout)
}
