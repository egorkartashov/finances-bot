package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/messages"
)

var (
	SummaryResponseTime = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "finances",
			Subsystem: "bot",
			Name:      "summary_response_time_seconds",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.01,
				0.99: 0.001,
			},
			MaxAge: 1 * time.Minute,
		},
		[]string{"userId", "command"},
	)
	SummaryProcessingTime = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "finances",
			Subsystem: "bot",
			Name:      "summary_processing_time_seconds",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.01,
				0.99: 0.001,
			},
			MaxAge: 1 * time.Minute,
		},
		[]string{"userId", "command"},
	)
)

func Middleware(handler messages.MessageHandler) messages.MessageHandler {
	return messages.NewHandleFunc(
		func(ctx context.Context, msg messages.Message) messages.HandleResult {
			sendTime := msg.SendTime
			start := time.Now()
			res := handler.Handle(ctx, msg)

			responseTime := time.Since(sendTime)
			processingTime := time.Since(start)

			userID := strconv.FormatInt(msg.UserID, 10)
			SummaryResponseTime.
				WithLabelValues(userID, res.HandlerName).
				Observe(responseTime.Seconds())

			SummaryProcessingTime.
				WithLabelValues(userID, res.HandlerName).
				Observe(processingTime.Seconds())

			return res
		},
	)
}
