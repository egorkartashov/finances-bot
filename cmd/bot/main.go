package main

import (
	"context"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/cbrf"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/presenters"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/rates"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage/tx"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

const baseCurrency = currency.RUB

var supportedCurrencies = []entities.Currency{currency.EUR, currency.CNY, currency.USD}

func main() {
	_ = godotenv.Load(".env")

	cfg, err := config.New(baseCurrency)
	if err != nil {
		logger.Fatal("config init failed: ", err)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("tg client init failed: ", err)
	}

	db := mustConnectToDb(cfg.Dsn())
	dbTxStorage := tx.New(db)

	expenseStorage := storage.NewExpenses(dbTxStorage)
	userStorage := storage.NewUsers(dbTxStorage)
	ratesStorage := storage.NewRates(dbTxStorage)
	limitStorage := storage.NewLimits(dbTxStorage)

	ratesApi := &cbrf.RatesApi{}
	ratesProvider := rates.NewProvider(cfg, ratesApi, ratesStorage)

	userUc := users.NewUsecase(cfg, userStorage)
	currencyConverter := currency.NewConverter(cfg, ratesProvider, userStorage)
	limitUc := limits.NewUsecase(limitStorage, expenseStorage, currencyConverter)
	expensesUc := expenses.NewUsecase(cfg, dbTxStorage, expenseStorage, userStorage, currencyConverter, limitUc)

	msgHandlers := []messages.MessageHandler{
		handlers.NewStart(userUc, tgClient),
		handlers.NewAddExpense(expensesUc, tgClient),
		handlers.NewReport(expensesUc, presenters.NewReport(), tgClient),
		handlers.NewGetCurrencyOptions(tgClient),
		handlers.NewSetCurrency(userUc, tgClient),
		handlers.NewSetLimit(limitUc, tgClient),
		handlers.NewRemoveLimit(limitUc, tgClient),
	}

	msgModel := messages.New(tgClient, msgHandlers)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		rateFetchFreq := time.Duration(cfg.RateFetchFreqMinutes()) * time.Minute
		ratesProvider.UpdateRates(ctx, rateFetchFreq, supportedCurrencies)
	}()

	go func() {
		tgClient.ListenUpdates(ctx, msgModel)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Ready for stop signal")
	sig := <-sigChan

	logger.Infof("Received %v, gracefully shutting down...\n", sig)
	cancel()

	time.Sleep(10 * time.Second)
}

func mustConnectToDb(dsn string) *sqlx.DB {
	return sqlx.MustConnect("postgres", dsn)
}
