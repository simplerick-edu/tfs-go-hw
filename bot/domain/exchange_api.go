package domain

type ExchangeAPI interface {
	Buy(symbol string, price float64, size int64) (*SendOrderResponse, error)
	Sell(symbol string, price float64, size int64) (*SendOrderResponse, error)
	CancelOrders() (*CancelOrdersResponse, error)
	Subscribe(...string) (<-chan string, error)
	Unsubscribe() error
}
