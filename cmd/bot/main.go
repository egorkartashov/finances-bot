package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/cbrf"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/presenters"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/rates"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

const baseCurrency = currency.RUB

var supportedCurrencies = []currency.Currency{currency.EUR, currency.CNY, currency.USD}

func main() {
	l := log.New(os.Stdout, "finances-bot: ", log.LstdFlags)

	cfg, err := config.New()
	if err != nil {
		l.Fatal("config init failed: ", err)
	}

	tgClient, err := tg.New(cfg, l)
	if err != nil {
		l.Fatal("tg client init failed: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	expenseStorage := storage.NewExpenses()
	userStorage := storage.NewUsers()
	ratesStorage, err := storage.NewRates()
	if err != nil {
		l.Fatalln(errors.WithMessage(err, "failed to init rates storage"))
	}

	ratesApi := &cbrf.RatesApi{}
	ratesProvider := rates.NewProvider(l, baseCurrency, ratesApi, ratesStorage)
	go func() {
		rateFetchFreq := time.Duration(cfg.RateFetchFreqMinutes()) * time.Minute
		ratesProvider.UpdateRates(ctx, rateFetchFreq, supportedCurrencies)
	}()

	currencyConverter := currency.NewConverter(ratesProvider)
	userUc := users.NewUsecase(userStorage)
	expensesUc := expenses.NewUsecase(expenseStorage, baseCurrency, userUc, currencyConverter, l)

	msgHandlers := []messages.MessageHandler{
		handlers.NewStart(tgClient),
		handlers.NewAddExpense(expensesUc, tgClient),
		handlers.NewReport(expensesUc, presenters.NewReport(), tgClient),
		handlers.NewGetCurrencyOptions(tgClient),
		handlers.NewSetCurrency(userUc, tgClient),
	}

	msgModel := messages.New(tgClient, msgHandlers)

	go func() {
		tgClient.ListenUpdates(ctx, msgModel)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	l.Println("Waiting for stop signal")
	sig := <-sigChan

	l.Printf("Received %v, gracefully shutting down...\n", sig)
	cancel()

	time.Sleep(10 * time.Second)
}
