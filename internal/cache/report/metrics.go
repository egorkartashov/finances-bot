package report

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
)

var (
	GetReportCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "finances",
			Subsystem: "bot",
			Name:      "get_report_counter",
		},
		[]string{"userId", "period", "format", "result"},
	)
)

type GetReportResult string

const (
	MetricsMiss GetReportResult = "miss"
	MetricsHit  GetReportResult = "hit"
	MetricsErr  GetReportResult = "err"
)

func writeMetrics(req *reports.NewReportRequest, res GetReportResult) {
	GetReportCounter.WithLabelValues(
		strconv.FormatInt(req.UserID, 10),
		string(req.Period),
		string(req.Format),
		string(res),
	)
}
