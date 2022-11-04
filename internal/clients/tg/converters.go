package tg

//goland:noinspection SpellCheckingInspection
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
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
		msg = extractCommonFields(update)
		msg.CallbackData = update.CallbackData()
		ok = true
		return
	} else if update.Message != nil {
		msg = extractCommonFields(update)
		msg.Text = update.Message.Text
		ok = true
		return
	}

	ok = false
	return
}

func extractCommonFields(update tgbotapi.Update) messages.Message {
	return messages.Message{
		UserID:   update.Message.From.ID,
		UserName: update.Message.From.UserName,
		SendTime: update.Message.Time(),
	}
}
