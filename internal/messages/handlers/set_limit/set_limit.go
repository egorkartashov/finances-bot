package set_limit

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_limit"
	"go.uber.org/zap"
)

type SetLimit struct {
	usecase       usecase
	messageSender messages.MessageSender
}

func New(limitUc usecase, messageSender messages.MessageSender) *SetLimit {
	return &SetLimit{
		usecase:       limitUc,
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
	req, parseErr := h.parseReq(params, msg.UserID, time.Now())
	if parseErr != nil {
		badFormatMessage := parseErr.Error() + "\n" + AddLimitHelp
		err := h.messageSender.SendText(badFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "send parsing error to user failed"))
	}

	resp, err := h.usecase.SetLimit(ctx, req)
	if err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "usecase failed"))
	}

	response := constructResponseMessage(resp)
	if err = h.messageSender.SendText(response, msg.UserID); err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "send success message failed"))
	}

	return utils.HandleSuccess
}

func (h *SetLimit) parseReq(params string, userID int64, date time.Time) (req set_limit.Req, err error) {
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
		logger.Error(errMsg, zap.Error(convErr))
		return
	}

	req = set_limit.Req{
		UserID:            userID,
		Category:          category,
		SumInUserCurrency: decimal.NewFromInt(sum),
		SetAt:             date,
	}
	return
}

func constructResponseMessage(resp set_limit.Resp) string {
	limitSumInfo := fmt.Sprintf("%v %s", resp.InputSum.Sum, resp.InputSum.Currency)
	if resp.SavedSum.Currency != resp.InputSum.Currency {
		limitSumInfo += fmt.Sprintf(" (%v %s)", resp.SavedSum.Sum, resp.SavedSum.Currency)
	}
	response := fmt.Sprintf("Успешно задали лимит %s по категории \"%s\"", limitSumInfo, resp.Category)
	return response
}
