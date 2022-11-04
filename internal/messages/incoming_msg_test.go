package messages_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/mocks"
)

func Test_ShouldSendSomethingHasBrokenMessage_WhenHandledWithErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.TODO()
	msg := messages.Message{
		Text:   "some command",
		UserID: int64(123),
	}

	sender := messages_mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendText(messages.SomethingHasBroken, msg.UserID)

	handler := messages_mocks.NewMockMessageHandler(ctrl)
	handler.EXPECT().Handle(ctx, msg).Return(messages.HandleResult{Skipped: false, Err: errors.New("")})

	model := messages.NewModel(sender, handler)
	model.IncomingMessage(ctx, msg)
}
