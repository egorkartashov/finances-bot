package handlers

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/callbacks"
)

type userUc interface {
	SetCurrency(userID int64, curr currency.Currency) error
}

type SetCurrency struct {
	userUc userUc
	sender messages.MessageSender
}

func (s *SetCurrency) Handle(msg messages.Message) messages.HandleResult {
	if !strings.HasPrefix(msg.CallbackData, callbacks.SetCurrencyPrefix) {
		return messages.HandleResult{Skipped: true}
	}

	selectedCurr := currency.Currency(strings.TrimPrefix(msg.CallbackData, callbacks.SetCurrencyPrefix))
	if err := s.userUc.SetCurrency(msg.UserID, selectedCurr); err != nil {
		return handleWithErrorOrNil(err)
	}

	successText := fmt.Sprintf("Валюта %s успешно установлена", selectedCurr)
	err := s.sender.SendText(successText, msg.UserID)

	return handleWithErrorOrNil(err)
}

func (s *SetCurrency) Name() string {
	return "SetCurrency"
}

func NewSetCurrency(userUc userUc, sender messages.MessageSender) *SetCurrency {
	return &SetCurrency{
		userUc: userUc,
		sender: sender,
	}
}
