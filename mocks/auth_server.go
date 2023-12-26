// Code generated by MockGen. DO NOT EDIT.
// Source: authserver.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	authserver "github.com/konstellation-io/kli/authserver"
)

// MockAuthServerInterface is a mock of Authenticator interface.
type MockAuthServerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerInterfaceMockRecorder
}

// MockAuthServerInterfaceMockRecorder is the mock recorder for MockAuthServerInterface.
type MockAuthServerInterfaceMockRecorder struct {
	mock *MockAuthServerInterface
}

// NewMockAuthServerInterface creates a new mock instance.
func NewMockAuthServerInterface(ctrl *gomock.Controller) *MockAuthServerInterface {
	mock := &MockAuthServerInterface{ctrl: ctrl}
	mock.recorder = &MockAuthServerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServerInterface) EXPECT() *MockAuthServerInterfaceMockRecorder {
	return m.recorder
}

// StartServer mocks base method.
func (m *MockAuthServerInterface) Login(config authserver.KeycloakConfig) (*authserver.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", config)
	ret0, _ := ret[0].(*authserver.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartServer indicates an expected call of StartServer.
func (mr *MockAuthServerInterfaceMockRecorder) StartServer(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServerInterface)(nil).Login), config)
}
