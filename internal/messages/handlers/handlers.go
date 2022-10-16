package handlers

import "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"

type base struct {
	messageSender messages.MessageSender
}

func handleWithErrorOrNil(err error) messages.HandleResult {
	return messages.HandleResult{Skipped: false, Err: err}
}

var (
	handleSkipped = messages.HandleResult{Skipped: true, Err: nil}
)
