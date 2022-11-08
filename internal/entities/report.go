package entities

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ReportPeriod byte

const (
	ReportFor1Week  ReportPeriod = 1
	ReportFor1Month ReportPeriod = 2
	ReportFor1Year  ReportPeriod = 3
)

type Report struct {
	Cur     Currency      `json:"curr"`
	Entries []ReportEntry `json:"entries"`
}

type ReportEntry struct {
	Category string          `json:"category"`
	TotalSum decimal.Decimal `json:"totalSum"`
}

func GetAllReportPeriods() []ReportPeriod {
	return []ReportPeriod{
		ReportFor1Week,
		ReportFor1Month,
		ReportFor1Year,
	}
}

func (p ReportPeriod) GetStartDate(endDate time.Time) (time.Time, error) {
	switch p {
	case ReportFor1Week:
		return endDate.AddDate(0, 0, -7), nil
	case ReportFor1Month:
		return endDate.AddDate(0, -1, 0), nil
	case ReportFor1Year:
		return endDate.AddDate(-1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("getReportStartDate for period %v is not implemented", p)
	}
}
