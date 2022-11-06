package remove_limit

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type RemoveLimit struct {
	usecase       usecase
	messageSender messages.MessageSender
}

func (h *RemoveLimit) Name() string {
	return "RemoveLimit"
}

const (
	removeLimitKeyword = "убрать лимит"
	RemoveLimitHelp    = "Чтобы удалить лимит, введи:\n" +
		"\"" + removeLimitKeyword + " <категория>\""
)

func (h *RemoveLimit) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.Text, removeLimitKeyword) {
		return utils.HandleSkipped
	}

	category := strings.Trim(strings.TrimPrefix(msg.Text, removeLimitKeyword+" "), " ")
	if category == "" {
		err := h.messageSender.SendText(RemoveLimitHelp, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}
	if err := h.usecase.RemoveLimit(ctx, msg.UserID, category); err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "usecase RemoveLimit failed"))
	}

	response := "Успешно удалили лимит по категории \"" + category + "\""
	if err := h.messageSender.SendText(response, msg.UserID); err != nil {
		return utils.HandleWithErrorOrNil(errors.WithMessage(err, "send success message failed"))
	}

	return utils.HandleSuccess
}

func New(limitUc usecase, messageSender messages.MessageSender) *RemoveLimit {
	return &RemoveLimit{
		usecase:       limitUc,
		messageSender: messageSender,
	}
}
