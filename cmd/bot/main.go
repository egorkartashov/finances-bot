package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/cbrf"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/aggregate"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/middleware/logging"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/middleware/metrics"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/middleware/tracing"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/presenters"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/rates"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/storage/tx"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/add_expense"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/get_expenses_report"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/register_user"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/remove_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_limit"
	"go.uber.org/zap"
)

const baseCurrency = currency.RUB

var supportedCurrencies = []entities.Currency{currency.EUR, currency.CNY, currency.USD}

//goland:noinspection GoUnusedGlobalVariable
var devMode = flag.Bool("devmode", false, "Start bot in development mode")

func main() {
	flag.Parse()
	logger.InitLogger(*devMode)

	_ = godotenv.Load(".env")

	cfg, err := config.New(baseCurrency)
	if err != nil {
		logger.Fatal("config init failed", zap.Error(err))
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("tg client init failed", zap.Error(err))
	}

	db := mustConnectToDb(cfg.Dsn())
	dbTxStorage := tx.New(db)

	expenseStorage := storage.NewExpenses(dbTxStorage)
	userStorage := storage.NewUsers(dbTxStorage)
	ratesStorage := storage.NewRates(dbTxStorage)
	limitStorage := storage.NewLimits(dbTxStorage, cfg.BaseCurrency())

	ratesApi := &cbrf.RatesApi{}
	ratesProvider := rates.NewProvider(cfg, ratesApi, ratesStorage)
	currencyConverter := currency.NewConverter(cfg, ratesProvider, userStorage)
	limitChecker := limits.NewChecker(limitStorage, expenseStorage, userStorage, currencyConverter)

	registerUserUc := register_user.NewUsecase(cfg, userStorage)
	setCurrencyUc := set_currency.NewUsecase(cfg, userStorage)
	addExpenseUc := add_expense.NewUsecase(dbTxStorage, expenseStorage, userStorage, currencyConverter, limitChecker)
	getExpensesReportUc := get_expenses_report.NewUsecase(expenseStorage, userStorage, currencyConverter)
	setLimitUc := set_limit.NewUsecase(limitStorage, userStorage, currencyConverter)
	removeLimitUc := remove_limit.NewUsecase(limitStorage)

	var handler messages.MessageHandler = aggregate.NewAggregate(
		handlers.NewStart(registerUserUc, tgClient),
		handlers.NewAddExpense(addExpenseUc, tgClient),
		handlers.NewGetReport(getExpensesReportUc, presenters.NewReport(), tgClient),
		handlers.NewGetCurrencyOptions(tgClient),
		handlers.NewSetCurrency(setCurrencyUc, tgClient),
		handlers.NewSetLimit(setLimitUc, tgClient),
		handlers.NewRemoveLimit(removeLimitUc, tgClient),
		handlers.NewUnknownCommand(tgClient),
	)
	handler = logging.Middleware(handler)
	handler = tracing.Middleware(handler, cfg)
	handler = metrics.Middleware(handler)

	msgModel := messages.NewModel(tgClient, handler)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		startHttpServer(ctx)
	}()
	go func() {
		defer wg.Done()
		rateFetchFreq := time.Duration(cfg.RateFetchFreqMinutes()) * time.Minute
		ratesProvider.UpdateRates(ctx, rateFetchFreq, supportedCurrencies)
	}()
	go func() {
		defer wg.Done()
		tgClient.ListenUpdates(ctx, msgModel)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Ready for stop signal")
	sig := <-sigChan

	logger.Info("Received %v, gracefully shutting down...\n", zap.String("signal", sig.String()))
	cancel()

	waitWithTimeout(&wg, 10*time.Second)
}

func mustConnectToDb(dsn string) *sqlx.DB {
	return sqlx.MustConnect("postgres", dsn)
}

func startHttpServer(ctx context.Context) {
	srv := http.Server{
		Addr: ":9876",
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("failed to shutdown server", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) (timedOut bool) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
