//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package expenses

import (
	"fmt"
	"sort"
	"time"
)

type Storage interface {
	AddExpense(userID int64, exp Expense)
	GetCategoriesTotals(userID int64, minTime time.Time) []CategoryTotal
}

type Model struct {
	storage Storage
}

func New(storage Storage) *Model {
	return &Model{
		storage: storage,
	}
}

type Expense struct {
	Category string
	SumRub   int32
	Date     time.Time
}

type CategoryTotal struct {
	Category    string
	TotalSumRub int32
}

func (m *Model) AddExpense(userID int64, exp Expense) {
	m.storage.AddExpense(userID, exp)
}

type ReportPeriod byte

const (
	ReportFor1Week  ReportPeriod = 1
	ReportFor1Month ReportPeriod = 2
	ReportFor1Year  ReportPeriod = 3
)

type Report struct {
	Entries []ReportEntry
}

type ReportEntry struct {
	Category    string
	TotalSumKop int32
}

func (m *Model) GenerateReport(userID int64, reportPeriod ReportPeriod) (*Report, error) {
	minTime, err := getReportMinTime(reportPeriod)
	if err != nil {
		return nil, err
	}

	categoriesTotals := m.storage.GetCategoriesTotals(userID, minTime)
	entries := make([]ReportEntry, len(categoriesTotals))
	for i, cat := range categoriesTotals {
		entries[i] = ReportEntry{
			Category:    cat.Category,
			TotalSumKop: cat.TotalSumRub,
		}
	}

	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Category < entries[j].Category })
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].TotalSumKop > entries[j].TotalSumKop })

	return &Report{Entries: entries}, nil
}

func getReportMinTime(period ReportPeriod) (time.Time, error) {
	now := time.Now().UTC()
	switch period {
	case ReportFor1Week:
		return now.AddDate(0, 0, -7), nil
	case ReportFor1Month:
		return now.AddDate(0, -1, 0), nil
	case ReportFor1Year:
		return now.AddDate(-1, 0, 0), nil
	default:
		return time.Now(), fmt.Errorf("getReportMinTime for period %v is not implemented", period)
	}
}
