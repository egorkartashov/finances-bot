//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package messages

import (
	"context"

	"github.com/pkg/errors"
)

type MessageSender interface {
	SendText(text string, userID int64) error
	SendMessage(userID int64, msg Message) error
}

type MessageHandler interface {
	Handle(ctx context.Context, msg Message) HandleResult
	Name() string
}

type Model struct {
	messageSender MessageSender
	handlers      []MessageHandler
}

func New(tgClient MessageSender, handlers []MessageHandler) *Model {
	return &Model{
		messageSender: tgClient,
		handlers:      handlers,
	}
}

type Message struct {
	Text                  string
	UserID                int64
	UserName              string
	CallbackData          string
	InlineKeyboardButtons [][]InlineKeyboardButton
}

type InlineKeyboardButton struct {
	Label        string
	CallbackData *string
}

type HandleResult struct {
	Skipped bool
	Err     error
}

const UnknownCommandMessage = "не знаю эту команду"

func (m *Model) IncomingMessage(msg Message) error {
	for _, h := range m.handlers {
		res := h.Handle(context.Background(), msg)
		if res.Skipped {
			continue
		}
		if res.Err != nil {
			_ = m.messageSender.SendText("Что-то сломалось :o", msg.UserID)
			return errors.WithMessage(res.Err, "IncomingMessage: error in "+h.Name()+" handler")
		}
		return nil
	}

	return m.messageSender.SendText(UnknownCommandMessage, msg.UserID)
}
