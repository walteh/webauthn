// Code generated by mockery v2.39.1. DO NOT EDIT.

package mockery

import mock "github.com/stretchr/testify/mock"

// MockProvider_relyingparty is an autogenerated mock type for the Provider type
type MockProvider_relyingparty struct {
	mock.Mock
}

type MockProvider_relyingparty_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProvider_relyingparty) EXPECT() *MockProvider_relyingparty_Expecter {
	return &MockProvider_relyingparty_Expecter{mock: &_m.Mock}
}

// RPDisplayName provides a mock function with given fields:
func (_m *MockProvider_relyingparty) RPDisplayName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPDisplayName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockProvider_relyingparty_RPDisplayName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RPDisplayName'
type MockProvider_relyingparty_RPDisplayName_Call struct {
	*mock.Call
}

// RPDisplayName is a helper method to define mock.On call
func (_e *MockProvider_relyingparty_Expecter) RPDisplayName() *MockProvider_relyingparty_RPDisplayName_Call {
	return &MockProvider_relyingparty_RPDisplayName_Call{Call: _e.mock.On("RPDisplayName")}
}

func (_c *MockProvider_relyingparty_RPDisplayName_Call) Run(run func()) *MockProvider_relyingparty_RPDisplayName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockProvider_relyingparty_RPDisplayName_Call) Return(_a0 string) *MockProvider_relyingparty_RPDisplayName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProvider_relyingparty_RPDisplayName_Call) RunAndReturn(run func() string) *MockProvider_relyingparty_RPDisplayName_Call {
	_c.Call.Return(run)
	return _c
}

// RPID provides a mock function with given fields:
func (_m *MockProvider_relyingparty) RPID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockProvider_relyingparty_RPID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RPID'
type MockProvider_relyingparty_RPID_Call struct {
	*mock.Call
}

// RPID is a helper method to define mock.On call
func (_e *MockProvider_relyingparty_Expecter) RPID() *MockProvider_relyingparty_RPID_Call {
	return &MockProvider_relyingparty_RPID_Call{Call: _e.mock.On("RPID")}
}

func (_c *MockProvider_relyingparty_RPID_Call) Run(run func()) *MockProvider_relyingparty_RPID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockProvider_relyingparty_RPID_Call) Return(_a0 string) *MockProvider_relyingparty_RPID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProvider_relyingparty_RPID_Call) RunAndReturn(run func() string) *MockProvider_relyingparty_RPID_Call {
	_c.Call.Return(run)
	return _c
}

// RPOrigin provides a mock function with given fields:
func (_m *MockProvider_relyingparty) RPOrigin() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPOrigin")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockProvider_relyingparty_RPOrigin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RPOrigin'
type MockProvider_relyingparty_RPOrigin_Call struct {
	*mock.Call
}

// RPOrigin is a helper method to define mock.On call
func (_e *MockProvider_relyingparty_Expecter) RPOrigin() *MockProvider_relyingparty_RPOrigin_Call {
	return &MockProvider_relyingparty_RPOrigin_Call{Call: _e.mock.On("RPOrigin")}
}

func (_c *MockProvider_relyingparty_RPOrigin_Call) Run(run func()) *MockProvider_relyingparty_RPOrigin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockProvider_relyingparty_RPOrigin_Call) Return(_a0 string) *MockProvider_relyingparty_RPOrigin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProvider_relyingparty_RPOrigin_Call) RunAndReturn(run func() string) *MockProvider_relyingparty_RPOrigin_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProvider_relyingparty creates a new instance of MockProvider_relyingparty. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProvider_relyingparty(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProvider_relyingparty {
	mock := &MockProvider_relyingparty{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
