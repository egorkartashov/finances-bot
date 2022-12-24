package report

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/reports"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/utils"
)

type Report struct {
	client *redis.Client
}

func New(client *redis.Client) *Report {
	return &Report{client}
}

func (r *Report) Get(ctx context.Context, req *reports.NewReportRequest) (*reports.FormattedReport, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache-get-report")
	defer span.Finish()

	span.SetTag("userID", req.UserID)
	span.SetTag("period", req.Period)
	span.SetTag("type", req.Format)

	key := getKey(req)
	reportPayload, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		writeMetrics(req, MetricsMiss)
		return nil, nil
	}
	if err != nil {
		writeMetrics(req, MetricsErr)
		return nil, err
	}

	writeMetrics(req, MetricsHit)
	report := &reports.FormattedReport{
		NewReportRequest: *req,
		Payload:          reportPayload,
	}
	return report, nil
}

func (r *Report) Save(ctx context.Context, report *reports.FormattedReport) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache-save-report")
	defer span.Finish()

	span.SetTag("userID", report.UserID)
	span.SetTag("period", report.Period)
	span.SetTag("type", report.Format)

	key := getKey(&report.NewReportRequest)
	expiration := time.Until(utils.GetTomorrow(time.Now()))

	return r.client.Set(ctx, key, report.Payload, expiration).Err()
}

func (r *Report) DeleteAffected(ctx context.Context, userID int64, newExpenseDate time.Time) error {
	allPeriods := entities.GetAllReportPeriods()
	now := time.Now().UTC()
	for _, period := range allPeriods {
		periodStart, err := period.GetStartDate(now)
		if err != nil || periodStart.After(newExpenseDate) {
			continue
		}

		prefix := getUserWithPeriodPrefix(userID, period)
		if err := r.deleteByPrefix(ctx, prefix); err != nil {
			return err
		}
	}
	return nil
}

func getKey(req *reports.NewReportRequest) string {
	return fmt.Sprintf("report_%v_%v_%s_%v", req.UserID, req.Period, req.Currency, req.Format)
}

func getUserWithPeriodPrefix(userID int64, period entities.ReportPeriod) string {
	return fmt.Sprintf("report_%v_%v", userID, period)
}

func (r *Report) deleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	match := prefix + "*"
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			return err
		}
		r.client.Del(ctx, keys...)

		if cursor == 0 {
			break
		}
	}
	return nil
}
