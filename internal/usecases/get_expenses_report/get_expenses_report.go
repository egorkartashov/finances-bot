package get_expenses_report

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/set_currency"
)

type Usecase struct {
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userStorage       userStorage
}

func (u *Usecase) GenerateReport(ctx context.Context, userID int64, reportPeriod entities.ReportPeriod) (
	report *entities.Report, err error,
) {
	startDate, err := getReportStartDate(reportPeriod)
	if err != nil {
		return
	}

	expenses, err := u.expenseStorage.GetExpenses(ctx, userID, startDate)
	if err != nil {
		return
	}

	user, ok, err := u.userStorage.Get(ctx, userID)
	if err != nil {
		return
	}
	if !ok {
		err = set_currency.NewUserNotFoundErr(userID)
		return
	}

	sumByCategory, err := u.calculateSumByCategoryInUserCurrency(ctx, expenses, user.Currency)
	if err != nil {
		return
	}
	entries := mapToReportEntries(sumByCategory)
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Category < entries[j].Category })
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].TotalSum.GreaterThan(entries[j].TotalSum) })

	return &entities.Report{Cur: user.Currency, Entries: entries}, nil
}

func mapToReportEntries(sumByCategory map[string]decimal.Decimal) []entities.ReportEntry {
	entries := make([]entities.ReportEntry, len(sumByCategory))
	var i = 0
	for cat, sum := range sumByCategory {
		entries[i] = entities.ReportEntry{Category: cat, TotalSum: sum}
		i += 1
	}
	return entries
}

func getReportStartDate(period entities.ReportPeriod) (time.Time, error) {
	now := time.Now().UTC()
	switch period {
	case entities.ReportFor1Week:
		return now.AddDate(0, 0, -7), nil
	case entities.ReportFor1Month:
		return now.AddDate(0, -1, 0), nil
	case entities.ReportFor1Year:
		return now.AddDate(-1, 0, 0), nil
	default:
		return time.Now(), fmt.Errorf("getReportStartDate for period %v is not implemented", period)
	}
}

func (u *Usecase) calculateSumByCategoryInUserCurrency(
	ctx context.Context, expenses []entities.Expense, userCurrency entities.Currency,
) (sumByCategory map[string]decimal.Decimal, err error) {
	sumByCategory = make(map[string]decimal.Decimal)
	for _, e := range expenses {
		var convertedSum decimal.Decimal
		convertedSum, err = u.currencyConverter.FromBase(ctx, userCurrency, e.Sum, e.Date)
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

func NewUsecase(es expenseStorage, us userStorage, cc currencyConverter) *Usecase {
	return &Usecase{
		expenseStorage:    es,
		currencyConverter: cc,
		userStorage:       us,
	}
}
