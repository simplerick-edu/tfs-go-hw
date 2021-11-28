package domain

// Result
const (
	Success = "success"
	Error   = "error"
)

// Status
const (
	Placed    = "placed"
	Cancelled = "cancelled"
)

type BaseResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
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
	Type   string `json:"type"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
	Order  Order  `json:"orderPriorExecution"`
}

type Order struct {
	OrderID    string  `json:"orderId"`
	Symbol     string  `json:"symbol"`
	Side       string  `json:"side"`
	Type       string  `json:"type"`
	LimitPrice float64 `json:"limit_price"`
	Quantity   float64 `json:"quantity"`
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
