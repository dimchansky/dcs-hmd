// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MessageHandler is an autogenerated mock type for the MessageHandler type
type MessageHandler struct {
	mock.Mock
}

// HandleMessage provides a mock function with given fields: msg
func (_m *MessageHandler) HandleMessage(msg []byte) {
	_m.Called(msg)
}

type mockConstructorTestingTNewMessageHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageHandler creates a new instance of MessageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageHandler(t mockConstructorTestingTNewMessageHandler) *MessageHandler {
	mock := &MessageHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}