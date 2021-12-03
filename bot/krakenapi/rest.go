package krakenapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	// not-std
	"bot/domain"
)

func (k *KrakenAPI) GetPositions() (*domain.OpenPositionsResponse, error) {
	u, err := url.Parse(OpenPositionsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), strings.NewReader(u.RawQuery))
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	req, err = k.privateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't make private request: %w", err)
	}
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}
	resp := &domain.OpenPositionsResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (k *KrakenAPI) SendOrder(order domain.Order) (*domain.SendOrderResponse, error) {
	u, _ := url.Parse(SendOrderURL)
	values := url.Values{}
	values.Set("symbol", order.Symbol)
	//values.Set("limitPrice", strconv.FormatFloat(order.LimitPrice, 'f', 2, 64))
	values.Set("limitPrice", strconv.FormatInt(int64(order.LimitPrice), 10))
	values.Set("size", strconv.FormatInt(int64(order.Quantity), 10))
	values.Set("side", string(order.Side))
	values.Set("orderType", string(order.Type))
	u.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(u.RawQuery))
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	req, err = k.privateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't make private request: %w", err)
	}
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
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
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	req, err = k.privateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't make private request: %w", err)
	}
	respBody, err := k.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}
	resp := &domain.CancelOrdersResponse{}
	err = json.Unmarshal(respBody, resp)
	return resp, err
}

// Helpers

func (k *KrakenAPI) privateRequest(req *http.Request) (*http.Request, error) {
	// Get sign
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	postData, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	endpointPath := []byte(req.URL.Path[12:])
	sign := generateSign(postData, endpointPath, k.privateKey)
	// Add headers
	req.Header.Add("APIKey", k.publicKey)
	req.Header.Add("Authent", sign)
	return req, nil
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
