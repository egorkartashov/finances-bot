package expenses_test

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
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	expenses_mocks "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses/mocks"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type inputData struct {
	ctx    context.Context
	userID int64
	dto    expenses.AddExpenseDto
}

type deps struct {
	cfg               *expenses_mocks.Mockcfg
	expenseStorage    *expenses_mocks.MockexpenseStorage
	currencyConverter *expenses_mocks.MockcurrencyConverter
	limitUc           *expenses_mocks.MocklimitUc
}

func TestUsecase_AddExpense_WhenErrConvertingToBaseCurrency_ReturnsErr(t *testing.T) {
	tests := []inputData{
		{
			userID: 1,
			dto: expenses.AddExpenseDto{
				Category: "продукты",
				Sum:      100,
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			dto: expenses.AddExpenseDto{
				Category: "1232313131",
				Sum:      -371837931,
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testWithThisArrange(
					t, tc,
					func(inputData inputData, deps deps) (wantRes expenses.AddExpenseResult, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						dto := inputData.dto

						convertToBaseCurrencyErr := errors.New("convertToBaseCurrencyErr")

						deps.currencyConverter.EXPECT().
							ConvertToBaseCurrency(ctx, decimal.NewFromInt32(dto.Sum), userID, dto.Date).
							Return(decimal.Zero, entities.Rate{}, convertToBaseCurrencyErr)

						wantRes = expenses.AddExpenseResult{}
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
			dto: expenses.AddExpenseDto{
				Category: "продукты",
				Sum:      100,
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			dto: expenses.AddExpenseDto{
				Category: "1232313131",
				Sum:      -371837931,
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		input := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testWithThisArrange(
					t, input,
					func(inputData inputData, deps deps) (wantRes expenses.AddExpenseResult, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						dto := inputData.dto

						sumDecimal := decimal.NewFromInt32(dto.Sum)
						rate := entities.Rate{
							Date:  dto.Date,
							From:  currency.RUB,
							To:    currency.RUB,
							Value: decimal.NewFromInt32(1),
						}
						deps.currencyConverter.EXPECT().
							ConvertToBaseCurrency(gomock.Any(), sumDecimal, userID, dto.Date).
							Return(sumDecimal, rate, nil)

						expense := entities.Expense{
							Category: dto.Category,
							Sum:      sumDecimal,
							Date:     dto.Date,
						}
						limitCheckErr := errors.New("limitCheckErr")
						deps.limitUc.EXPECT().
							CheckLimit(ctx, userID, expense).
							Return(limits.LimitCheckResult{}, limitCheckErr)

						wantRes = expenses.AddExpenseResult{}
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
			dto: expenses.AddExpenseDto{
				Category: "продукты",
				Sum:      100,
				Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			userID: 101010,
			dto: expenses.AddExpenseDto{
				Category: "1232313131",
				Sum:      -371837931,
				Date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		input := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testWithThisArrange(
					t, input,
					func(inputData inputData, deps deps) (wantRes expenses.AddExpenseResult, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						dto := inputData.dto

						sumDecimal := decimal.NewFromInt32(dto.Sum)
						rate := entities.Rate{
							Date:  dto.Date,
							From:  currency.RUB,
							To:    currency.RUB,
							Value: decimal.NewFromInt32(1),
						}
						deps.currencyConverter.EXPECT().
							ConvertToBaseCurrency(gomock.Any(), sumDecimal, userID, dto.Date).
							Return(sumDecimal, rate, nil)

						expense := entities.Expense{
							Category: dto.Category,
							Sum:      sumDecimal,
							Date:     dto.Date,
						}
						deps.limitUc.EXPECT().
							CheckLimit(ctx, userID, expense).
							Return(limits.LimitCheckResult{}, nil)

						addExpenseErr := errors.New("addExpenseErr")
						deps.expenseStorage.EXPECT().
							AddExpense(ctx, userID, expense).
							Return(addExpenseErr)

						wantRes = expenses.AddExpenseResult{}
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
		convertedSum     decimal.Decimal
		rate             entities.Rate
		limitCheckResult limits.LimitCheckResult
	}{
		{
			inputData: inputData{
				userID: 1,
				dto: expenses.AddExpenseDto{
					Category: "продукты",
					Sum:      100,
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			convertedSum: decimal.NewFromInt32(rand.Int31()),
			rate: entities.Rate{
				From:  currency.EUR,
				To:    currency.RUB,
				Value: decimal.NewFromFloat32(rand.Float32()),
				Date:  time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:          limits.StatusLimitExceeded,
				CurrentTotalSum: decimal.New(12345, -2),
				Limit:           decimal.NewFromInt32(100),
			},
		},
		{
			inputData: inputData{
				userID: 1,
				dto: expenses.AddExpenseDto{
					Category: "машина",
					Sum:      25,
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			convertedSum: decimal.NewFromInt32(rand.Int31()),
			rate: entities.Rate{
				From:  currency.USD,
				To:    currency.RUB,
				Value: decimal.NewFromFloat32(rand.Float32()),
				Date:  time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:          limits.StatusLimitSatisfied,
				CurrentTotalSum: decimal.New(10, 0),
				Limit:           decimal.NewFromInt32(100),
			},
		},
		{
			inputData: inputData{
				userID: 1,
				dto: expenses.AddExpenseDto{
					Category: "учеба",
					Sum:      99,
					Date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			convertedSum: decimal.NewFromInt32(rand.Int31()),
			rate: entities.Rate{
				From:  currency.RUB,
				To:    currency.EUR,
				Value: decimal.NewFromFloat32(rand.Float32()),
				Date:  time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			limitCheckResult: limits.LimitCheckResult{
				Status:          limits.StatusLimitNotSet,
				CurrentTotalSum: decimal.New(234, 0),
				Limit:           decimal.Zero,
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(
			"", func(t *testing.T) {
				t.Parallel()

				testWithThisArrange(
					t, tc.inputData,
					func(inputData inputData, deps deps) (wantRes expenses.AddExpenseResult, wantErr error) {
						ctx := inputData.ctx
						userID := inputData.userID
						dto := inputData.dto

						deps.currencyConverter.EXPECT().
							ConvertToBaseCurrency(gomock.Any(), decimal.NewFromInt32(dto.Sum), userID, dto.Date).
							Return(tc.convertedSum, tc.rate, nil)

						expense := entities.Expense{
							Category: dto.Category,
							Sum:      tc.convertedSum,
							Date:     dto.Date,
						}
						limitCheckResult := tc.limitCheckResult
						deps.limitUc.EXPECT().
							CheckLimit(ctx, userID, expense).
							Return(limitCheckResult, nil)

						deps.expenseStorage.EXPECT().
							AddExpense(ctx, userID, expense).
							Return(nil)

						wantRes = expenses.AddExpenseResult{
							SumInUserCurrency: dto.Sum,
							Expense:           expense,
							Rate:              tc.rate,
							LimitCheckResult:  limitCheckResult,
						}
						wantErr = nil

						return
					},
				)
			},
		)
	}
}

func testWithThisArrange(
	t *testing.T,
	inputData inputData,
	arrange func(inputData inputData, deps deps) (wantRes expenses.AddExpenseResult, wantErr error),
) {
	ctrl := gomock.NewController(t)
	cfg := expenses_mocks.NewMockcfg(ctrl)
	tx := getTxMock(ctrl)
	expenseStorage := expenses_mocks.NewMockexpenseStorage(ctrl)
	currencyConverter := expenses_mocks.NewMockcurrencyConverter(ctrl)
	limitUc := expenses_mocks.NewMocklimitUc(ctrl)

	deps := deps{cfg, expenseStorage, currencyConverter, limitUc}
	wantRes, wantErr := arrange(inputData, deps)

	expensesModel := expenses.NewUsecase(cfg, tx, expenseStorage, nil, currencyConverter, limitUc)
	gotRes, gotErr := expensesModel.AddExpense(inputData.ctx, inputData.userID, inputData.dto)

	assert.Equal(t, wantRes, gotRes)
	assert.Equal(t, wantErr, gotErr)
}

func getTxMock(ctrl *gomock.Controller) *expenses_mocks.Mocktx {
	txMock := expenses_mocks.NewMocktx(ctrl)
	txMock.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		)
	return txMock
}
