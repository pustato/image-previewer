// Code generated by mockery v2.10.2. DO NOT EDIT.

package mockresizer

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Resizer is an autogenerated mock type for the Resizer type
type Resizer struct {
	mock.Mock
}

// Resize provides a mock function with given fields: i, w, h
func (_m *Resizer) Resize(i io.Reader, w int, h int) ([]byte, error) {
	ret := _m.Called(i, w, h)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(io.Reader, int, int) []byte); ok {
		r0 = rf(i, w, h)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader, int, int) error); ok {
		r1 = rf(i, w, h)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
