// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocklru

import (
	lru "github.com/pustato/image-previewer/internal/cache/lru"
	mock "github.com/stretchr/testify/mock"
)

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *Cache) Get(key string) (*lru.Item, bool) {
	ret := _m.Called(key)

	var r0 *lru.Item
	if rf, ok := ret.Get(0).(func(string) *lru.Item); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*lru.Item)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Set provides a mock function with given fields: key, item
func (_m *Cache) Set(key string, item *lru.Item) bool {
	ret := _m.Called(key, item)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, *lru.Item) bool); ok {
		r0 = rf(key, item)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
