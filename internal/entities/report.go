package entities

import "github.com/shopspring/decimal"

type ReportPeriod byte

const (
	ReportFor1Week  ReportPeriod = 1
	ReportFor1Month ReportPeriod = 2
	ReportFor1Year  ReportPeriod = 3
)

type Report struct {
	Cur     Currency
	Entries []ReportEntry
}

type ReportEntry struct {
	Category string
	TotalSum decimal.Decimal
}
