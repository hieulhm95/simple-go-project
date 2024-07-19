package ad_listing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	BaseUrl = "https://gateway.chotot.com/v1/public/ad-listing"
	CateVeh = "2000"
	CatePty = "1000"
)

type Option func(*client)

func NewClient(baseUrl string, retryTimes int, log *log.Logger) *client {
	// TODO #4 refactor NewClient using functional options
	return &client{
		httpClient: http.DefaultClient,
		baseUrl:    baseUrl,
		retryTimes: retryTimes,
		logger:     log,
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func WithUrl(url string) Option {
	return func(c *client) {
		c.baseUrl = url
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(c *client) {
		c.logger = logger
	}
}

func WithRetryTimes(retryTimes int) Option {
	return func(c *client) {
		c.retryTimes = retryTimes
	}
}

func NewClient2(opt ...Option) *client {
	c := new(client)
	for _, o := range opt {
		o(c)
	}

	// Set a default client if one was not provided, same as pattern :v
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.logger == nil {
		c.logger = log.Default()
	}

	return c
}

type client struct {
	httpClient *http.Client
	baseUrl    string
	retryTimes int
	logger     *log.Logger
}

func (c *client) GetAdByCate(ctx context.Context, cate string) (*AdsResponse, error) {
	now := time.Now()
	defer func() {
		c.logger.Printf("GetAdByCate Request - Cate %v, Duration: %v", cate, time.Since(now).String())
	}()

	url := fmt.Sprintf("%v?cg=%v&limit=10", BaseUrl, cate)

	// TODO #3 implement retry if StatusCode = 5xx
	resp, err := c.httpClient.Get(url)
	for i := 0; i < c.retryTimes; i++ {
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 500 {
			break
		}
		resp, err = c.httpClient.Get(url)
	}
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nResponse %v", string(b))

	var adResp AdsResponse
	// TODO #2 unmarshal json
	if err := json.Unmarshal(b, &adResp); err != nil {
		return nil, err
	}

	return &adResp, nil
}

type AdsResponse struct {
	Total int  `json:"total"`
	Ads   []Ad `json:"ads"`
}

type Ad struct {
	AdId int `json:"ad_id"`
	//TODO #1 Define struct
	// list_id , account_name, subject, list_time
	ListId      int    `json:"list_id"`
	AccountName string `json:"account_name"`
	Subject     string `json:"subject"`
	ListTime    int    `json:"list_time"`
}
