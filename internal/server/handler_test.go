package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockapp "github.com/pustato/image-previewer/internal/app/mocks"
	mocklogger "github.com/pustato/image-previewer/internal/logger/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
	t.Parallel()

	testData := []struct {
		rq   *http.Request
		url  string
		w, h int
	}{
		{
			httptest.NewRequest(http.MethodGet, "http://x/100/200/www.example.com/image.jpg", nil),
			"http://www.example.com/image.jpg",
			100, 200,
		},
		{
			httptest.NewRequest(http.MethodGet, "http://x/1/1/www.example.com/image.jpg?param=not_encoded", nil),
			"http://www.example.com/image.jpg",
			1, 1,
		},
		{
			httptest.NewRequest(http.MethodGet, "http://x/1/1/www.example.com/image.jpg?#fragment_not_encoded", nil),
			"http://www.example.com/image.jpg",
			1, 1,
		},
		{
			httptest.NewRequest(http.MethodGet, "http://x/1024/768/www.example.com/Image.JPG", nil),
			"http://www.example.com/image.jpg",
			1024, 768,
		},
		{
			httptest.NewRequest(http.MethodGet, "http://x/12/14/www.example.com/Image.JPG%3Fv%3D1", nil),
			"http://www.example.com/image.jpg?v=1",
			12, 14,
		},
		{
			httptest.NewRequest(
				http.MethodGet,
				"http://x/12/14/www.example.com/Image.jpg%3Fa_order%3Dfirst%26b_order%3Dsecond",
				nil,
			),
			"http://www.example.com/image.jpg?a_order=first&b_order=second",
			12, 14,
		},
		{
			httptest.NewRequest(
				http.MethodGet,
				"http://x/12/14/www.example.com/Image.jpg%3Fb_order%3Dsecond%26a_order%3Dfirst",
				nil,
			),
			"http://www.example.com/image.jpg?a_order=first&b_order=second",
			12, 14,
		},
		{
			httptest.NewRequest(http.MethodGet, "http://x/1/1/www.example.com/Image.JPG%3Fv%3D1%23fragment", nil),
			"http://www.example.com/image.jpg?v=1",
			1, 1,
		},

		{
			func() *http.Request {
				rq := httptest.NewRequest(http.MethodGet, "http://x/1/1/www.example.com%2FImage.JPG%3Fv%3D1%23fragment", nil)

				rq.Header.Add("key", "value")
				rq.Header.Add("key2", "value2")

				return rq
			}(),
			"http://www.example.com/image.jpg?v=1",
			1, 1,
		},
	}

	for i, td := range testData {
		td := td
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			t.Parallel()
			result := []byte("success result")

			app := &mockapp.App{}
			app.
				On("GetAndResize", td.rq.Context(), td.url, td.w, td.h, td.rq.Header).
				Once().
				Return(result, nil)

			h := Handler{
				app: app,
				log: &mocklogger.Logger{},
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, td.rq)

			rsp := w.Result()
			body, _ := io.ReadAll(rsp.Body)

			require.Equal(t, http.StatusOK, rsp.StatusCode)
			require.EqualValues(t, result, body)

			rsp.Body.Close()
		})
	}
}

func TestHandler_ServeHTTP_PathErrors(t *testing.T) {
	t.Parallel()

	testData := []struct {
		url string
		err error
	}{
		{"http://x/", ErrMalformedRequestPath},
		{"http://x/1", ErrMalformedRequestPath},
		{"http://x/1/1", ErrMalformedRequestPath},

		{"http://x/nan/1/x", ErrWidthIsNotANumber},
		{"http://x/1.6/1/x", ErrWidthIsNotANumber},

		{"http://x/1/nan/x", ErrHeightIsNotANumber},
		{"http://x/1/22.4/x", ErrHeightIsNotANumber},

		{"http://x/1/22.4/x", ErrHeightIsNotANumber},
	}

	for i, td := range testData {
		td := td
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			t.Parallel()

			logg := &mocklogger.Logger{}
			logg.On("Warn", mock.MatchedBy(func(msg string) bool {
				return strings.Contains(msg, td.err.Error())
			}))

			w := httptest.NewRecorder()
			rq := httptest.NewRequest(http.MethodGet, td.url, nil)

			h := Handler{
				app: &mockapp.App{},
				log: logg,
			}

			h.ServeHTTP(w, rq)

			rsp := w.Result()
			body, _ := io.ReadAll(rsp.Body)

			require.Equal(t, http.StatusNotFound, rsp.StatusCode)
			require.Contains(t, string(body), td.err.Error())

			rsp.Body.Close()
		})
	}
}

func TestHandler_ServeHTTP_AppError(t *testing.T) {
	rq := httptest.NewRequest(http.MethodGet, "http://x/10/11/www.example.com/image.jpg", nil)
	w := httptest.NewRecorder()
	testError := errors.New("some app error")

	logg := &mocklogger.Logger{}
	logg.On("Warn", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, testError.Error())
	}))

	app := &mockapp.App{}
	app.
		On("GetAndResize", rq.Context(), "http://www.example.com/image.jpg", 10, 11, rq.Header).
		Once().
		Return(nil, testError)

	h := Handler{
		app: app,
		log: logg,
	}

	h.ServeHTTP(w, rq)

	rsp := w.Result()
	body, _ := io.ReadAll(rsp.Body)

	require.Equal(t, http.StatusBadGateway, rsp.StatusCode)
	require.Equal(t, string(body), badRequestText)

	rsp.Body.Close()
}
