package app

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mockclient "github.com/pustato/image-previewer/internal/client/mocks"
	mockresizer "github.com/pustato/image-previewer/internal/resizer/mocks"
	"github.com/stretchr/testify/require"
)

var (
	ctx     = context.Background()
	url     = "http://test.url/"
	headers = http.Header{}
)

type bodyStub struct{}

func (t *bodyStub) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (t *bodyStub) Close() error {
	return nil
}

func TestResizerApp_GetAndResize_Success(t *testing.T) {
	body := &bodyStub{}
	rsp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	expectedResult := []byte("expected result")

	client := &mockclient.Client{}
	client.
		On("GetWithHeaders", ctx, url, headers).
		Once().
		Return(rsp, nil)

	resizer := &mockresizer.Resizer{}
	resizer.
		On("Resize", body, 100, 100).
		Once().
		Return(expectedResult, nil)

	app := NewResizerApp(client, resizer)

	res, err := app.GetAndResize(ctx, url, 100, 100, headers)
	require.NoError(t, err)
	require.EqualValues(t, expectedResult, res)
}

func TestResizerApp_GetAndResize_Errors(t *testing.T) {
	t.Parallel()

	t.Run("client GetWithHeaders error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")

		client := &mockclient.Client{}
		client.
			On("GetWithHeaders", ctx, url, headers).
			Once().
			Return(nil, expectedError)

		resizer := &mockresizer.Resizer{}

		app := NewResizerApp(client, resizer)
		res, err := app.GetAndResize(ctx, url, 100, 100, headers)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("client GetWithHeaders request error", func(t *testing.T) {
		t.Parallel()

		rsp := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       &bodyStub{},
		}

		client := &mockclient.Client{}
		client.
			On("GetWithHeaders", ctx, url, headers).
			Once().
			Return(rsp, nil)

		resizer := &mockresizer.Resizer{}
		app := NewResizerApp(client, resizer)
		res, err := app.GetAndResize(ctx, url, 100, 100, headers)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrRequestError)
	})

	t.Run("resizer error", func(t *testing.T) {
		t.Parallel()

		body := &bodyStub{}
		rsp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}
		expectedError := errors.New("expected error")

		client := &mockclient.Client{}
		client.
			On("GetWithHeaders", ctx, url, headers).
			Once().
			Return(rsp, nil)

		resizer := &mockresizer.Resizer{}
		resizer.
			On("Resize", body, 100, 100).
			Once().
			Return(nil, expectedError)

		app := NewResizerApp(client, resizer)
		res, err := app.GetAndResize(ctx, url, 100, 100, headers)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
	})
}
