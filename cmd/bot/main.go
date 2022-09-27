package main

import (
	"log"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/model/messages"
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

	msgModel := messages.New(tgClient)

	tgClient.ListenUpdates(msgModel)
}
