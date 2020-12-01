// Code generated by MockGen. DO NOT EDIT.
// Source: server_client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/konstellation-io/kli/api"
	reflect "reflect"
)

// MockServerClienter is a mock of ServerClienter interface
type MockServerClienter struct {
	ctrl     *gomock.Controller
	recorder *MockServerClienterMockRecorder
}

// MockServerClienterMockRecorder is the mock recorder for MockServerClienter
type MockServerClienterMockRecorder struct {
	mock *MockServerClienter
}

// NewMockServerClienter creates a new mock instance
func NewMockServerClienter(ctrl *gomock.Controller) *MockServerClienter {
	mock := &MockServerClienter{ctrl: ctrl}
	mock.recorder = &MockServerClienterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServerClienter) EXPECT() *MockServerClienterMockRecorder {
	return m.recorder
}

// ListRuntimes mocks base method
func (m *MockServerClienter) ListRuntimes() (api.RuntimeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRuntimes")
	ret0, _ := ret[0].(api.RuntimeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRuntimes indicates an expected call of ListRuntimes
func (mr *MockServerClienterMockRecorder) ListRuntimes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRuntimes", reflect.TypeOf((*MockServerClienter)(nil).ListRuntimes))
}
