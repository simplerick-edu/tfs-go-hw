package krakenapi

import (
	"bot/domain"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (k *KrakenAPI) GetPositions() (*domain.OpenPositionsResponse, error) {
	u, _ := url.Parse(OpenPositionsURL)
	req, err := http.NewRequest(http.MethodGet, u.String(), strings.NewReader(u.RawQuery))
	if err != nil {
		return nil, err
	}
	req = k.privateRequest(req)
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, err
	}
	resp := &domain.OpenPositionsResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	return resp, err
}

//func (k *KrakenAPI) Buy(symbol string, price float64, size int64) (*domain.SendOrderResponse, error) {
//	return k.SendOrder(symbol, domain.Buy, IocType, price, size)
//}
//
//func (k *KrakenAPI) Sell(symbol string, price float64, size int64) (*domain.SendOrderResponse, error) {
//	return k.SendOrder(symbol, domain.Sell, IocType, price, size)
//}

func (k *KrakenAPI) SendOrder(symbol string, side domain.Action, orderType domain.OrderType, price float64, size int64) (*domain.SendOrderResponse, error) {
	u, _ := url.Parse(SendOrderURL)
	values := url.Values{}
	values.Set("symbol", symbol)
	values.Set("limitPrice", strconv.FormatFloat(price, 'f', 3, 64))
	values.Set("size", strconv.FormatInt(size, 10))
	values.Set("side", string(side))
	values.Set("orderType", string(orderType))
	u.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(u.RawQuery))
	if err != nil {
		return nil, err
	}
	req = k.privateRequest(req)
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, err
	}
	resp := &domain.SendOrderResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	return resp, err
}

func (k *KrakenAPI) CancelOrders() (*domain.CancelOrdersResponse, error) {
	u, _ := url.Parse(CancelOrdersURL)
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(u.RawQuery))
	if err != nil {
		return nil, err
	}
	req = k.privateRequest(req)
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, err
	}
	resp := &domain.CancelOrdersResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	return resp, err
}

// Helpers

func (k *KrakenAPI) privateRequest(req *http.Request) *http.Request {
	// Get sign
	body, _ := req.GetBody()
	postData, _ := io.ReadAll(body)
	endpointPath := []byte(req.URL.Path[12:])
	sign := generateSign(postData, endpointPath, k.privateKey)
	// Add headers
	req.Header.Add("APIKey", k.publicKey)
	req.Header.Add("Authent", sign)
	return req
}

func (k *KrakenAPI) sendRequest(req *http.Request) ([]byte, error) {
	resp, err := k.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	return bodyBytes, err
}

func generateSign(postData []byte, endpointPath []byte, privateKey string) string {
	// Concatenate postData + nonce + endpointPath
	src := append(postData, endpointPath...)
	//Hash the result of step 1 with the SHA-256 algorithm
	sha := sha256.New()
	sha.Write(src)
	//Base64-decode your api_secret
	apiSecret, _ := base64.StdEncoding.DecodeString(privateKey)
	//Use the result of step 3 to hash the result of the step 2 with the HMAC-SHA-512 algorithm
	h := hmac.New(sha512.New, apiSecret)
	h.Write(sha.Sum(nil))
	//Base64-encode the result of step 4
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}
