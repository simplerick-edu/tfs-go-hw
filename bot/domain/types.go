package domain

import "time"

type BaseResponse struct {
	Result Result  `json:"result"`
	Error  *string `json:"error,omitempty"`
}

type CancelOrdersResponse struct {
	BaseResponse
	CancelStatus Status `json:"cancelStatus"`
}

type SendOrderResponse struct {
	BaseResponse
	SendStatus Status `json:"sendStatus"`
}

type Status struct {
	Status      string       `json:"status"`
	OrderEvents []OrderEvent `json:"orderEvents"`
}

type OrderEvent struct {
	Type      string  `json:"type"`
	Price     float64 `json:"price,omitempty"`
	Amount    int64   `json:"amount,omitempty"`
	ExecOrder *Order  `json:"orderPriorExecution,omitempty"`
	Order     *Order  `json:"order,omitempty"`
}

type Order struct {
	OrderID    string    `json:"orderId"`
	Symbol     string    `json:"symbol"`
	Side       Action    `json:"side"`
	Type       string    `json:"type"`
	LimitPrice float64   `json:"limitPrice"`
	Quantity   float64   `json:"quantity"`
	TS         time.Time `json:"timestamp"`
}

type OpenPositionsResponse struct {
	BaseResponse
	OpenPositions []Position `json:"openPositions"`
}

type Position struct {
	Side   string  `json:"side"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Size   int64   `json:"size"`
}

type Message struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIDs []string `json:"product_ids"`
}

type Ticker struct {
	Time         int64   `json:"time"`
	Feed         string  `json:"feed"`
	ProductId    string  `json:"product_id"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	BidSize      int     `json:"bid_size"`
	AskSize      int     `json:"ask_size"` //
	Volume       int     `json:"volume"`   // volume for 24 hours
	Dtm          int     `json:"dtm"`      // the days until maturity
	Leverage     string  `json:"leverage"` // leverage
	Last         int     `json:"last"`     // last trade price
	Change       float64 `json:"change"`   // 24h change
	OpenInterest int     `json:"openInterest"`
}
