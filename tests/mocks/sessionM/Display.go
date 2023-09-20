// Code generated by mockery v2.32.0. DO NOT EDIT.

package sessionM

import mock "github.com/stretchr/testify/mock"

// Display is an autogenerated mock type for the Display type
type Display struct {
	mock.Mock
}

// GetLoginFromSession provides a mock function with given fields: _a0
func (_m *Display) GetLoginFromSession(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSession provides a mock function with given fields: login
func (_m *Display) NewSession(login string) (string, error) {
	ret := _m.Called(login)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSessionLifeTime provides a mock function with given fields: login
func (_m *Display) UpdateSessionLifeTime(login string) {
	_m.Called(login)
}

// NewDisplay creates a new instance of Display. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDisplay(t interface {
	mock.TestingT
	Cleanup(func())
}) *Display {
	mock := &Display{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}