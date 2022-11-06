package add_expense_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/dtos"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/add_expense"
	add_expense_mocks "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/add_expense/mocks"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_currency"
)

const baseCurr = currency.RUB

type inputData struct {
	ctx    context.Context
	userID int64
	req    add_expense.AddExpenseReq
}

type deps struct {
	expenseStorage    *add_expense_mocks.MockexpenseStorage
	userStorage       *add_expense_mocks.MockuserStorage
	currencyConverter *add_expense_mocks.MockcurrencyConverter
	limitUc           *add_expense_mocks.MocklimitUc
}

func TestUsecase_AddExpense_WhenErrGettingUser_ReturnsErr(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			req: add_expense.AddExpenseReq{
				Category: "продукты",
				Sum:      decimal.NewFromInt32(100),
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			req: add_expense.AddExpenseReq{
				Category: "1232313131",
				Sum:      decimal.NewFromInt32(-371837931),
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, tc,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID

						getUserErr := errors.New("getUserErr")
						deps.userStorage.EXPECT().
							Get(ctx, userID).
							Return(entities.User{}, false, getUserErr)

						wantRes = add_expense.AddExpenseResp{}
						wantErr = getUserErr

						return
					},
				)
			},
		)
	}
}

func TestUsecase_AddExpense_WhenUserDoesNotExist_ReturnsErrUserNotFound(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			req: add_expense.AddExpenseReq{
				Category: "продукты",
				Sum:      decimal.NewFromInt32(100),
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			req: add_expense.AddExpenseReq{
				Category: "1232313131",
				Sum:      decimal.NewFromInt32(-371837931),
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, tc,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID

						deps.userStorage.EXPECT().
							Get(ctx, userID).
							Return(entities.User{}, false, nil)

						wantRes = add_expense.AddExpenseResp{}
						wantErr = set_currency.NewUserNotFoundErr(userID)

						return
					},
				)
			},
		)
	}
}

func TestUsecase_AddExpense_WhenErrConvertingToBaseCurrency_ReturnsErr(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			req: add_expense.AddExpenseReq{
				Category: "продукты",
				Sum:      decimal.NewFromInt32(100),
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			req: add_expense.AddExpenseReq{
				Category: "1232313131",
				Sum:      decimal.NewFromInt32(-371837931),
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, tc,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						req := inputData.req

						user := entities.User{
							ID:       userID,
							Currency: currency.EUR,
						}
						deps.userStorage.EXPECT().Get(ctx, userID).Return(user, true, nil)

						convertToBaseCurrencyErr := errors.New("convertToBaseCurrencyErr")
						deps.currencyConverter.EXPECT().
							ToBase(ctx, user.Currency, req.Sum, req.Date).
							Return(decimal.Zero, entities.Currency(""), convertToBaseCurrencyErr)

						wantRes = add_expense.AddExpenseResp{}
						wantErr = convertToBaseCurrencyErr

						return
					},
				)
			},
		)
	}
}

func TestUsecase_AddExpense_WhenErrInCheckLimit_ReturnsErr(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			req: add_expense.AddExpenseReq{
				Category: "продукты",
				Sum:      decimal.NewFromInt32(100),
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			req: add_expense.AddExpenseReq{
				Category: "1232313131",
				Sum:      decimal.NewFromInt32(-371837931),
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		input := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, input,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						req := inputData.req

						user := entities.User{
							ID:       userID,
							Currency: currency.EUR,
						}
						deps.userStorage.EXPECT().Get(ctx, userID).Return(user, true, nil)

						sumInBaseCurr := decimal.NewFromInt32(rand.Int31())
						deps.currencyConverter.EXPECT().
							ToBase(ctx, user.Currency, req.Sum, req.Date).
							Return(sumInBaseCurr, baseCurr, nil)

						expense := entities.Expense{
							Category: req.Category,
							Sum:      sumInBaseCurr,
							Date:     req.Date,
						}
						limitCheckErr := errors.New("limitCheckErr")
						deps.limitUc.EXPECT().
							Check(ctx, userID, expense).
							Return(limits.LimitCheckResult{}, limitCheckErr)

						wantRes = add_expense.AddExpenseResp{}
						wantErr = limitCheckErr

						return
					},
				)
			},
		)
	}
}

