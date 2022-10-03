package main

import (
	"log"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	handlers "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/presenters"
	expensesstorage "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage/expenses"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed: ", err)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed: ", err)
	}

	expensesModel := expenses.New(expensesstorage.NewInMem())
	msgHandlers := []messages.MessageHandler{
		handlers.NewStartHandler(tgClient),
		handlers.NewExpenseHandler(expensesModel, tgClient),
		handlers.NewReportHandler(expensesModel, presenters.NewReportPresenter(), tgClient),
	}

	msgModel := messages.New(tgClient, msgHandlers)

	tgClient.ListenUpdates(msgModel)
}
