package set_limit

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type SetLimit struct {
	limitUc       limitUc
	messageSender messages.MessageSender
}

func NewSetLimit(limitUc limitUc, messageSender messages.MessageSender) *SetLimit {
	return &SetLimit{
		limitUc:       limitUc,
		messageSender: messageSender,
	}
}

func (h *SetLimit) Name() string {
	return "SetLimit"
}

const (
	addLimitKeyword = "задать лимит"
	AddLimitHelp    = "Чтобы задать лимит по категории на месяц, введи:\n" +
		"\"" + addLimitKeyword + " <категория> <сумма>\""
)

func (h *SetLimit) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.Text, addLimitKeyword) {
		return utils.HandleSkipped
	}

	params := strings.Trim(strings.TrimPrefix(msg.Text, addLimitKeyword+" "), " ")
	limit, parseErr := h.parseLimit(params, msg.UserID, time.Now())
	if parseErr != nil {
		badFormatMessage := parseErr.Error() + "\n" + AddLimitHelp
		err := h.messageSender.SendText(badFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "send parsing error to user failed"))
	}

	res, err := h.limitUc.SetLimit(ctx, limit)
	if err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "usecase failed"))
	}

	response := constructResponseMessage(res)
	if err = h.messageSender.SendText(response, msg.UserID); err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "send success message failed"))
	}

	return utils.HandleSuccess
}

func (h *SetLimit) parseLimit(params string, userID int64, date time.Time) (
	limit entities.MonthBudgetLimit, err error,
) {
	split := strings.Split(params, " ")
	if len(split) != 2 {
		err = errors.New("Некорректное число параметров")
		return
	}

	category := split[0]
	var sum int64
	sum, convErr := strconv.ParseInt(split[1], 10, 32)
	if convErr != nil {
		errMsg := "Сумма лимита не является корректным целым числом"
		err = errors.New(errMsg)
		logger.Errorf(errMsg+": %s", convErr)
		return
	}

	limit = entities.MonthBudgetLimit{
		UserID:   userID,
		Category: category,
		Sum:      decimal.NewFromInt(sum),
		SetAt:    date,
	}
	return
}

func constructResponseMessage(res limits.SetLimitResult) string {
	limitSumInfo := fmt.Sprintf("%v %s", res.SumInUserCurrency, res.ExchangeRate.From)
	if res.ExchangeRate.From != res.ExchangeRate.To {
		limitSumInfo += fmt.Sprintf(" (%v %s)", res.Limit.Sum, res.ExchangeRate.To)
	}
	response := fmt.Sprintf("Успешно задали лимит %s по категории \"%s\"", limitSumInfo, res.Limit.Category)
	return response
}
