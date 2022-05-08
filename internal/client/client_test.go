package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestGetWithHeaders(t *testing.T) {
	t.Run("proxy headers", func(t *testing.T) {
		url := "https://domain.zone/some/path/?query=param#fragment"
		headers := http.Header{}
		headers.Add("key1", "value1")
		headers.Add("key1", "value2")
		headers.Add("key2", "value1")

		client := NewHTTPClient(time.Second).
			WithRoundTripFunc(func(rq *http.Request) (*http.Response, error) {
				v1 := rq.Header.Values("key1")
				require.Len(t, v1, 2)
				require.Equal(t, v1[0], "value1")
				require.Equal(t, v1[1], "value2")

				v2 := rq.Header.Values("key2")
				require.Len(t, v2, 1)
				require.Equal(t, v2[0], "value1")

				require.Equal(t, url, rq.URL.String())

				return nil, nil
			})

		_, _ = client.GetWithHeaders(ctx, url, headers) //nolint:bodyclose
	})

	t.Run("proxy error", func(t *testing.T) {
		expectedErr := errors.New("test error")

		client := NewHTTPClient(time.Second).
			WithRoundTripFunc(func(rq *http.Request) (*http.Response, error) {
				return nil, expectedErr
			})

		rsp, err := client.GetWithHeaders(ctx, "http://google.com", http.Header{}) //nolint:bodyclose
		require.Nil(t, rsp)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("response as is", func(t *testing.T) {
		expectedRsp := &http.Response{}

		client := NewHTTPClient(time.Second).WithRoundTripFunc(func(rq *http.Request) (*http.Response, error) {
			return expectedRsp, nil
		})

		rsp, err := client.GetWithHeaders(ctx, "http://google.com", http.Header{})
		require.NoError(t, err)
		require.Equal(t, expectedRsp, rsp)

		rsp.Body.Close()
	})
}
