package add_expense

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/add_expense"
)

type usecase interface {
	AddExpense(ctx context.Context, userID int64, req add_expense.AddExpenseReq) (add_expense.AddExpenseResp, error)
}
