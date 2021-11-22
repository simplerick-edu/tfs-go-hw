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

type OrderEvent struct {
	Type   string `json:"type"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type Status struct {
	Status      string       `json:"status"`
	OrderEvents []OrderEvent `json:"orderEvents"`
}

type CancelOrdersResponse struct {
	BaseResponse
	CancelStatus Status `json:"cancelStatus"`
}

type SendOrderResponse struct {
	BaseResponse
	SendStatus Status `json:"sendStatus"`
}
