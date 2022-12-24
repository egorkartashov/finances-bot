package add_expense

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/dtos"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/users"
)

type Usecase struct {
	tx                tx
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userStorage       userStorage
	limitChecker      limitChecker
	reportCache       reportCache
}

type AddExpenseReq struct {
	Category string
	Sum      decimal.Decimal
	Date     time.Time
}

type AddExpenseResp struct {
	UserInputSum     dtos.SumWithCurrency
	SavedSum         dtos.SumWithCurrency
	Category         string
	Date             time.Time
	LimitCheckResult limits.LimitCheckResult
}

func (u *Usecase) AddExpense(ctx context.Context, userID int64, req AddExpenseReq) (AddExpenseResp, error) {
	user, ok, err := u.userStorage.Get(ctx, userID)
	if err != nil {
		return AddExpenseResp{}, err
	}
	if !ok {
		return AddExpenseResp{}, users.NewUserNotFoundErr(userID)
	}

	expense, baseCurr, err := u.convertToExpenseInBaseCurrency(ctx, req, user.Currency)
	if err != nil {
		return AddExpenseResp{}, err
	}

	err = u.reportCache.DeleteAffected(ctx, userID, req.Date)
	if err != nil {
		return AddExpenseResp{}, errors.WithMessage(err, "failed to delete affected cached reports")
	}

	var res AddExpenseResp
	txErr := u.tx.WithTransaction(
		ctx, func(txCtx context.Context) error {
			limitCheckResult, err := u.limitChecker.Check(txCtx, userID, expense)
			if err != nil {
				return err
			}

			if err = u.expenseStorage.AddExpense(txCtx, userID, expense); err != nil {
				return err
			}

			res = AddExpenseResp{
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
			return nil
		},
	)

	return res, txErr
}

func (u *Usecase) convertToExpenseInBaseCurrency(
	ctx context.Context, req AddExpenseReq, userCurr entities.Currency,
) (exp entities.Expense, baseCurr entities.Currency, err error) {
	subRub, baseCurr, err := u.currencyConverter.ToBase(ctx, userCurr, req.Sum, req.Date)
	if err != nil {
		return entities.Expense{}, "", err
	}

	exp = entities.Expense{
		Category: req.Category,
		Sum:      subRub,
		Date:     req.Date,
	}
	return exp, baseCurr, nil
}

func NewUsecase(
	tx tx, es expenseStorage, us userStorage, cc currencyConverter, lc limitChecker, rc reportCache,
) *Usecase {
	return &Usecase{
		tx:                tx,
		expenseStorage:    es,
		userStorage:       us,
		currencyConverter: cc,
		limitChecker:      lc,
		reportCache:       rc,
	}
}
