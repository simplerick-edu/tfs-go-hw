package domain

import "time"

type OrderType string

// Order Types
const (
	LmtType OrderType = "lmt"
	IocType OrderType = "ioc"
	MktType OrderType = "mkt"
)

type Action string

// Actions (sides)
const (
	Buy  Action = "buy"
	Sell Action = "sell"
	None Action = "none"
)

type OrderEvent struct {
	Type      string  `json:"type"`
	Price     float64 `json:"price,omitempty"`  // execution price
	Amount    int64   `json:"amount,omitempty"` // execution quantity
	ExecOrder *Order  `json:"orderPriorExecution,omitempty"`
	Order     *Order  `json:"order,omitempty"`
}

type Order struct {
	OrderID    string    `json:"orderId"`
	Symbol     string    `json:"symbol"`
	Side       Action    `json:"side"`
	Type       OrderType `json:"type"`
	LimitPrice float64   `json:"limitPrice"`
	Quantity   float64   `json:"quantity"`
	TS         time.Time `json:"timestamp"`
}

func NewOrder(symbol string, side Action, orderType OrderType, price float64, quantity int64) *Order {
	return &Order{
		Symbol:     symbol,
		Side:       side,
		Type:       orderType,
		LimitPrice: price,
		Quantity:   float64(quantity),
		TS:         time.Now(),
	}
}
