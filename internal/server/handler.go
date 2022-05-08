package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/goware/urlx"
	"github.com/pustato/image-previewer/internal/app"
	"github.com/pustato/image-previewer/internal/logger"
)

const (
	pathPartsExpected  = 4
	pathPartsWidthIdx  = 1
	pathPartsHeightIdx = 2
	pathPartsURLIdx    = 3
	badRequestText     = "bad request"
)

type Handler struct {
	app app.App
	log logger.Logger
}

type request struct {
	w, h int
	url  string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rq, err := parsePath(r.URL.Path)
	if err != nil {
		h.log.Warn("parse path " + r.URL.Path + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	resized, err := h.app.GetAndResize(r.Context(), rq.url, rq.w, rq.h, r.Header)
	if err != nil {
		h.log.Warn("get and resize: " + err.Error())
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(badRequestText))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resized)
}

func parsePath(path string) (*request, error) {
	parts := strings.SplitN(path, `/`, pathPartsExpected)
	if len(parts) != pathPartsExpected {
		return nil, ErrMalformedRequestPath
	}

	w, err := strconv.Atoi(parts[pathPartsWidthIdx])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parts[pathPartsWidthIdx], ErrWidthIsNotANumber)
	}

	h, err := strconv.Atoi(parts[pathPartsHeightIdx])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parts[pathPartsHeightIdx], ErrHeightIsNotANumber)
	}

	u, err := normalizeURL("http://" + parts[pathPartsURLIdx])
	if err != nil {
		return nil, err
	}

	return &request{w, h, u}, nil
}

func normalizeURL(u string) (string, error) {
	uu, err := urlx.Parse(strings.ToLower(u))
	if err != nil {
		return "", fmt.Errorf("parse url %s: %w: %s", u, ErrInvalidURL, err.Error())
	}

	uu.Fragment = ""

	normalURL, err := urlx.Normalize(uu)
	if err != nil {
		return "", fmt.Errorf("normalize url %s: %w: %s", u, ErrInvalidURL, err.Error())
	}

	return normalURL, nil
}
