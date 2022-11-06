package add_expense

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/add_expense"
)

type AddExpense struct {
	uc     usecase
	sender messages.MessageSender
}

func New(uc usecase, sender messages.MessageSender) *AddExpense {
	return &AddExpense{
		uc:     uc,
		sender: sender,
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
	req, err := parseReq(expenseParams)
	if err != nil {
		err := h.sender.SendText(incorrectFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}

	res, err := h.uc.AddExpense(ctx, msg.UserID, *req)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	response := constructResponse(res)
	err = h.sender.SendText(response, msg.UserID)

	return utils.HandleWithErrorOrNil(err)
}

func parseReq(paramsStr string) (*add_expense.AddExpenseReq, error) {
	params := strings.Split(paramsStr, " ")
	if len(params) < 2 || len(params) > 3 {
		return nil, incorrectParamsCountErr
	}

	category := params[0]
	sum, err := decimal.NewFromString(params[1])
	if err != nil {
		return nil, incorrectSumErr
	}

	var date time.Time
	if date, err = parseDate(params); err != nil {
		return nil, incorrectDateErr
	}

	return &add_expense.AddExpenseReq{
		Category: category,
		Sum:      sum,
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

func constructResponse(res add_expense.AddExpenseResp) string {
	sumInfo := fmt.Sprintf("%v %s", res.UserInputSum.Sum, res.UserInputSum.Currency)
	if res.SavedSum.Currency != res.UserInputSum.Currency {
		sumInfo += fmt.Sprintf(" (%v %s)", res.SavedSum.Sum, res.SavedSum.Currency)
	}

	dateStr := res.Date.Format("02.01.2006")
	expenseAddedMsg := fmt.Sprintf(
		"Успешно добавили трату: категория \"%s\", сумма %s, дата %s", res.Category, sumInfo, dateStr,
	)

	limitCheckResultStr := limitCheckResultToString(res.LimitCheckResult)
	if limitCheckResultStr == "" {
		return expenseAddedMsg
	}
	return limitCheckResultStr + "\n\n" + expenseAddedMsg
}

func limitCheckResultToString(result limits.LimitCheckResult) string {
	if result.Status != limits.StatusLimitExceeded {
		return ""
	}

	limitStatusString := fmt.Sprintf(
		"лимит %v %[2]s, потрачено %v %[2]s", result.Limit, result.Currency, result.SumWithNewExpense,
	)
	return fmt.Sprintf(
		"Превышен лимит по тратам в данной категории за месяц (%s), будьте аккуратнее со своими тратами!",
		limitStatusString,
	)
}

func (h *AddExpense) Name() string {
	return "AddExpense"
}
