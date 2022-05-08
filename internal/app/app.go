package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pustato/image-previewer/internal/client"
	"github.com/pustato/image-previewer/internal/resizer"
)

var ErrRequestError = errors.New("request error")

type App interface {
	GetAndResize(ctx context.Context, url string, w, h int, headers http.Header) ([]byte, error)
}

func NewResizerApp(c client.Client, r resizer.Resizer) *ResizerApp {
	return &ResizerApp{c, r}
}

type ResizerApp struct {
	client  client.Client
	resizer resizer.Resizer
}

func (a *ResizerApp) GetAndResize(ctx context.Context, url string, w, h int, headers http.Header) ([]byte, error) {
	rsp, err := a.client.GetWithHeaders(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("ResizerApp get %s: %w", url, err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, ErrRequestError
	}

	result, err := a.resizer.Resize(rsp.Body, w, h)
	if err != nil {
		return nil, fmt.Errorf("ResizerApp resize: %w", err)
	}

	return result, err
}
