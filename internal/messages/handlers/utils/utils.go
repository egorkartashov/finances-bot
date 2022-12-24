package utils

import (
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
)

func HandleWithErrorOrNil(err error) messages.HandleResult {
	return messages.HandleResult{Skipped: false, Err: err}
}

var (
	HandleSkipped = messages.HandleResult{Skipped: true, Err: nil}
	HandleSuccess = messages.HandleResult{Skipped: false, Err: nil}
)
