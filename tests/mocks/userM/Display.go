// Code generated by mockery v2.32.0. DO NOT EDIT.

package userM

import mock "github.com/stretchr/testify/mock"

// Display is an autogenerated mock type for the Display type
type Display struct {
	mock.Mock
}

// AuthenticationUser provides a mock function with given fields: email, password
func (_m *Display) AuthenticationUser(email string, password string) (string, error) {
	ret := _m.Called(email, password)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangePassword provides a mock function with given fields: email, password, newPassword
func (_m *Display) ChangePassword(email string, password string, newPassword string) error {
	ret := _m.Called(email, password, newPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(email, password, newPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Display) Close() {
	_m.Called()
}

// DeleteUser provides a mock function with given fields: id, email, password
func (_m *Display) DeleteUser(id string, email string, password string) error {
	ret := _m.Called(id, email, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(id, email, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsUser provides a mock function with given fields: email
func (_m *Display) IsUser(email string) bool {
	ret := _m.Called(email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewUser provides a mock function with given fields: email, password
func (_m *Display) NewUser(email string, password string) (string, error) {
	ret := _m.Called(email, password)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
