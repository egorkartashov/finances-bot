//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package messages

import (
	"github.com/pkg/errors"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type MessageHandler interface {
	Handle(msg Message) MessageHandleResult
	Name() string
}

type Model struct {
	tgClient MessageSender
	handlers []MessageHandler
}

func New(tgClient MessageSender, handlers []MessageHandler) *Model {
	return &Model{
		tgClient: tgClient,
		handlers: handlers,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type MessageHandleResult struct {
	Skipped bool
	Err     error
}

const UnknownCommandMessage = "не знаю эту команду"

func (m *Model) IncomingMessage(msg Message) error {
	for _, h := range m.handlers {
		res := h.Handle(msg)
		if res.Skipped {
			continue
		}
		if res.Err != nil {
			return errors.WithMessage(res.Err, "IncomingMessage: error in "+h.Name())
		}
		return nil
	}

	return m.tgClient.SendMessage(UnknownCommandMessage, msg.UserID)
}
