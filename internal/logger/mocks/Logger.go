// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocklogger

import mock "github.com/stretchr/testify/mock"

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: msg
func (_m *Logger) Debug(msg string) {
	_m.Called(msg)
}

// Error provides a mock function with given fields: msg
func (_m *Logger) Error(msg string) {
	_m.Called(msg)
}

// Info provides a mock function with given fields: msg
func (_m *Logger) Info(msg string) {
	_m.Called(msg)
}

// Warn provides a mock function with given fields: msg
func (_m *Logger) Warn(msg string) {
	_m.Called(msg)
}
