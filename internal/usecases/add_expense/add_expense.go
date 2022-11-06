package add_expense

import (
	"context"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/dtos"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
)

type Usecase struct {
	tx                tx
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userStorage       userStorage
	limitChecker      limitUc
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

func (u *Usecase) AddExpense(outsideCtx context.Context, userID int64, req AddExpenseReq) (AddExpenseResp, error) {
	var res AddExpenseResp
	txErr := u.tx.WithTransaction(
		outsideCtx, func(ctx context.Context) error {
			user, ok, err := u.userStorage.Get(ctx, userID)
			if err != nil {
				return err
			}
			if !ok {
				return users.NewUserNotFoundErr(userID)
			}

			expense, baseCurr, err := u.convertToExpenseInBaseCurrency(ctx, req, user.Currency)
			if err != nil {
				return err
			}

			// Пока просто сохраняем статус проверки лимита, т.к. запрещать вводить трату не будем.
			// Если лимит превышен, просто сообщим об этом пользователю
			limitCheckResult, err := u.limitChecker.Check(ctx, userID, expense)
			if err != nil {
				return err
			}

			if err = u.expenseStorage.AddExpense(ctx, userID, expense); err != nil {
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

func NewUsecase(tx tx, es expenseStorage, us userStorage, cc currencyConverter, lu limitUc) *Usecase {
	return &Usecase{
		tx:                tx,
		expenseStorage:    es,
		userStorage:       us,
		currencyConverter: cc,
		limitChecker:      lu,
	}
}
