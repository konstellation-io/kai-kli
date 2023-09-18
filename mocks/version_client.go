// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	configuration "github.com/konstellation-io/kli/internal/services/configuration"

	mock "github.com/stretchr/testify/mock"
)

// MockVersionClient is an autogenerated mock type for the VersionClient type
type MockVersionClient struct {
	mock.Mock
}

type MockVersionClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockVersionClient) EXPECT() *MockVersionClient_Expecter {
	return &MockVersionClient_Expecter{mock: &_m.Mock}
}

// Push provides a mock function with given fields: server, product, krtFilePath
func (_m *MockVersionClient) Push(server *configuration.Server, product string, krtFilePath string) (string, error) {
	ret := _m.Called(server, product, krtFilePath)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) (string, error)); ok {
		return rf(server, product, krtFilePath)
	}
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) string); ok {
		r0 = rf(server, product, krtFilePath)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*configuration.Server, string, string) error); ok {
		r1 = rf(server, product, krtFilePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockVersionClient_Push_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Push'
type MockVersionClient_Push_Call struct {
	*mock.Call
}

// Push is a helper method to define mock.On call
//   - server *configuration.Server
//   - product string
//   - krtFilePath string
func (_e *MockVersionClient_Expecter) Push(server interface{}, product interface{}, krtFilePath interface{}) *MockVersionClient_Push_Call {
	return &MockVersionClient_Push_Call{Call: _e.mock.On("Push", server, product, krtFilePath)}
}

func (_c *MockVersionClient_Push_Call) Run(run func(server *configuration.Server, product string, krtFilePath string)) *MockVersionClient_Push_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockVersionClient_Push_Call) Return(startedVersionTag string, err error) *MockVersionClient_Push_Call {
	_c.Call.Return(startedVersionTag, err)
	return _c
}

func (_c *MockVersionClient_Push_Call) RunAndReturn(run func(*configuration.Server, string, string) (string, error)) *MockVersionClient_Push_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: server, productID, versionTag, comment
func (_m *MockVersionClient) Start(server *configuration.Server, productID string, versionTag string, comment string) (string, error) {
	ret := _m.Called(server, productID, versionTag, comment)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string, string) (string, error)); ok {
		return rf(server, productID, versionTag, comment)
	}
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string, string) string); ok {
		r0 = rf(server, productID, versionTag, comment)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*configuration.Server, string, string, string) error); ok {
		r1 = rf(server, productID, versionTag, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockVersionClient_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type MockVersionClient_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - server *configuration.Server
//   - productID string
//   - versionTag string
//   - comment string
func (_e *MockVersionClient_Expecter) Start(server interface{}, productID interface{}, versionTag interface{}, comment interface{}) *MockVersionClient_Start_Call {
	return &MockVersionClient_Start_Call{Call: _e.mock.On("Start", server, productID, versionTag, comment)}
}

func (_c *MockVersionClient_Start_Call) Run(run func(server *configuration.Server, productID string, versionTag string, comment string)) *MockVersionClient_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockVersionClient_Start_Call) Return(_a0 string, _a1 error) *MockVersionClient_Start_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockVersionClient_Start_Call) RunAndReturn(run func(*configuration.Server, string, string, string) (string, error)) *MockVersionClient_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockVersionClient creates a new instance of MockVersionClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockVersionClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockVersionClient {
	mock := &MockVersionClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