func TestUsecase_AddExpense_WhenErrSavingExpense_ReturnsErr(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			req: add_expense.AddExpenseReq{
				Category: "продукты",
				Sum:      decimal.NewFromInt32(100),
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			req: add_expense.AddExpenseReq{
				Category: "1232313131",
				Sum:      decimal.NewFromInt32(-371837931),
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		input := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, input,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						req := inputData.req

						user := entities.User{
							ID:       userID,
							Currency: currency.EUR,
						}
						deps.userStorage.EXPECT().Get(ctx, userID).Return(user, true, nil)

						sumInBaseCurr := decimal.NewFromInt32(rand.Int31())
						deps.currencyConverter.EXPECT().
							ToBase(ctx, user.Currency, req.Sum, req.Date).
							Return(sumInBaseCurr, baseCurr, nil)

						expense := entities.Expense{
							Category: req.Category,
							Sum:      sumInBaseCurr,
							Date:     req.Date,
						}
						deps.limitUc.EXPECT().
							Check(ctx, userID, expense).
							Return(limits.LimitCheckResult{}, nil)

						addExpenseErr := errors.New("addExpenseErr")
						deps.expenseStorage.EXPECT().
							AddExpense(ctx, userID, expense).
							Return(addExpenseErr)

						wantRes = add_expense.AddExpenseResp{}
						wantErr = addExpenseErr

						return
					},
				)
			},
		)
	}
}

func TestUsecase_AddExpense_WhenNoErr_ReturnsCorrectResult(t *testing.T) {
	tests := []struct {
		inputData        inputData
		limitCheckResult limits.LimitCheckResult
	}{
		{
			inputData: inputData{
				userID: 1,
				req: add_expense.AddExpenseReq{
					Category: "продукты",
					Sum:      decimal.NewFromInt32(100),
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:            limits.StatusLimitExceeded,
				SumWithNewExpense: decimal.New(12345, -2),
				Limit:             decimal.NewFromInt32(100),
				Currency:          baseCurr,
			},
		},
		{
			inputData: inputData{
				userID: 1,
				req: add_expense.AddExpenseReq{
					Category: "машина",
					Sum:      decimal.NewFromInt32(25),
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:            limits.StatusLimitSatisfied,
				SumWithNewExpense: decimal.New(10, 0),
				Limit:             decimal.NewFromInt32(100),
				Currency:          baseCurr,
			},
		},
		{
			inputData: inputData{
				userID: 1,
				req: add_expense.AddExpenseReq{
					Category: "учеба",
					Sum:      decimal.NewFromInt32(99),
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:            limits.StatusLimitNotSet,
				SumWithNewExpense: decimal.New(234, 0),
				Limit:             decimal.Zero,
				Currency:          baseCurr,
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testAddExpenseWithThisArrange(
					t, tc.inputData,
					func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						req := inputData.req

						user := entities.User{
							ID:       userID,
							Currency: currency.EUR,
						}
						deps.userStorage.EXPECT().Get(ctx, userID).Return(user, true, nil)

						sumInBaseCurr := decimal.NewFromInt32(rand.Int31())
						deps.currencyConverter.EXPECT().
							ToBase(ctx, user.Currency, req.Sum, req.Date).
							Return(sumInBaseCurr, baseCurr, nil)

						expense := entities.Expense{
							Category: req.Category,
							Sum:      sumInBaseCurr,
							Date:     req.Date,
						}
						limitCheckResult := tc.limitCheckResult
						deps.limitUc.EXPECT().
							Check(ctx, userID, expense).
							Return(limitCheckResult, nil)

						deps.expenseStorage.EXPECT().
							AddExpense(ctx, userID, expense).
							Return(nil)

						wantRes = add_expense.AddExpenseResp{
							UserInputSum: dtos.SumWithCurrency{
								Sum:      req.Sum,
								Currency: user.Currency,
							},
							SavedSum: dtos.SumWithCurrency{
								Sum:      expense.Sum,
								Currency: baseCurr,
							},
							Date:             req.Date,
							Category:         req.Category,
							LimitCheckResult: limitCheckResult,
						}
						wantErr = nil

						return
					},
				)
			},
		)
	}
}

func testAddExpenseWithThisArrange(
	t *testing.T,
	inputData inputData,
	arrange func(inputData inputData, deps deps) (wantRes add_expense.AddExpenseResp, wantErr error),
) {
	ctrl := gomock.NewController(t)
	tx := getTxMock(ctrl)
	expenseStorage := add_expense_mocks.NewMockexpenseStorage(ctrl)
	userStorage := add_expense_mocks.NewMockuserStorage(ctrl)
	currencyConverter := add_expense_mocks.NewMockcurrencyConverter(ctrl)
	limitUc := add_expense_mocks.NewMocklimitUc(ctrl)

	deps := deps{expenseStorage, userStorage, currencyConverter, limitUc}
	wantRes, wantErr := arrange(inputData, deps)

	expensesModel := add_expense.NewUsecase(tx, expenseStorage, userStorage, currencyConverter, limitUc)
	gotRes, gotErr := expensesModel.AddExpense(inputData.ctx, inputData.userID, inputData.req)

	assert.Equal(t, wantRes, gotRes)
	assert.Equal(t, wantErr, gotErr)
}

func getTxMock(ctrl *gomock.Controller) *add_expense_mocks.Mocktx {
	txMock := add_expense_mocks.NewMocktx(ctrl)
	txMock.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		)
	return txMock
}
