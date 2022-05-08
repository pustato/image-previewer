// Code generated by mockery v2.10.2. DO NOT EDIT.

package mockclient

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetWithHeaders provides a mock function with given fields: ctx, url, headers
func (_m *Client) GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	ret := _m.Called(ctx, url, headers)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string, http.Header) *http.Response); ok {
		r0 = rf(ctx, url, headers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, http.Header) error); ok {
		r1 = rf(ctx, url, headers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
