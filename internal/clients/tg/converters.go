package tg

//goland:noinspection SpellCheckingInspection
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
)

func convertToTgInlineKeyboard(buttons [][]messages.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	tgRows := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for rIdx, row := range buttons {
		tgRows[rIdx] = make([]tgbotapi.InlineKeyboardButton, len(row))
		for i, btn := range row {
			tgRows[rIdx][i] = tgbotapi.InlineKeyboardButton{
				Text:         btn.Label,
				CallbackData: btn.CallbackData,
			}
		}
	}

	tgInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
	return tgInlineKeyboard
}

func convertToMessage(update tgbotapi.Update) (msg messages.Message, ok bool) {
	if update.CallbackQuery != nil {
		data := update.CallbackQuery
		msg = messages.Message{
			UserID:       data.From.ID,
			UserName:     data.From.UserName,
			SendTime:     data.Message.Time(),
			CallbackData: update.CallbackData(),
		}
		ok = true
		return
	} else if update.Message != nil {
		data := update.Message
		msg = messages.Message{
			UserID:   data.From.ID,
			UserName: data.From.UserName,
			Text:     data.Text,
			SendTime: data.Time(),
		}
		ok = true
		return
	}

	ok = false
	return
}
