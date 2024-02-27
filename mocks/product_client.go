// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	kai "github.com/konstellation-io/kli/api/kai"
	configuration "github.com/konstellation-io/kli/internal/services/configuration"

	mock "github.com/stretchr/testify/mock"
)

// MockProductClient is an autogenerated mock type for the ProductClient type
type MockProductClient struct {
	mock.Mock
}

type MockProductClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProductClient) EXPECT() *MockProductClient_Expecter {
	return &MockProductClient_Expecter{mock: &_m.Mock}
}

// AddMaintainerToProduct provides a mock function with given fields: server, product, userEmail
func (_m *MockProductClient) AddMaintainerToProduct(server *configuration.Server, product string, userEmail string) error {
	ret := _m.Called(server, product, userEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) error); ok {
		r0 = rf(server, product, userEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductClient_AddMaintainerToProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMaintainerToProduct'
type MockProductClient_AddMaintainerToProduct_Call struct {
	*mock.Call
}

// AddMaintainerToProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - product string
//   - userEmail string
func (_e *MockProductClient_Expecter) AddMaintainerToProduct(server interface{}, product interface{}, userEmail interface{}) *MockProductClient_AddMaintainerToProduct_Call {
	return &MockProductClient_AddMaintainerToProduct_Call{Call: _e.mock.On("AddMaintainerToProduct", server, product, userEmail)}
}

func (_c *MockProductClient_AddMaintainerToProduct_Call) Run(run func(server *configuration.Server, product string, userEmail string)) *MockProductClient_AddMaintainerToProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockProductClient_AddMaintainerToProduct_Call) Return(_a0 error) *MockProductClient_AddMaintainerToProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductClient_AddMaintainerToProduct_Call) RunAndReturn(run func(*configuration.Server, string, string) error) *MockProductClient_AddMaintainerToProduct_Call {
	_c.Call.Return(run)
	return _c
}

// AddUserToProduct provides a mock function with given fields: server, product, userEmail
func (_m *MockProductClient) AddUserToProduct(server *configuration.Server, product string, userEmail string) error {
	ret := _m.Called(server, product, userEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) error); ok {
		r0 = rf(server, product, userEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductClient_AddUserToProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUserToProduct'
type MockProductClient_AddUserToProduct_Call struct {
	*mock.Call
}

// AddUserToProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - product string
//   - userEmail string
func (_e *MockProductClient_Expecter) AddUserToProduct(server interface{}, product interface{}, userEmail interface{}) *MockProductClient_AddUserToProduct_Call {
	return &MockProductClient_AddUserToProduct_Call{Call: _e.mock.On("AddUserToProduct", server, product, userEmail)}
}

func (_c *MockProductClient_AddUserToProduct_Call) Run(run func(server *configuration.Server, product string, userEmail string)) *MockProductClient_AddUserToProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockProductClient_AddUserToProduct_Call) Return(_a0 error) *MockProductClient_AddUserToProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductClient_AddUserToProduct_Call) RunAndReturn(run func(*configuration.Server, string, string) error) *MockProductClient_AddUserToProduct_Call {
	_c.Call.Return(run)
	return _c
}

// CreateProduct provides a mock function with given fields: server, name, description
func (_m *MockProductClient) CreateProduct(server *configuration.Server, name string, description string) error {
	ret := _m.Called(server, name, description)

	var r0 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) error); ok {
		r0 = rf(server, name, description)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductClient_CreateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProduct'
type MockProductClient_CreateProduct_Call struct {
	*mock.Call
}

// CreateProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - name string
//   - description string
func (_e *MockProductClient_Expecter) CreateProduct(server interface{}, name interface{}, description interface{}) *MockProductClient_CreateProduct_Call {
	return &MockProductClient_CreateProduct_Call{Call: _e.mock.On("CreateProduct", server, name, description)}
}

func (_c *MockProductClient_CreateProduct_Call) Run(run func(server *configuration.Server, name string, description string)) *MockProductClient_CreateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockProductClient_CreateProduct_Call) Return(_a0 error) *MockProductClient_CreateProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductClient_CreateProduct_Call) RunAndReturn(run func(*configuration.Server, string, string) error) *MockProductClient_CreateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetProduct provides a mock function with given fields: server, id
func (_m *MockProductClient) GetProduct(server *configuration.Server, id string) (*kai.Product, error) {
	ret := _m.Called(server, id)

	var r0 *kai.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string) (*kai.Product, error)); ok {
		return rf(server, id)
	}
	if rf, ok := ret.Get(0).(func(*configuration.Server, string) *kai.Product); ok {
		r0 = rf(server, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kai.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(*configuration.Server, string) error); ok {
		r1 = rf(server, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProductClient_GetProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProduct'
type MockProductClient_GetProduct_Call struct {
	*mock.Call
}

// GetProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - id string
func (_e *MockProductClient_Expecter) GetProduct(server interface{}, id interface{}) *MockProductClient_GetProduct_Call {
	return &MockProductClient_GetProduct_Call{Call: _e.mock.On("GetProduct", server, id)}
}

func (_c *MockProductClient_GetProduct_Call) Run(run func(server *configuration.Server, id string)) *MockProductClient_GetProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string))
	})
	return _c
}

