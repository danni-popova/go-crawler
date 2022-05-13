// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// SyncedList is an autogenerated mock type for the SyncedList type
type SyncedList struct {
	mock.Mock
}

// Add provides a mock function with given fields: url
func (_m *SyncedList) Add(url string) {
	_m.Called(url)
}

// Contains provides a mock function with given fields: url
func (_m *SyncedList) Contains(url string) bool {
	ret := _m.Called(url)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PrintContents provides a mock function with given fields:
func (_m *SyncedList) PrintContents() {
	_m.Called()
}

// NewSyncedList creates a new instance of SyncedList. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSyncedList(t testing.TB) *SyncedList {
	mock := &SyncedList{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}