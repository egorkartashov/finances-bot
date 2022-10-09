package handlers

import "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"

type base struct {
	messageSender messages.MessageSender
}
