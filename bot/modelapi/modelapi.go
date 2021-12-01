package modelapi

import (
	"bot/domain"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type vector []float64

type tfServeData struct {
	SignatureName string     `json:"signature_name"`
	Instances     [][]vector `json:"instances"`
}

type tfServeResponse struct {
	Predictions []vector `json:"predictions"`
}

type ModelService struct {
	url string
}

func New(url string) *ModelService {
	return &ModelService{url}
}

func (m *ModelService) Predict(tickers ...domain.Ticker) (float64, error) {
	sequence := toSequence(tickers)
	data := tfServeData{"serving_default",
		[][]vector{sequence}}
	postBody, _ := json.Marshal(data)
	resp, err := http.Post(m.url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	tfResp := &tfServeResponse{}
	json.Unmarshal(bodyBytes, tfResp)
	return tfResp.Predictions[0][0], nil
}

func toSequence(tickers []domain.Ticker) []vector {
	sequence := make([]vector, 0, len(tickers))
	for _, ticker := range tickers {
		vec := vector{
			ticker.Bid,
			ticker.Ask,
			float64(ticker.BidSize),
			float64(ticker.AskSize),
			float64(ticker.Volume),
			float64(ticker.Dtm),
			float64(ticker.Last),
			ticker.Change,
			float64(ticker.OpenInterest),
		}
		sequence = append(sequence, vec)
	}
	return sequence
}
