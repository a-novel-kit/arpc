// Code generated by mockery v2.46.0. DO NOT EDIT.

package arpcmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockExecService is an autogenerated mock type for the ExecService type
type MockExecService[In interface{}, Out interface{}] struct {
	mock.Mock
}

type MockExecService_Expecter[In interface{}, Out interface{}] struct {
	mock *mock.Mock
}

func (_m *MockExecService[In, Out]) EXPECT() *MockExecService_Expecter[In, Out] {
	return &MockExecService_Expecter[In, Out]{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockExecService[In, Out]) Exec(ctx context.Context, data In) (Out, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 Out
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, In) (Out, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, In) Out); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(Out)
	}

	if rf, ok := ret.Get(1).(func(context.Context, In) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExecService_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockExecService_Exec_Call[In interface{}, Out interface{}] struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data In
func (_e *MockExecService_Expecter[In, Out]) Exec(ctx interface{}, data interface{}) *MockExecService_Exec_Call[In, Out] {
	return &MockExecService_Exec_Call[In, Out]{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockExecService_Exec_Call[In, Out]) Run(run func(ctx context.Context, data In)) *MockExecService_Exec_Call[In, Out] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(In))
	})
	return _c
}

func (_c *MockExecService_Exec_Call[In, Out]) Return(_a0 Out, _a1 error) *MockExecService_Exec_Call[In, Out] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExecService_Exec_Call[In, Out]) RunAndReturn(run func(context.Context, In) (Out, error)) *MockExecService_Exec_Call[In, Out] {
	_c.Call.Return(run)
	return _c
}

// NewMockExecService creates a new instance of MockExecService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecService[In interface{}, Out interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecService[In, Out] {
	mock := &MockExecService[In, Out]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
