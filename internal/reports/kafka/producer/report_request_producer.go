package producer

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
)

type ReportRequestProducer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewReportRequestProducer(sp sarama.SyncProducer, topic string) *ReportRequestProducer {
	return &ReportRequestProducer{
		syncProducer: sp,
		topic:        topic,
	}
}

func (p *ReportRequestProducer) Send(ctx context.Context, req *reports.NewReportRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "report-request-kafka")
	defer span.Finish()

	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(bytes),
	}
	if _, _, err = p.syncProducer.SendMessage(&msg); err != nil {
		return errors.WithMessage(err, "failed to produce report request")
	}
	return nil
}
