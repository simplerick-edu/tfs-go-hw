package domain

type Ticker struct {
	Time         int64   `json:"time"`
	Feed         string  `json:"feed"`
	ProductId    string  `json:"product_id"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	BidSize      float64 `json:"bid_size"`
	AskSize      float64 `json:"ask_size"` //
	Volume       float64 `json:"volume"`   // volume for 24 hours
	Dtm          int     `json:"dtm"`      // the days until maturity
	Leverage     string  `json:"leverage"` // leverage
	Last         float64 `json:"last"`     // last trade price
	Change       float64 `json:"change"`   // 24h change
	OpenInterest float64 `json:"openInterest"`
}
