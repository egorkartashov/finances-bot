package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	rates_cache "gitlab.ozon.dev/egor.linkinked/finances-bot/internal/cache/rates"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/clients/cbrf"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/config"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/send_report"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/rates"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports"
	report_generator "gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/generator"
	report_presenters "gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/generator/presenters"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/grpc/client"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/kafka"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports/kafka/consumer"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/storage/tx"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	KafkaConsumerGroup = "report-consumer-group"
)

const (
	ServiceName = "reportconsumer"
)

var devMode = flag.Bool("devmode", false, "Start bot in development mode")
var port = flag.Int("httpport", 9871, "Port to start HTTP server for metrics")

func main() {
	flag.Parse()
	logger.InitLogger(*devMode)
	initTracing(ServiceName)

	_ = godotenv.Load(".env")

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("consumerCfg init failed", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	utils.WithGracefulShutdown(
		cancel,
		func() {
			mustStartConsumerGroup(ctx, cfg)
		},
		func() {
			utils.StartMetricsHttpServer(ctx, *port)
		},
	)
}

func mustStartConsumerGroup(ctx context.Context, cfg *config.Service) {
	reportRequestConsumer := constructConsumer(cfg)
	consumerGroup, err := constructConsumerGroup(cfg)
	if err != nil {
		logger.Fatal("failed to start consumer group", zap.Error(err))
	}

	err = consumerGroup.Consume(ctx, []string{kafka.ReportRequestTopic}, reportRequestConsumer)
	if err != nil {
		logger.Fatal("failed to start consuming messages", zap.Error(err))
	}
}

func constructConsumer(cfg *config.Service) *consumer.Consumer {
	db := mustConnectToDb(cfg)
	ratesApi := &cbrf.RatesApi{}
	grpcConn, err := connectToGrpcServer(cfg)
	if err != nil {
		logger.Fatal("failed to connect to gRPC for sending reports", zap.Error(err))
	}

	dbTxStorage := tx.New(db)
	expenseStorage := storage.NewExpenses(dbTxStorage)
	userStorage := storage.NewUsers(dbTxStorage)
	ratesStorage := storage.NewRates(dbTxStorage)

	ratesProvider := rates.NewProvider(cfg, ratesApi, ratesStorage)
	cachingRatesProvider := rates_cache.NewInMemCacheDecorator(ratesProvider)
	currencyConverter := currency.NewConverter(cfg, cachingRatesProvider, userStorage)

	reportGenerator := report_generator.New(
		expenseStorage, currencyConverter,
		[]report_generator.ReportPresenter{
			report_presenters.NewFormatMessage(),
		},
	)
	finishedReportSender := constructFinishedReportSender(grpcConn)
	cleanUpFunc := func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close db during cleanup", zap.Error(err))
		}
	}
	reportRequestConsumer := consumer.MustNew(reportGenerator, finishedReportSender, cleanUpFunc)
	return reportRequestConsumer
}

func constructFinishedReportSender(grpcConn *grpc.ClientConn) reports.FinishedReportSender {
	sendReportClient := send_report.NewReportSenderClient(grpcConn)
	finishedReportSender := client.NewGrpcAdapter(sendReportClient)
	return finishedReportSender
}

func constructConsumerGroup(cfg *config.Service) (sarama.ConsumerGroup, error) {
	consumerCfg := sarama.NewConfig()
	consumerCfg.Version = sarama.V2_5_0_0
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}

	log.Printf("Kafka brokers: %s", strings.Join(cfg.KafkaBrokers(), ", "))
	return sarama.NewConsumerGroup(cfg.KafkaBrokers(), KafkaConsumerGroup, consumerCfg)
}
