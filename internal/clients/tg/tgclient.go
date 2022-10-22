package tg

//goland:noinspection SpellCheckingInspection
import (
	"context"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.WithMessage(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendText(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.WithMessage(err, "client.Send")
	}

	return nil
}

func (c *Client) SendMessage(userID int64, msg messages.Message) error {
	tgMsg := tgbotapi.NewMessage(userID, msg.Text)
	if msg.InlineKeyboardButtons != nil {
		tgMsg.ReplyMarkup = convertToTgInlineKeyboard(msg.InlineKeyboardButtons)
	}

	_, err := c.client.Send(tgMsg)
	if err != nil {
		return errors.WithMessage(err, "client.Send")
	}

	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("listening for messages")

	for {
		select {
		case update := <-updates:
			msg, ok := convertToMessage(update)
			if ok {
				logger.Infof("[%s] text: %s, callback_data: %s", msg.UserName, msg.Text, msg.CallbackData)
				err := msgModel.IncomingMessage(msg)

				if err != nil {
					logger.Infof("error processing message:", err)
				}
			}
		case <-ctx.Done():
			log.Println("ListenUpdates: exiting due to ctx.Done()")
			return
		}
	}
}
