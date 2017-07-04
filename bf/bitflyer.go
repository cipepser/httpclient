package bf

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cipepser/httpclient/sdk"
)

const (
	// TimeLayout is a format of the date contains bitflyer
	// api response and it is used in time package.
	TimeLayout = "2006-01-02T03:04:05"

	timeout = 10 // [sec]

	// URL is a endpoint of bitflyer api.
	URL = "https://api.bitflyer.jp"
)

// Client is a wrapper of sdk.Client for bitflyer api.
type Client struct {
	sdk.Client
}

// ************** public API **************

// Board is a json struct for market board information.
type Board struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []struct {
		Price float64 `json:"price"`
		Size  float64 `json:"size"`
	} `json:"bids"`
	Asks []struct {
		Price float64 `json:"price"`
		Size  float64 `json:"size"`
	} `json:"asks"`
}

// GetBoard returns makert board information.
// product is a paramter represented the makert you want to get information.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
func (c *Client) GetBoard(product string) (Board, error) {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// prepare query parameters.
	vals := url.Values{}
	if product != "" {
		vals.Add("product_code", product)
	}

	// make new request to get market board information.
	req, err := c.NewRequest(ctx, "GET", "/v1/getboard", nil)
	if err != nil {
		return Board{}, err
	}

	// embed vals to URL as query paramters.
	req.URL.RawQuery = vals.Encode()

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Board{}, err
	}

	// decode http response as type Board.
	b := Board{}
	err = sdk.DecodeBody(resp, &b)
	if err != nil {
		return Board{}, err
	}

	return b, nil
}

// Ticker is a json struct for market ticker information.
type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          float64 `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// GetTicker returns makert ticker information.
// product is a paramter represented the makert you want to get information.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
func (c *Client) GetTicker(product string) (Ticker, error) {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// prepare query parameters.
	vals := url.Values{}
	if product != "" {
		vals.Add("product_code", product)
	}

	// make new request to get market ticker information.
	req, err := c.NewRequest(ctx, "GET", "/v1/getticker", nil)
	if err != nil {
		return Ticker{}, err
	}

	// embed vals to URL as query paramters.
	req.URL.RawQuery = vals.Encode()

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Ticker{}, err
	}

	// decode http response as type Ticker.
	t := Ticker{}
	err = sdk.DecodeBody(resp, &t)
	if err != nil {
		return Ticker{}, err
	}

	return t, nil
}

// Execution is a json struct for market executions information.
type Execution struct {
	ID                         float64 `json:"id"`
	Side                       string  `json:"side"`
	Price                      float64 `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

// GetExecutions returns makert executions.
// [PARAMTERS]
// product : makert you want to get information.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
// count : the number of result.
// before : get the result which have smaller `id` than the `before`.
// after : get the result which have bigger `id` than the `after`.
func (c *Client) GetExecutions(product, count, before, after string) ([]Execution, error) {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// prepare query parameters.
	vals := url.Values{}
	if product != "" {
		vals.Add("product_code", product)
	}
	if count != "" {
		vals.Add("count", count)
	}
	if before != "" {
		vals.Add("before", before)
	}
	if after != "" {
		vals.Add("after", after)
	}

	// make new request to get market executions.
	req, err := c.NewRequest(ctx, "GET", "/v1/getexecutions", nil)
	if err != nil {
		return nil, err
	}

	// embed vals to URL as query paramters.
	req.URL.RawQuery = vals.Encode()

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// decode http response as type []Execution.
	es := []Execution{}
	err = sdk.DecodeBody(resp, &es)
	if err != nil {
		return nil, err
	}

	return es, nil
}

// ************** private API **************

// Balance is a json struct for your private blance information.
type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

// GetBalances returns your private balances information.
func (c *Client) GetBalances() ([]Balance, error) {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// make new request to get your private collateral information.
	spath := "/v1/me/getbalance"
	method := "GET"
	req, err := c.NewRequest(ctx, method, spath, nil)
	if err != nil {
		return nil, err
	}

	// set authentication header to req
	setPrivateHeader(req, method, spath, "")

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// decode http response as type []Balance.
	bs := []Balance{}
	err = sdk.DecodeBody(resp, &bs)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

// Collateral is a json struct for your private collateral information.
type Collateral struct {
	Collateral        float64 `json:"collateral"`
	OpenPositionPnl   float64 `json:"open_position_pnl"`
	RequireCollateral float64 `json:"require_collateral"`
	KeepRate          float64 `json:"keep_rate"`
}

// GetCollateral returns your private collateral information.
func (c *Client) GetCollateral() (Collateral, error) {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// make new request to get your private collateral information.
	spath := "/v1/me/getcollateral"
	method := "GET"
	req, err := c.NewRequest(ctx, method, spath, nil)
	if err != nil {
		return Collateral{}, err
	}

	// set authentication header to req
	setPrivateHeader(req, method, spath, "")

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Collateral{}, err
	}

	// decode http response as type Collateral.
	col := Collateral{}
	err = sdk.DecodeBody(resp, &col)
	if err != nil {
		return Collateral{}, err
	}

	return col, nil
}