func (_c *MockProductClient_GetProduct_Call) Return(_a0 *kai.Product, _a1 error) *MockProductClient_GetProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductClient_GetProduct_Call) RunAndReturn(run func(*configuration.Server, string) (*kai.Product, error)) *MockProductClient_GetProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetProducts provides a mock function with given fields: server
func (_m *MockProductClient) GetProducts(server *configuration.Server) ([]kai.Product, error) {
	ret := _m.Called(server)

	var r0 []kai.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(*configuration.Server) ([]kai.Product, error)); ok {
		return rf(server)
	}
	if rf, ok := ret.Get(0).(func(*configuration.Server) []kai.Product); ok {
		r0 = rf(server)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]kai.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(*configuration.Server) error); ok {
		r1 = rf(server)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProductClient_GetProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProducts'
type MockProductClient_GetProducts_Call struct {
	*mock.Call
}

// GetProducts is a helper method to define mock.On call
//   - server *configuration.Server
func (_e *MockProductClient_Expecter) GetProducts(server interface{}) *MockProductClient_GetProducts_Call {
	return &MockProductClient_GetProducts_Call{Call: _e.mock.On("GetProducts", server)}
}

func (_c *MockProductClient_GetProducts_Call) Run(run func(server *configuration.Server)) *MockProductClient_GetProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server))
	})
	return _c
}

func (_c *MockProductClient_GetProducts_Call) Return(_a0 []kai.Product, _a1 error) *MockProductClient_GetProducts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductClient_GetProducts_Call) RunAndReturn(run func(*configuration.Server) ([]kai.Product, error)) *MockProductClient_GetProducts_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveMaintainerFromProduct provides a mock function with given fields: server, product, userEmail
func (_m *MockProductClient) RemoveMaintainerFromProduct(server *configuration.Server, product string, userEmail string) error {
	ret := _m.Called(server, product, userEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) error); ok {
		r0 = rf(server, product, userEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductClient_RemoveMaintainerFromProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveMaintainerFromProduct'
type MockProductClient_RemoveMaintainerFromProduct_Call struct {
	*mock.Call
}

// RemoveMaintainerFromProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - product string
//   - userEmail string
func (_e *MockProductClient_Expecter) RemoveMaintainerFromProduct(server interface{}, product interface{}, userEmail interface{}) *MockProductClient_RemoveMaintainerFromProduct_Call {
	return &MockProductClient_RemoveMaintainerFromProduct_Call{Call: _e.mock.On("RemoveMaintainerFromProduct", server, product, userEmail)}
}

func (_c *MockProductClient_RemoveMaintainerFromProduct_Call) Run(run func(server *configuration.Server, product string, userEmail string)) *MockProductClient_RemoveMaintainerFromProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockProductClient_RemoveMaintainerFromProduct_Call) Return(_a0 error) *MockProductClient_RemoveMaintainerFromProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductClient_RemoveMaintainerFromProduct_Call) RunAndReturn(run func(*configuration.Server, string, string) error) *MockProductClient_RemoveMaintainerFromProduct_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUserFromProduct provides a mock function with given fields: server, product, userEmail
func (_m *MockProductClient) RemoveUserFromProduct(server *configuration.Server, product string, userEmail string) error {
	ret := _m.Called(server, product, userEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(*configuration.Server, string, string) error); ok {
		r0 = rf(server, product, userEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductClient_RemoveUserFromProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUserFromProduct'
type MockProductClient_RemoveUserFromProduct_Call struct {
	*mock.Call
}

// RemoveUserFromProduct is a helper method to define mock.On call
//   - server *configuration.Server
//   - product string
//   - userEmail string
func (_e *MockProductClient_Expecter) RemoveUserFromProduct(server interface{}, product interface{}, userEmail interface{}) *MockProductClient_RemoveUserFromProduct_Call {
	return &MockProductClient_RemoveUserFromProduct_Call{Call: _e.mock.On("RemoveUserFromProduct", server, product, userEmail)}
}

func (_c *MockProductClient_RemoveUserFromProduct_Call) Run(run func(server *configuration.Server, product string, userEmail string)) *MockProductClient_RemoveUserFromProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*configuration.Server), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockProductClient_RemoveUserFromProduct_Call) Return(_a0 error) *MockProductClient_RemoveUserFromProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductClient_RemoveUserFromProduct_Call) RunAndReturn(run func(*configuration.Server, string, string) error) *MockProductClient_RemoveUserFromProduct_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProductClient creates a new instance of MockProductClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProductClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProductClient {
	mock := &MockProductClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
