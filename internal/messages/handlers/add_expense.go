package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type AddExpense struct {
	expensesUc *expenses.Usecase
	base
}

func NewAddExpense(expensesUc *expenses.Usecase, sender messages.MessageSender) *AddExpense {
	return &AddExpense{
		expensesUc: expensesUc,
		base: base{
			MessageSender: sender,
		},
	}
}

const (
	expenseKeyword   = "трата"
	AddExpenseFormat = "трата <категория> <сумма> [<дата ДД.ММ.ГГГГ>]. \n" +
		"Например: трата продукты 310 15.01.2022. Если дата не указана, то трата сохранится за сегодняшний день"
	AddExpenseHelp         = "Чтобы добавить трату, введи ее в следующем формате: \n" + AddExpenseFormat
	incorrectFormatMessage = "Трата введена в некорректном формате. " +
		"Правильный формат: " + AddExpenseFormat
)

var (
	incorrectParamsCountErr = errors.New("expense: incorrect params count")
	incorrectSumErr         = errors.New("expense: failed to parse sum")
	incorrectDateErr        = errors.New("expense: failed to parse date")
)

func (h *AddExpense) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.Text, expenseKeyword) {
		return utils.HandleSkipped
	}

	expenseParams := strings.TrimPrefix(msg.Text, expenseKeyword)
	expenseParams = strings.Trim(expenseParams, " ")
	exp, err := parseExpense(expenseParams)
	if err != nil {
		err := h.MessageSender.SendText(incorrectFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}

	res, err := h.expensesUc.AddExpense(ctx, msg.UserID, *exp)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	response := constructResponse(res)
	err = h.MessageSender.SendText(response, msg.UserID)

	return utils.HandleWithErrorOrNil(err)
}

func parseExpense(paramsStr string) (*expenses.AddExpenseDto, error) {
	params := strings.Split(paramsStr, " ")
	if len(params) < 2 || len(params) > 3 {
		return nil, incorrectParamsCountErr
	}

	category := params[0]
	sum, err := strconv.ParseInt(params[1], 10, 32)
	if err != nil {
		return nil, incorrectSumErr
	}

	var date time.Time
	if date, err = parseDate(params); err != nil {
		return nil, incorrectDateErr
	}

	return &expenses.AddExpenseDto{
		Category: category,
		Sum:      int32(sum),
		Date:     date,
	}, nil
}

func parseDate(params []string) (time.Time, error) {
	if len(params) >= 3 {
		date, err := time.ParseInLocation("02.01.2006", params[2], time.UTC)
		return date, err
	} else {
		year, month, day := time.Now().UTC().Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	}
}

func constructResponse(res expenses.AddExpenseResult) string {
	baseCurrency := res.Rate.To
	sumInfo := fmt.Sprintf("%v %s", res.SumInUserCurrency, res.Rate.From)
	if res.Rate.From != baseCurrency {
		sumInfo += fmt.Sprintf(" (%v %s)", res.Expense.Sum, baseCurrency)
	}

	dateStr := res.Expense.Date.Format("02.01.2006")
	expenseAddedMsg := fmt.Sprintf(
		"Успешно добавили трату: категория \"%s\", сумма %s, дата %s", res.Expense.Category, sumInfo, dateStr,
	)

	limitRes := res.LimitCheckResult
	if limitRes.Status != limits.StatusLimitExceeded {
		return expenseAddedMsg
	}

	totalSumWithNewExpense := res.LimitCheckResult.CurrentTotalSum.Add(res.Expense.Sum)
	limitExceededMsg := fmt.Sprintf(
		"Превышен лимит по тратам в данной категории за месяц (лимит %v %s, потрачено %v %s), "+
			"будьте аккуратнее со своими тратами!", limitRes.Limit, baseCurrency, totalSumWithNewExpense,
		baseCurrency,
	)
	return limitExceededMsg + "\n\n" + expenseAddedMsg
}

func (h *AddExpense) Name() string {
	return "AddExpense"
}
