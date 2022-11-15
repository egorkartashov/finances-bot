//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package messages

import (
	"context"
	"time"
)

type MessageSender interface {
	SendText(text string, userID int64) error
	SendMessage(userID int64, msg *Message) error
}

type MessageHandler interface {
	Handle(ctx context.Context, msg Message) HandleResult
}

type NamedHandler interface {
	MessageHandler
	Name() string
}

type Model struct {
	messageSender MessageSender
	handler       MessageHandler
}

func NewModel(tgClient MessageSender, handler MessageHandler) *Model {
	return &Model{
		messageSender: tgClient,
		handler:       handler,
	}
}

type Message struct {
	UserID                int64
	UserName              string
	SendTime              time.Time
	Text                  string
	CallbackData          string
	InlineKeyboardButtons [][]InlineKeyboardButton
}

type InlineKeyboardButton struct {
	Label        string
	CallbackData *string
}

type HandleResult struct {
	Skipped     bool
	Err         error
	HandlerName string
}

const SomethingHasBroken = "Что-то сломалось :o"

func (m *Model) IncomingMessage(ctx context.Context, msg Message) {
	res := m.handler.Handle(ctx, msg)
	if res.Err != nil {
		_ = m.messageSender.SendText(SomethingHasBroken, msg.UserID)
		return
	}
}
