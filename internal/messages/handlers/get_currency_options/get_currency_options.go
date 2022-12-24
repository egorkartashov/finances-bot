package get_currency_options

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers/callbacks"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers/utils"
)

type GetCurrencyOptions struct {
	sender messages.MessageSender
}

const changeCurrencyCommand = "/change_currency"

func (g *GetCurrencyOptions) Handle(_ context.Context, msg messages.Message) messages.HandleResult {
	if msg.Text != changeCurrencyCommand {
		return utils.HandleSkipped
	}

	currencies := []entities.Currency{currency.RUB, currency.EUR, currency.USD, currency.CNY}

	buttons := make([][]messages.InlineKeyboardButton, len(currencies))
	for i, cur := range currencies {
		curStr := string(cur)
		callbackData := callbacks.SetCurrencyPrefix + curStr
		buttons[i] = []messages.InlineKeyboardButton{{Label: curStr, CallbackData: &callbackData}}
	}

	resp := messages.Message{
		Text:                  "Выберите нужную вам валюту",
		InlineKeyboardButtons: buttons,
	}

	err := g.sender.SendMessage(msg.UserID, &resp)
	return utils.HandleWithErrorOrNil(err)
}

func (g *GetCurrencyOptions) Name() string {
	return "GetCurrencyOptions"
}

func New(sender messages.MessageSender) *GetCurrencyOptions {
	return &GetCurrencyOptions{
		sender: sender,
	}
}
