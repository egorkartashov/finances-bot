package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
	"go.uber.org/zap"
)

type reportGenerator interface {
	Generate(ctx context.Context, req *reports.NewReportRequest) (*reports.GeneratedReportResponse, error)
}

type reportSender interface {
	reports.FinishedReportSender
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	reportGenerator reportGenerator
	client          reportSender
	cleanUpFunc     func()
}

func MustNew(g reportGenerator, c reportSender, cleanUpFunc func()) *Consumer {
	return &Consumer{
		reportGenerator: g,
		client:          c,
		cleanUpFunc:     cleanUpFunc,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	logger.Info("consumer - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.Info("consumer - cleanup start")
	consumer.cleanUpFunc()
	logger.Info("consumer - cleanup end")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		consumer.processMessage(session, message)
	}

	return nil
}

func (consumer *Consumer) processMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	span, ctx := opentracing.StartSpanFromContext(session.Context(), "consume-report-request")
	defer span.Finish()

	failed := false
	defer func() {
		if !failed {
			session.MarkMessage(message, "")
		}
	}()

	req, err := deserializeRequest(message.Value)
	if err != nil {
		logDeserializeFail(err, message)
		return
	}

	generatedReportResponse, err := consumer.reportGenerator.Generate(ctx, req)
	if err != nil {
		ackMessage := isExpectedError(err)
		logger.Error(
			fmt.Sprintf("failed to generate report, ack this message: %v", ackMessage),
			zap.Error(err),
		)
		failed = !ackMessage
		return
	}

	finishedReportReq := reports.FormattedReport{
		NewReportRequest: reports.NewReportRequest{
			UserID:   req.UserID,
			Format:   req.Format,
			Period:   req.Period,
			Date:     req.Date,
			Currency: req.Currency,
		},
		Payload: generatedReportResponse.Payload,
	}
	if err = consumer.client.Send(ctx, &finishedReportReq); err != nil {
		logger.Error("failed to send report", zap.Error(err))
		failed = true
		return
	}
}

func deserializeRequest(bytes []byte) (*reports.NewReportRequest, error) {
	req := reports.NewReportRequest{}
	err := json.Unmarshal(bytes, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func logDeserializeFail(err error, message *sarama.ConsumerMessage) {
	logger.Error(
		"failed to deserialize message",
		zap.Error(err),
		zap.ByteString("key", message.Key),
		zap.Int64("offset", message.Offset),
	)
}

func isExpectedError(err error) bool {
	if _, ok := err.(reports.ErrUnsupportedFormat); ok {
		return true
	}
	return false
}
