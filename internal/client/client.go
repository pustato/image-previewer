package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var _ Client = (*HttpClient)(nil)

type Client interface {
	GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error)
}

type HttpClient struct {
	client *http.Client
}

func New(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HttpClient) GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("HttpClient create request: %w", err)
	}

	for key, hh := range headers {
		for _, h := range hh {
			rq.Header.Add(key, h)
		}
	}

	rsp, err := c.client.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("HttpClient do request: %w", err)
	}

	return rsp, nil
}
