package unknown_command

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
)

type UnknownCommand struct {
	sender messages.MessageSender
}

func New(sender messages.MessageSender) *UnknownCommand {
	return &UnknownCommand{sender}
}

const UnknownCommandMessage = "не знаю эту команду"

func (u *UnknownCommand) Handle(_ context.Context, msg messages.Message) messages.HandleResult {
	err := u.sender.SendText(UnknownCommandMessage, msg.UserID)

	return messages.HandleResult{
		Skipped: false,
		Err:     err,
	}
}

func (u *UnknownCommand) Name() string {
	return "UnknownCommand"
}
