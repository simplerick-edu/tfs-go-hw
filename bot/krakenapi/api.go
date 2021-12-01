package krakenapi

import (
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"time"
)

const (
	WebSocketURL     = "wss://demo-futures.kraken.com/ws/v1"
	SendOrderURL     = "https://demo-futures.kraken.com/derivatives/api/v3/sendorder?"
	CancelOrdersURL  = "https://demo-futures.kraken.com/derivatives/api/v3/cancelallorders"
	OpenPositionsURL = "https://demo-futures.kraken.com/derivatives//api/v3/openpositions"
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

func NewWithConfig(filePath string) (*KrakenAPI, error) {
	var c map[string]string
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	timeout, _ := time.ParseDuration(c["timeout"])
	return New(c["public_key"], c["private_key"], timeout), nil
}
