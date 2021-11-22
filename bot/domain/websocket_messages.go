package domain

type Message struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIDs []string `json:"product_ids"`
}
