package set_currency

import (
	"context"
	"fmt"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers/callbacks"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers/utils"
)

type userUc interface {
	SetCurrency(ctx context.Context, userID int64, curr entities.Currency) error
}

type SetCurrency struct {
	userUc userUc
	sender messages.MessageSender
}

func New(userUc userUc, sender messages.MessageSender) *SetCurrency {
	return &SetCurrency{
		userUc: userUc,
		sender: sender,
	}
}

func (s *SetCurrency) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.CallbackData, callbacks.SetCurrencyPrefix) {
		return messages.HandleResult{Skipped: true}
	}

	selectedCurr := entities.Currency(strings.TrimPrefix(msg.CallbackData, callbacks.SetCurrencyPrefix))
	if err := s.userUc.SetCurrency(ctx, msg.UserID, selectedCurr); err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	successText := fmt.Sprintf("Валюта %s успешно установлена", selectedCurr)
	err := s.sender.SendText(successText, msg.UserID)

	return utils.HandleWithErrorOrNil(err)
}

func (s *SetCurrency) Name() string {
	return "SetCurrency"
}
