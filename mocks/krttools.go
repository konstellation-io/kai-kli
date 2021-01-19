// Code generated by MockGen. DO NOT EDIT.
// Source: krttools.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockKrtTooler is a mock of KrtTooler interface
type MockKrtTooler struct {
	ctrl     *gomock.Controller
	recorder *MockKrtToolerMockRecorder
}

// MockKrtToolerMockRecorder is the mock recorder for MockKrtTooler
type MockKrtToolerMockRecorder struct {
	mock *MockKrtTooler
}

// NewMockKrtTooler creates a new mock instance
func NewMockKrtTooler(ctrl *gomock.Controller) *MockKrtTooler {
	mock := &MockKrtTooler{ctrl: ctrl}
	mock.recorder = &MockKrtToolerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKrtTooler) EXPECT() *MockKrtToolerMockRecorder {
	return m.recorder
}

// Validate mocks base method
func (m *MockKrtTooler) Validate(yamlPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", yamlPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockKrtToolerMockRecorder) Validate(yamlPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockKrtTooler)(nil).Validate), yamlPath)
}

// Build mocks base method
func (m *MockKrtTooler) Build(src, target string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", src, target)
	ret0, _ := ret[0].(error)
	return ret0
}

// Build indicates an expected call of Build
func (mr *MockKrtToolerMockRecorder) Build(src, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockKrtTooler)(nil).Build), src, target)
}
