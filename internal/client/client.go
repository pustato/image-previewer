package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var _ Client = (*HTTPClient)(nil)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

type Client interface {
	GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error)
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) WithRoundTripFunc(f RoundTripperFunc) *HTTPClient {
	c.client.Transport = f

	return c
}

func (c *HTTPClient) GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTPClient create request: %w", err)
	}

	for key, hh := range headers {
		for _, h := range hh {
			rq.Header.Add(key, h)
		}
	}

	rsp, err := c.client.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("HTTPClient do request: %w", err)
	}

	return rsp, nil
}
