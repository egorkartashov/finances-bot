package messages_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/mocks"
)

func Test_ShouldGoThroughAllHandlers_IfAllHandlersSkip(t *testing.T) {
	ctrl := gomock.NewController(t)

	msg := messages.Message{
		Text:   "some command",
		UserID: int64(123),
	}

	sender := messages_mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage(messages.UnknownCommandMessage, msg.UserID)

	handlers := make([]messages.MessageHandler, 3)
	for i := range handlers {
		h := messages_mocks.NewMockMessageHandler(ctrl)
		h.EXPECT().Handle(msg).Return(messages.MessageHandleResult{Skipped: true, Err: nil})
		handlers[i] = h
	}

	model := messages.New(sender, handlers)
	err := model.IncomingMessage(msg)

	assert.NoError(t, err)
}

func Test_ShouldStopAfterFirstNotSkippedHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	msg := messages.Message{
		Text:   "some command",
		UserID: int64(123),
	}

	sender := messages_mocks.NewMockMessageSender(ctrl)

	handlers := make([]messages.MessageHandler, 3)
	for i := range handlers {
		h := messages_mocks.NewMockMessageHandler(ctrl)
		handlers[i] = h
		if i == 0 {
			h.EXPECT().Handle(msg).Return(messages.MessageHandleResult{Skipped: true, Err: nil})
		} else if i == 1 {
			h.EXPECT().Handle(msg).Return(messages.MessageHandleResult{Skipped: false, Err: nil})
		} else {
			break
		}
	}

	model := messages.New(sender, handlers)
	err := model.IncomingMessage(msg)

	assert.NoError(t, err)
}
