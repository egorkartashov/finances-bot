package handlers

import (
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
)

type Start struct {
	base
}

func NewStart(sender messages.MessageSender) *Start {
	return &Start{
		base: base{sender},
	}
}

const welcomeMessage = "Привет! Я - телеграм бот для учета финансов. \n" +
	"Пока я могу только сохранять твои траты и формировать отчет по сумма трат по категориям. \n\n" +
	"Чтобы добавить трату, введи ее в следующем формате: " + ExpenseFormat + "\n\n" +
	"Чтобы получить отчет, введи команду: " + ReportFormatMessage

func (h *Start) Handle(msg messages.Message) messages.HandleResult {
	if msg.Text != "/start" {
		return handleSkipped
	}

	err := h.messageSender.SendText(welcomeMessage, msg.UserID)
	return handleWithErrorOrNil(err)
}

func (h *Start) Name() string {
	return "Start"
}
