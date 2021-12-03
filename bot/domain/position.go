package domain

type Position struct {
	Side   string  `json:"side"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Size   int64   `json:"size"`
}
