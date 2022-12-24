package main

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/clients/tg"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/config"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"go.uber.org/zap"
)

func mustConnectToTg(cfg *config.Service) *tg.Client {
	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("tg client init failed", zap.Error(err))
	}
	logger.Info("connected to Telegram API")
	return tgClient
}

func mustConnectToDb(dsn string) *sqlx.DB {
	db := sqlx.MustConnect("postgres", dsn)
	logger.Info("connected to Postgres")
	return db
}

func mustConnectToRedisCache(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Fatal("failed to parse Redis URL", zap.Error(err))
	}
	redisClient := redis.NewClient(opt)
	logger.Info("connected to Redis")
	return redisClient
}

func mustCreateKafkaSyncProducer(brokersUrls []string) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_5_0_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Retry.Backoff = time.Millisecond * 250
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokersUrls, cfg)
	if err != nil {
		logger.Fatal("failed to create kafka sync producer", zap.Error(err))
	}
	logger.Info("Connected to Kafka")
	return producer
}
