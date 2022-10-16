package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
)

type AddExpense struct {
	expensesUc *expenses.Usecase
	base
}

func NewAddExpense(expensesUc *expenses.Usecase, sender messages.MessageSender) *AddExpense {
	return &AddExpense{
		expensesUc: expensesUc,
		base: base{
			messageSender: sender,
		},
	}
}

const (
	expenseKeyword = "трата"
	ExpenseFormat  = "трата <категория> <сумма в руб.> [<дата ДД.ММ.ГГГГ>]. \n" +
		"Например: трата продукты 310 15.01.2022. Если дата не указана, то трата сохранится за сегодняшний день"
	incorrectFormatMessage = "Трата введена в некорректном формате. " +
		"Правильный формат: " + ExpenseFormat
)

var (
	incorrectParamsCountErr = errors.New("expense: incorrect params count")
	incorrectSumErr         = errors.New("expense: failed to parse sum")
	incorrectDateErr        = errors.New("expense: failed to parse date")
)

func (h *AddExpense) Handle(msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.Text, expenseKeyword) {
		return handleSkipped
	}

	expenseParams := strings.TrimPrefix(msg.Text, expenseKeyword)
	expenseParams = strings.Trim(expenseParams, " ")
	exp, err := parseExpense(expenseParams)
	if err != nil {
		err := h.messageSender.SendText(incorrectFormatMessage, msg.UserID)
		return handleWithErrorOrNil(err)
	}

	res, err := h.expensesUc.AddExpense(msg.UserID, *exp)
	if err != nil {
		return handleWithErrorOrNil(err)
	}

	dateStr := res.Date.Format("02.01.2006")
	successMsg := fmt.Sprintf("Успешно добавили трату: категория \"%s\", сумма %v %s, дата %s",
		res.Category, res.Sum, res.Cur, dateStr)

	err = h.messageSender.SendText(successMsg, msg.UserID)

	return handleWithErrorOrNil(err)
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

func (h *AddExpense) Name() string {
	return "AddExpense"
}
