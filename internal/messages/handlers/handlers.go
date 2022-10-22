package handlers

import (
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/remove_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/set_limit"
)

type base struct {
	MessageSender messages.MessageSender
}

var NewSetLimit = set_limit.NewSetLimit
var NewRemoveLimit = remove_limit.NewRemoveLimit
