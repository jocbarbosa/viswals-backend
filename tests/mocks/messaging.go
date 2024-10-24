// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	port "github.com/jocbarbosa/viswals-backend/internals/core/port"
	mock "github.com/stretchr/testify/mock"
)

// Messaging is an autogenerated mock type for the Messaging type
type Messaging struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Messaging) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: handler
func (_m *Messaging) Consume(handler port.MessageHandler) error {
	ret := _m.Called(handler)

	if len(ret) == 0 {
		panic("no return value specified for Consume")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(port.MessageHandler) error); ok {
		r0 = rf(handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Write provides a mock function with given fields: msg
func (_m *Messaging) Write(msg port.Message) error {
	ret := _m.Called(msg)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(port.Message) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMessaging creates a new instance of Messaging. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMessaging(t interface {
	mock.TestingT
	Cleanup(func())
}) *Messaging {
	mock := &Messaging{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
