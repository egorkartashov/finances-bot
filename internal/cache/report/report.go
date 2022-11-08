package report

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/utils"
)

type Report struct {
	client *redis.Client
}

func New(client *redis.Client) *Report {
	return &Report{client}
}

func (r *Report) Get(ctx context.Context, userID int64, period entities.ReportPeriod) (*entities.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache-get-report")
	defer span.Finish()

	span.SetTag("userID", userID)
	span.SetTag("period", period)

	key := getKey(userID, period)
	jsonVal, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		writeMetrics(userID, period, MetricsMiss)
		return nil, nil
	}
	if err != nil {
		writeMetrics(userID, period, MetricsErr)
		return nil, err
	}

	var report entities.Report
	if jsonErr := json.Unmarshal([]byte(jsonVal), &report); jsonErr != nil {
		writeMetrics(userID, period, MetricsErr)
		return nil, jsonErr
	}

	writeMetrics(userID, period, MetricsHit)
	return &report, nil
}

func (r *Report) Save(
	ctx context.Context, userID int64, period entities.ReportPeriod, report *entities.Report,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache-save-report")
	defer span.Finish()

	span.SetTag("userID", userID)
	span.SetTag("period", period)

	key := getKey(userID, period)
	jsonBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	expiration := time.Until(utils.GetTomorrow(time.Now()))
	return r.client.Set(ctx, key, string(jsonBytes), expiration).Err()
}

func (r *Report) DeleteAffected(ctx context.Context, userID int64, newExpenseDate time.Time) error {
	periods := entities.GetAllReportPeriods()
	now := time.Now().UTC()
	for _, period := range periods {
		startDate, err := period.GetStartDate(now)
		if err != nil {
			continue
		}
		if startDate.Before(newExpenseDate) || startDate.Equal(newExpenseDate) {
			if err := r.Delete(ctx, userID, period); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Report) Delete(ctx context.Context, userID int64, period entities.ReportPeriod) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache-del-report")
	defer span.Finish()

	span.SetTag("userID", userID)
	span.SetTag("period", period)

	key := getKey(userID, period)
	return r.client.Del(ctx, key).Err()
}

func getKey(userID int64, period entities.ReportPeriod) string {
	return fmt.Sprintf("%s_%v", getPrefix(userID), period)
}

func getPrefix(userID int64) string {
	return fmt.Sprintf("report_%v", userID)
}
