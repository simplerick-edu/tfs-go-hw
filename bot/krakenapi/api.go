package krakenapi

import (
	"net/http"
	"time"
	// not-std
	"github.com/gorilla/websocket"
)

const (
	WebSocketURL     = "wss://demo-futures.kraken.com/ws/v1"
	SendOrderURL     = "https://demo-futures.kraken.com/derivatives/api/v3/sendorder?"
	CancelOrdersURL  = "https://demo-futures.kraken.com/derivatives/api/v3/cancelallorders"
	OpenPositionsURL = "https://demo-futures.kraken.com/derivatives/api/v3/openpositions"
)
const DefaultHttpTimeout = 10 * time.Second

type Config struct {
	PublicKey   string `yaml:"public_key"`
	PrivateKey  string `yaml:"private_key"`
	HttpTimeout string `yaml:"http_timeout"`
}

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

func NewWithConfig(config Config) *KrakenAPI {
	timeout, err := time.ParseDuration(config.HttpTimeout)
	if err != nil {
		timeout = DefaultHttpTimeout
	}
	return New(config.PublicKey, config.PrivateKey, timeout)
}
