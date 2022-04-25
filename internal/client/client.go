package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var _ Client = (*client)(nil)

type Client interface {
	GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error)
}

type client struct {
	client *http.Client
}

func New(timeout time.Duration) Client {
	return &client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client create request: %w", err)
	}

	for key, hh := range headers {
		for _, h := range hh {
			rq.Header.Add(key, h)
		}
	}

	rsp, err := c.client.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("client do request: %w", err)
	}

	return rsp, nil
}
