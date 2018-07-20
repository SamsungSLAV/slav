// Code generated by MockGen. DO NOT EDIT.
// Source: git.tizen.org/tools/slav/logger (interfaces: Filter)

package logger

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFilter is a mock of Filter interface
type MockFilter struct {
	ctrl     *gomock.Controller
	recorder *MockFilterMockRecorder
}

// MockFilterMockRecorder is the mock recorder for MockFilter
type MockFilterMockRecorder struct {
	mock *MockFilter
}

// NewMockFilter creates a new mock instance
func NewMockFilter(ctrl *gomock.Controller) *MockFilter {
	mock := &MockFilter{ctrl: ctrl}
	mock.recorder = &MockFilterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFilter) EXPECT() *MockFilterMockRecorder {
	return m.recorder
}

// Verify mocks base method
func (m *MockFilter) Verify(arg0 *Entry) (bool, error) {
	ret := m.ctrl.Call(m, "Verify", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify
func (mr *MockFilterMockRecorder) Verify(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockFilter)(nil).Verify), arg0)
}
