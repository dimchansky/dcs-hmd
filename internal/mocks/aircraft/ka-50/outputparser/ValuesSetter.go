// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ValuesSetter is an autogenerated mock type for the ValuesSetter type
type ValuesSetter struct {
	mock.Mock
}

// SetRotorPitch provides a mock function with given fields: val
func (_m *ValuesSetter) SetRotorPitch(val float64) {
	_m.Called(val)
}

// SetRotorRPM provides a mock function with given fields: val
func (_m *ValuesSetter) SetRotorRPM(val float64) {
	_m.Called(val)
}

// SetVerticalVelocity provides a mock function with given fields: val
func (_m *ValuesSetter) SetVerticalVelocity(val float64) {
	_m.Called(val)
}

type mockConstructorTestingTNewValuesSetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewValuesSetter creates a new instance of ValuesSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewValuesSetter(t mockConstructorTestingTNewValuesSetter) *ValuesSetter {
	mock := &ValuesSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
