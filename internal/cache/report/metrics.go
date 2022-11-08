package report

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

var (
	GetReportCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "finances",
			Subsystem: "bot",
			Name:      "get_report_counter",
		},
		[]string{"userId", "period", "result"},
	)
)

type GetReportResult string

const (
	MetricsMiss GetReportResult = "miss"
	MetricsHit  GetReportResult = "hit"
	MetricsErr  GetReportResult = "err"
)

func writeMetrics(userID int64, period entities.ReportPeriod, res GetReportResult) {
	GetReportCounter.WithLabelValues(strconv.FormatInt(userID, 10), string(period), string(res))
}
