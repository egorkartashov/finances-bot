package start

import (
	"context"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/add_expense"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/get_report"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/remove_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/set_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type Start struct {
	userUc userUc
	sender messages.MessageSender
}

func New(uu userUc, s messages.MessageSender) *Start {
	return &Start{
		userUc: uu,
		sender: s,
	}
}

var (
	helpMessages = []string{
		add_expense.AddExpenseHelp, get_report.ReportHelp, set_limit.AddLimitHelp, remove_limit.RemoveLimitHelp,
	}
	welcomeMessage1 = "Привет! Я - телеграм бот для учета финансов. Вот, что я могу:" + helpSeparator +
		strings.Join(helpMessages, helpSeparator)
)

const (
	helpSeparator = "\n\n"
)

func (h *Start) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	if msg.Text != "/start" {
		return utils.HandleSkipped
	}

	err := h.userUc.Register(ctx, msg.UserID)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	err = h.sender.SendText(welcomeMessage1, msg.UserID)
	return utils.HandleWithErrorOrNil(err)
}

func (h *Start) Name() string {
	return "Start"
}
