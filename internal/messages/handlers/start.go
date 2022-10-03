package handlers

import (
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
)

type StartHandler struct {
	baseHandler
}

func NewStartHandler(sender messages.MessageSender) *StartHandler {
	return &StartHandler{
		baseHandler: baseHandler{sender},
	}
}

func (h *StartHandler) Handle(msg messages.Message) messages.MessageHandleResult {
	if msg.Text != "/start" {
		return messages.MessageHandleResult{Skipped: true, Err: nil}
	}

	var welcomeMessage = "Привет! Я - телеграм бот для учета финансов. \n" +
		"Пока я могу только сохранять твои траты и формировать отчет по сумма трат по категориям. \n\n" +
		"Чтобы добавить трату, введи ее в следующем формате: " + ExpenseFormat + "\n\n" +
		"Чтобы получить отчет, введи команду: " + ReportFormatMessage

	err := h.messageSender.SendMessage(welcomeMessage, msg.UserID)
	return messages.MessageHandleResult{Skipped: false, Err: err}
}

func (h *StartHandler) Name() string {
	return "StartHandler"
}
