// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package expenses_mocks is a generated GoMock package.
package expenses_mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	decimal "github.com/shopspring/decimal"
	entities "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	limits "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

// Mockcfg is a mock of cfg interface.
type Mockcfg struct {
	ctrl     *gomock.Controller
	recorder *MockcfgMockRecorder
}

// MockcfgMockRecorder is the mock recorder for Mockcfg.
type MockcfgMockRecorder struct {
	mock *Mockcfg
}

// NewMockcfg creates a new mock instance.
func NewMockcfg(ctrl *gomock.Controller) *Mockcfg {
	mock := &Mockcfg{ctrl: ctrl}
	mock.recorder = &MockcfgMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcfg) EXPECT() *MockcfgMockRecorder {
	return m.recorder
}

// BaseCurrency mocks base method.
func (m *Mockcfg) BaseCurrency() entities.Currency {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseCurrency")
	ret0, _ := ret[0].(entities.Currency)
	return ret0
}

// BaseCurrency indicates an expected call of BaseCurrency.
func (mr *MockcfgMockRecorder) BaseCurrency() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseCurrency", reflect.TypeOf((*Mockcfg)(nil).BaseCurrency))
}

// Mocktx is a mock of tx interface.
type Mocktx struct {
	ctrl     *gomock.Controller
	recorder *MocktxMockRecorder
}

// MocktxMockRecorder is the mock recorder for Mocktx.
type MocktxMockRecorder struct {
	mock *Mocktx
}

// NewMocktx creates a new mock instance.
func NewMocktx(ctrl *gomock.Controller) *Mocktx {
	mock := &Mocktx{ctrl: ctrl}
	mock.recorder = &MocktxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktx) EXPECT() *MocktxMockRecorder {
	return m.recorder
}

// WithTransaction mocks base method.
func (m *Mocktx) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MocktxMockRecorder) WithTransaction(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*Mocktx)(nil).WithTransaction), ctx, fn)
}

// MockexpenseStorage is a mock of expenseStorage interface.
type MockexpenseStorage struct {
	ctrl     *gomock.Controller
	recorder *MockexpenseStorageMockRecorder
}

// MockexpenseStorageMockRecorder is the mock recorder for MockexpenseStorage.
type MockexpenseStorageMockRecorder struct {
	mock *MockexpenseStorage
}

// NewMockexpenseStorage creates a new mock instance.
func NewMockexpenseStorage(ctrl *gomock.Controller) *MockexpenseStorage {
	mock := &MockexpenseStorage{ctrl: ctrl}
	mock.recorder = &MockexpenseStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockexpenseStorage) EXPECT() *MockexpenseStorageMockRecorder {
	return m.recorder
}

// AddExpense mocks base method.
func (m *MockexpenseStorage) AddExpense(ctx context.Context, userID int64, exp entities.Expense) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddExpense", ctx, userID, exp)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddExpense indicates an expected call of AddExpense.
func (mr *MockexpenseStorageMockRecorder) AddExpense(ctx, userID, exp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddExpense", reflect.TypeOf((*MockexpenseStorage)(nil).AddExpense), ctx, userID, exp)
}

// GetExpenses mocks base method.
func (m *MockexpenseStorage) GetExpenses(ctx context.Context, userID int64, minTime time.Time) ([]entities.Expense, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpenses", ctx, userID, minTime)
	ret0, _ := ret[0].([]entities.Expense)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpenses indicates an expected call of GetExpenses.
func (mr *MockexpenseStorageMockRecorder) GetExpenses(ctx, userID, minTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpenses", reflect.TypeOf((*MockexpenseStorage)(nil).GetExpenses), ctx, userID, minTime)
}

// MockuserStorage is a mock of userStorage interface.
type MockuserStorage struct {
	ctrl     *gomock.Controller
	recorder *MockuserStorageMockRecorder
}

// MockuserStorageMockRecorder is the mock recorder for MockuserStorage.
type MockuserStorageMockRecorder struct {
	mock *MockuserStorage
}

// NewMockuserStorage creates a new mock instance.
func NewMockuserStorage(ctrl *gomock.Controller) *MockuserStorage {
	mock := &MockuserStorage{ctrl: ctrl}
	mock.recorder = &MockuserStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserStorage) EXPECT() *MockuserStorageMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockuserStorage) Get(ctx context.Context, id int64) (entities.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockuserStorageMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockuserStorage)(nil).Get), ctx, id)
}

// MockcurrencyConverter is a mock of currencyConverter interface.
type MockcurrencyConverter struct {
	ctrl     *gomock.Controller
	recorder *MockcurrencyConverterMockRecorder
}

// MockcurrencyConverterMockRecorder is the mock recorder for MockcurrencyConverter.
type MockcurrencyConverterMockRecorder struct {
	mock *MockcurrencyConverter
}

// NewMockcurrencyConverter creates a new mock instance.
func NewMockcurrencyConverter(ctrl *gomock.Controller) *MockcurrencyConverter {
	mock := &MockcurrencyConverter{ctrl: ctrl}
	mock.recorder = &MockcurrencyConverterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcurrencyConverter) EXPECT() *MockcurrencyConverterMockRecorder {
	return m.recorder
}

// Convert mocks base method.
func (m *MockcurrencyConverter) Convert(ctx context.Context, sum decimal.Decimal, from, to entities.Currency, date time.Time) (decimal.Decimal, entities.Rate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Convert", ctx, sum, from, to, date)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(entities.Rate)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Convert indicates an expected call of Convert.
func (mr *MockcurrencyConverterMockRecorder) Convert(ctx, sum, from, to, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Convert", reflect.TypeOf((*MockcurrencyConverter)(nil).Convert), ctx, sum, from, to, date)
}

// ConvertToBaseCurrency mocks base method.
func (m *MockcurrencyConverter) ConvertToBaseCurrency(ctx context.Context, sum decimal.Decimal, userID int64, date time.Time) (decimal.Decimal, entities.Rate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertToBaseCurrency", ctx, sum, userID, date)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(entities.Rate)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ConvertToBaseCurrency indicates an expected call of ConvertToBaseCurrency.
func (mr *MockcurrencyConverterMockRecorder) ConvertToBaseCurrency(ctx, sum, userID, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertToBaseCurrency", reflect.TypeOf((*MockcurrencyConverter)(nil).ConvertToBaseCurrency), ctx, sum, userID, date)
}

// MocklimitUc is a mock of limitUc interface.
type MocklimitUc struct {
	ctrl     *gomock.Controller
	recorder *MocklimitUcMockRecorder
}

// MocklimitUcMockRecorder is the mock recorder for MocklimitUc.
type MocklimitUcMockRecorder struct {
	mock *MocklimitUc
}

// NewMocklimitUc creates a new mock instance.
func NewMocklimitUc(ctrl *gomock.Controller) *MocklimitUc {
	mock := &MocklimitUc{ctrl: ctrl}
	mock.recorder = &MocklimitUcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklimitUc) EXPECT() *MocklimitUcMockRecorder {
	return m.recorder
}

// CheckLimit mocks base method.
func (m *MocklimitUc) CheckLimit(ctx context.Context, userID int64, expense entities.Expense) (limits.LimitCheckResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLimit", ctx, userID, expense)
	ret0, _ := ret[0].(limits.LimitCheckResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLimit indicates an expected call of CheckLimit.
func (mr *MocklimitUcMockRecorder) CheckLimit(ctx, userID, expense interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLimit", reflect.TypeOf((*MocklimitUc)(nil).CheckLimit), ctx, userID, expense)
}
