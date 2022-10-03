// Code generated by MockGen. DO NOT EDIT.
// Source: model.go

// Package expenses_mocks is a generated GoMock package.
package expenses_mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	expenses "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddExpense mocks base method.
func (m *MockStorage) AddExpense(userID int64, exp expenses.Expense) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddExpense", userID, exp)
}

// AddExpense indicates an expected call of AddExpense.
func (mr *MockStorageMockRecorder) AddExpense(userID, exp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddExpense", reflect.TypeOf((*MockStorage)(nil).AddExpense), userID, exp)
}

// GetCategoriesTotals mocks base method.
func (m *MockStorage) GetCategoriesTotals(userID int64, minTime time.Time) []expenses.CategoryTotal {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoriesTotals", userID, minTime)
	ret0, _ := ret[0].([]expenses.CategoryTotal)
	return ret0
}

// GetCategoriesTotals indicates an expected call of GetCategoriesTotals.
func (mr *MockStorageMockRecorder) GetCategoriesTotals(userID, minTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoriesTotals", reflect.TypeOf((*MockStorage)(nil).GetCategoriesTotals), userID, minTime)
}
