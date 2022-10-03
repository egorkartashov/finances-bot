package handlers

import "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"

type baseHandler struct {
	messageSender messages.MessageSender
}
