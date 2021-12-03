package modelapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	// not-std
	"bot/domain"
)

type vector []float64

const DefaultSignatureName = "serving_default"

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
	data := tfServeData{DefaultSignatureName,
		[][]vector{sequence}}
	postBody, _ := json.Marshal(data)
	resp, err := http.Post(m.url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	tfResp := &tfServeResponse{}
	err = json.Unmarshal(bodyBytes, tfResp)
	if err != nil {
		return 0, err
	}
	return tfResp.Predictions[0][0], nil
}

func toSequence(tickers []domain.Ticker) []vector {
	sequence := make([]vector, 0, len(tickers))
	for _, ticker := range tickers {
		vec := vector{
			ticker.Bid,
			ticker.Ask,
			ticker.BidSize,
			ticker.AskSize,
			ticker.Volume,
			float64(ticker.Dtm),
			ticker.Last,
			ticker.Change,
			ticker.OpenInterest,
		}
		sequence = append(sequence, vec)
	}
	return sequence
}
