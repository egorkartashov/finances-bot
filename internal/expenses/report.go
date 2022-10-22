package expenses

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type ReportPeriod byte

const (
	ReportFor1Week  ReportPeriod = 1
	ReportFor1Month ReportPeriod = 2
	ReportFor1Year  ReportPeriod = 3
)

type Report struct {
	Cur     entities.Currency
	Entries []ReportEntry
}

type ReportEntry struct {
	Category string
	TotalSum decimal.Decimal
}

func (u *Usecase) GenerateReport(ctx context.Context, userID int64, reportPeriod ReportPeriod) (
	report *Report, err error,
) {
	minTime, err := getReportMinTime(reportPeriod)
	if err != nil {
		return
	}

	expenses, err := u.expenseStorage.GetExpenses(ctx, userID, minTime)
	if err != nil {
		return
	}

	user, ok, err := u.userUc.Get(ctx, userID)
	if err != nil {
		return
	}
	if !ok {
		err = users.NewUserNotFoundErr(userID)
		return
	}

	sumByCategory, err := u.getSumByCategory(ctx, expenses, user)
	if err != nil {
		return
	}

	entries := make([]ReportEntry, len(sumByCategory))
	var i = 0
	for cat, sum := range sumByCategory {
		entries[i] = ReportEntry{Category: cat, TotalSum: sum}
		i += 1
	}

	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Category < entries[j].Category })
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].TotalSum.GreaterThan(entries[j].TotalSum) })

	return &Report{Cur: user.Currency, Entries: entries}, nil
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

func (u *Usecase) getSumByCategory(
	ctx context.Context, expenses []entities.Expense, user entities.User,
) (sumByCategory map[string]decimal.Decimal, err error) {
	baseCurrency := u.cfg.BaseCurrency()
	sumByCategory = make(map[string]decimal.Decimal)
	for _, e := range expenses {
		var convertedSum decimal.Decimal
		convertedSum, _, err = u.currencyConverter.Convert(ctx, e.Sum, baseCurrency, user.Currency, e.Date)
		if err != nil {
			return
		}
		if curSum, ok := sumByCategory[e.Category]; !ok {
			sumByCategory[e.Category] = convertedSum
		} else {
			sumByCategory[e.Category] = curSum.Add(convertedSum)
		}
	}
	return
}
