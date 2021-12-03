package domain

type Result string

// Result
const (
	Success Result = "success"
	Error   Result = "error"
)

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

type OpenPositionsResponse struct {
	BaseResponse
	OpenPositions []Position `json:"openPositions"`
}

type Status struct {
	Status      string       `json:"status"`
	OrderEvents []OrderEvent `json:"orderEvents"`
}
