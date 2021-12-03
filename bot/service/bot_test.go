package service

import (
	"bot/domain"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ExchangeMock struct {
	mock.Mock
}

func (exm *ExchangeMock) GetPositions() (*domain.OpenPositionsResponse, error) {
	args := exm.Called()
	return args.Get(0).(*domain.OpenPositionsResponse), args.Error(1)
}

func (exm *ExchangeMock) SendOrder(order domain.Order) (*domain.SendOrderResponse, error) {
	args := exm.Called(order)
	return args.Get(0).(*domain.SendOrderResponse), args.Error(1)
}

func (exm *ExchangeMock) Subscribe(instruments ...string) (<-chan domain.Ticker, error) {
	args := exm.Called(instruments)
	return args.Get(0).(chan domain.Ticker), args.Error(1)
}

func (exm *ExchangeMock) Unsubscribe() error {
	args := exm.Called()
	return args.Error(0)
}

type PredictorMock struct {
	mock.Mock
}

func (m *PredictorMock) Predict(tickers ...domain.Ticker) (float64, error) {
	args := m.Called(tickers)
	return args.Get(0).(float64), args.Error(1)
}

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) StoreEvent(ctx context.Context, event domain.OrderEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type NotifierMock struct {
	mock.Mock
}

func (m *NotifierMock) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *NotifierMock) Stop() {
	return
}

func (m *NotifierMock) Notify(text string) error {
	args := m.Called(text)
	return args.Error(0)
}

var defaultParams = Parameters{"str", 100, 2, 0.5, 10, 1}

var openPosSample = `{
     "result":"success",
     "openPositions":[
        {
         "side":"short",
         "symbol":"pi_xbtusd",
         "price":9392.749993345933,
         "fillTime":"2020-07-22T14:39:12.376Z",
         "size":10000,
         "unrealizedFunding":1.045432180096817E-5
         },
         {
          "side":"long",
          "symbol":"fi_xbtusd_201225",
          "price":9399.749966754434,
          "fillTime":"2020-07-22T14:39:12.376Z",
          "size":20000
          }
     ],
     "serverTime":"2020-07-22T14:39:12.376Z"
}`

var tickerSample = `{
  "time": 1612270825253,
  "feed": "ticker",
  "product_id": "PI_XBTUSD",
  "bid": 34832.5,
  "ask": 34847.5,
  "bid_size": 42864,
  "ask_size": 2300,
  "volume": 262306237,
  "dtm": 0,
  "leverage": "50x",
  "index": 34803.45,
  "premium": 0.1,
  "last": 34852,
  "change": 2.995109121267192,
  "funding_rate": 3.891007752e-9,
  "funding_rate_prediction": 4.2233756e-9,
  "suspended": false,
  "tag": "perpetual",
  "pair": "XBT:USD",
  "openInterest": 107706940,
  "markPrice": 34844.25,
  "maturityTime": 0,
  "relative_funding_rate": 0.000135046879166667,
  "relative_funding_rate_prediction": 0.000146960125,
  "next_funding_rate_time": 1612281600000
}`

var sendOrderRespSample = `{ 
   "result":"success",
   "sendStatus":{ 
      "order_id":"61ca5732-3478-42fe-8362-abbfd9465294",
      "status":"placed",
      "receivedTime":"2019-12-11T17:17:33.888Z",
      "orderEvents":[ 
         { 
            "executionId":"e1ec9f63-2338-4c44-b40a-43486c6732d7",
            "price":7244.5,
            "amount":10,
            "orderPriorEdit":null,
            "orderPriorExecution":{ 
               "orderId":"61ca5732-3478-42fe-8362-abbfd9465294",
               "cliOrdId":null,
               "type":"lmt",
               "symbol":"pi_xbtusd",
               "side":"buy",
               "quantity":10,
               "filled":0,
               "limitPrice":7500,
               "reduceOnly":false,
               "timestamp":"2019-12-11T17:17:33.888Z",
               "lastUpdateTimestamp":"2019-12-11T17:17:33.888Z"
            },
            "takerReducedQuantity":null,
            "type":"EXECUTION"
         }
      ]
   },
   "serverTime":"2019-12-11T17:17:33.888Z"
}`

func TestBot_FetchOpenPositions(t *testing.T) {
	exm := ExchangeMock{}
	var openPos domain.OpenPositionsResponse
	json.Unmarshal([]byte(openPosSample), &openPos)
	exm.On("GetPositions").Return(&openPos, nil)
	var bot = New(&exm, &NotifierMock{}, &StorageMock{}, &PredictorMock{}, defaultParams)
	err := bot.FetchOpenPositions()
	assert.Equal(t, err, nil)
}

func TestBot_Start(t *testing.T) {
	exm := ExchangeMock{}
	var openPos domain.OpenPositionsResponse
	json.Unmarshal([]byte(openPosSample), &openPos)
	exm.On("GetPositions").Return(&openPos, nil)
	var ticker domain.Ticker
	json.Unmarshal([]byte(tickerSample), &ticker)
	tickers := make(chan domain.Ticker)
	go func() {
		tickers <- ticker
	}()
	exm.On("Subscribe", mock.Anything).Return(tickers, nil)
	exm.On("Unsubscribe").Return(nil)
	nm := NotifierMock{}
	nm.On("Start").Return(nil)
	nm.On("Stop").Return()
	sm := StorageMock{}
	pm := PredictorMock{}
	pm.On("Predict", mock.Anything).Return(0.5, nil)
	var bot = New(&exm, &nm, &sm, &pm, defaultParams)
	err := bot.Start()
	assert.Equal(t, err, nil)
}

func TestBot_processSequence(t *testing.T) {
	exm := ExchangeMock{}
	var sendResp domain.SendOrderResponse
	json.Unmarshal([]byte(sendOrderRespSample), &sendResp)
	exm.On("SendOrder", mock.Anything).Return(&sendResp, nil)
	var ticker domain.Ticker
	json.Unmarshal([]byte(tickerSample), &ticker)
	tickers := []domain.Ticker{ticker, ticker}
	nm := NotifierMock{}
	nm.On("Notify", mock.Anything).Return(nil)
	sm := StorageMock{}
	sm.On("StoreEvent", mock.Anything, mock.Anything).Return(nil)
	pm := PredictorMock{}
	pm.On("Predict", mock.Anything).Return(0.6, nil)
	var bot = New(&exm, &nm, &sm, &pm, defaultParams)
	err := bot.processSequence(tickers)
	assert.Equal(t, err, nil)
}
