package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pustato/image-previewer/internal/client"
	"github.com/pustato/image-previewer/internal/resizer"
)

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
	var err error

	rsp, err := a.client.GetWithHeaders(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("ResizerApp get %s: %w", url, err)
	}
	defer func() {
		err = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("not ok") // todo
	}

	result, err := a.resizer.Resize(rsp.Body, w, h)
	if err != nil {
		return nil, fmt.Errorf("ResizerApp resize: %w", err)
	}

	return result, err
}