// ChildOrder is a json struct to send new order.
type ChildOrder struct {
	ProductCode string `json:"product_code"`
	// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".

	ChildOrderType string `json:"child_order_type"`
	// "LIMIT"(指値) or "MARKET"(成行).

	Side string `json:"side"`
	// "BUY" or "SELL".

	Price float64 `json:"price"`
	Size  float64 `json:"size"`

	MinuteToExpire float64 `json:"minute_to_expire"`
	// The time your order will be expired, default 43200[min].

	TimeInForce string `json:"time_in_force"`
	// "GTC", "IOC" or "FOK".
	// you can confirm the details in "https://lightning.bitflyer.jp/docs/specialorder#執行数量条件".
}

// ChildOrderResponse is a json struct for the response of func SendNewOrder().
// type ChildOrderResponse struct {
// 	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
// }

// SendNewOrder sent a new order to the market.
// If successfully ordered, returns 0, although it returns -1 when the order is failed.
func (c *Client) SendNewOrder(b ChildOrder) int {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// make body from ChildOrder b.
	body, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	// make new request to send order.
	method := "POST"
	spath := "/v1/me/sendchildorder"
	req, err := c.NewRequest(ctx, method, spath, strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	// set authentication header to req
	setPrivateHeader(req, method, spath, string(body))

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		return 0
	}

	return -1
}

// GetMyOrderResponse is a json struct for the response of func GetMyOrder(),
// strictly GetMyOrder() returns the slice of GetMyOrderResponse, []GetMyOrderResponse.
type GetMyOrderResponse struct {
	ID                     float64 `json:"id"`
	ChildOrderID           string  `json:"child_order_id"`
	ProductCode            string  `json:"product_code"`
	Side                   string  `json:"side"`
	ChildOrderType         string  `json:"child_order_type"`
	Price                  float64 `json:"price"`
	AveragePrice           float64 `json:"average_price"`
	Size                   float64 `json:"size"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
}

// GetMyOrder gets the list of your orders.
// [PARAMTERS]
// product : makert you want to get information.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
// count : the number of result.
// before : get the result which have smaller `id` than the `before`.
// after : get the result which have bigger `id` than the `after`.
// childOrderState: the order status you want to get, default "ACTIVE"
// "ACTIVE", "COMPLETED", "CANCELED", "EXPIRED", "REJECTED"
func (c *Client) GetMyOrder(product, count, before, after, childOrderState string) []GetMyOrderResponse {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// prepare query parameters.
	vals := url.Values{}
	if product != "" {
		vals.Add("product_code", product)
	}
	if count != "" {
		vals.Add("count", count)
	}
	if before != "" {
		vals.Add("before", before)
	}
	if after != "" {
		vals.Add("after", after)
	}
	if childOrderState != "" {
		vals.Add("child_order_state", childOrderState)
	}

	// make new request to send order.
	method := "GET"
	spath := "/v1/me/getchildorders"
	req, err := c.NewRequest(ctx, method, spath, nil)
	if err != nil {
		log.Fatal(err)
	}

	// embed vals to URL as query paramters.
	req.URL.RawQuery = vals.Encode()

	// set authentication header to req
	setPrivateHeader(req, method, req.URL.RequestURI(), "")

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// 	fmt.Println(resp.Status)
	// 	a, _ := ioutil.ReadAll(resp.Body)
	// 	fmt.Println(string(a))

	// decode http response as type []GetMyOrderResponse.
	odrs := []GetMyOrderResponse{}
	err = sdk.DecodeBody(resp, &odrs)
	if err != nil {
		log.Fatal(err)
	}

	return odrs
}

// ChildOrderCanceled is a json struct to cancel the product.
type ChildOrderCanceled struct {
	ProductCode string `json:"product_code"`
	// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
	ChildOrderID string `json:"child_order_id"`
	// The oder ID you want to cancel.
}

// CancelOrder cancels your ACTIVE orders specified by ChildOrderID.
// product is a paramter represented the makert you want to cancel the orders.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
func (c *Client) CancelOrder(b ChildOrderCanceled) error {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// make body from ChildOrder b.
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	// make new request to send order.
	method := "POST"
	spath := "/v1/me/cancelchildorder"
	req, err := c.NewRequest(ctx, method, spath, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	// set authentication header to req
	setPrivateHeader(req, method, spath, string(body))

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("HTTP response code: " + resp.Status)
}

// ChildOrderAllCanceled is a json struct to cancel the product.
type ChildOrderAllCanceled struct {
	ProductCode string `json:"product_code"`
	// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
}

// CancelAllOrder cancels all your ACTIVE orders.
// product is a paramter represented the makert you want to cancel the orders.
// e.g. "BTC_JPY", "FX_BTC_JPY", "ETH_BTC".
func (c *Client) CancelAllOrder(b ChildOrderAllCanceled) error {
	// set timeout timer by context package.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// make body from ChildOrder b.
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	// make new request to send order.
	method := "POST"
	spath := "/v1/me/cancelallchildorders"
	req, err := c.NewRequest(ctx, method, spath, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	// set authentication header to req
	setPrivateHeader(req, method, spath, string(body))

	// send a http request and get a response.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("HTTP response code: " + resp.Status)
}

// setPrivateHeader sets authentication header to req.
func setPrivateHeader(req *http.Request, method, spath, body string) {
	key := os.Getenv("BFKEY")
	secret := os.Getenv("BFSECRET")

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	sign := makeHMAC(timestamp+method+spath+body, secret)

	req.Header.Set("ACCESS-KEY", key)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("ACCESS-SIGN", sign)
	req.Header.Set("Content-Type", "application/json")
}

// makeHMAC returns a HMAC by sha256.
func makeHMAC(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
