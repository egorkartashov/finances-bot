package main

import (
	"context"
	"flag"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	rates_cache "gitlab.ozon.dev/egor.linkinked/finances-bot/internal/cache/rates"
	report_cache "gitlab.ozon.dev/egor.linkinked/finances-bot/internal/cache/report"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/clients/cbrf"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/config"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/limits"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/handlers/aggregate"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/middleware/logging"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/middleware/metrics"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages/middleware/tracing"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/rates"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/grpc/server"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/kafka"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/kafka/producer"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage/tx"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/add_expense"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/get_report"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/register_user"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/remove_limit"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/set_currency"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/usecases/set_limit"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/utils"
	"go.uber.org/zap"
)

var supportedCurrencies = []entities.Currency{currency.EUR, currency.CNY, currency.USD}

var devMode = flag.Bool("devmode", false, "Start bot in development mode")
var httpPort = flag.Int("httpport", 9870, "Port to start HTTP server for metrics")

func main() {
	flag.Parse()
	logger.InitLogger(*devMode)

	_ = godotenv.Load(".env")

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("config init failed", zap.Error(err))
	}

	// External connections
	tgClient := mustConnectToTg(cfg)
	db := mustConnectToDb(cfg.Dsn())
	redisClient := mustConnectToRedisCache(cfg.CacheURL())
	ratesApi := &cbrf.RatesApi{}
	syncProducer := mustCreateKafkaSyncProducer(cfg.KafkaBrokers())

	// Services
	reportCache := report_cache.New(redisClient)

	dbTxStorage := tx.New(db)
	expenseStorage := storage.NewExpenses(dbTxStorage)
	userStorage := storage.NewUsers(dbTxStorage)
	ratesStorage := storage.NewRates(dbTxStorage)
	limitStorage := storage.NewLimits(dbTxStorage, cfg.BaseCurrency())

	ratesProvider := rates.NewProvider(cfg, ratesApi, ratesStorage)
	cachingRatesProvider := rates_cache.NewInMemCacheDecorator(ratesProvider)

	currencyConverter := currency.NewConverter(cfg, cachingRatesProvider, userStorage)
	limitChecker := limits.NewChecker(limitStorage, expenseStorage, userStorage, currencyConverter)

	reportRequestProducer := producer.NewReportRequestProducer(syncProducer, kafka.ReportRequestTopic)

	// Use cases
	registerUserUc := register_user.NewUsecase(cfg, userStorage)
	setCurrencyUc := set_currency.NewUsecase(cfg, userStorage)
	addExpenseUc := add_expense.NewUsecase(
		dbTxStorage, expenseStorage, userStorage, currencyConverter, limitChecker, reportCache,
	)
	getExpensesReportUc := get_report.NewUsecase(userStorage, reportCache, reportRequestProducer)
	setLimitUc := set_limit.NewUsecase(limitStorage, userStorage, currencyConverter)
	removeLimitUc := remove_limit.NewUsecase(limitStorage)

	// Bot messages handlers
	getReportHandler := handlers.NewGetReport(getExpensesReportUc, tgClient)
	var handler messages.MessageHandler = aggregate.NewAggregate(
		handlers.NewStart(registerUserUc, tgClient),
		handlers.NewAddExpense(addExpenseUc, tgClient),
		getReportHandler,
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

	utils.WithGracefulShutdown(
		cancel,
		func() {
			utils.StartMetricsHttpServer(ctx, *httpPort)
		},
		func() {
			clientSecrets := []string{cfg.SendReportClientSecret()}
			startGrpcServer(ctx, cfg.SendReportPort(), clientSecrets, server.NewGrpc(getReportHandler))
		},
		func() {
			rateFetchFreq := time.Duration(cfg.RateFetchFreqMinutes()) * time.Minute
			ratesProvider.UpdateRates(ctx, rateFetchFreq, supportedCurrencies)
		},
		func() {
			tgClient.ListenUpdates(ctx, msgModel)
		},
	)
}
