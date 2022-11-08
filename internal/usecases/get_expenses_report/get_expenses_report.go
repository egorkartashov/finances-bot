package get_expenses_report

import (
	"context"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type Usecase struct {
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userStorage       userStorage
	reportCache       reportCache
}

func (u *Usecase) GenerateReport(ctx context.Context, userID int64, reportPeriod entities.ReportPeriod) (
	report *entities.Report, err error,
) {
	report, err = u.reportCache.Get(ctx, userID, reportPeriod)
	if err != nil {
		return nil, err
	}
	if report != nil {
		return report, nil
	}

	startDate, err := reportPeriod.GetStartDate(time.Now().UTC())
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
		err = users.NewUserNotFoundErr(userID)
		return
	}

	sumByCategory, err := u.calculateSumByCategoryInUserCurrency(ctx, expenses, user.Currency)
	if err != nil {
		return
	}
	entries := mapToReportEntries(sumByCategory)
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Category < entries[j].Category })
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].TotalSum.GreaterThan(entries[j].TotalSum) })

	report = &entities.Report{Cur: user.Currency, Entries: entries}
	err = u.reportCache.Save(ctx, userID, reportPeriod, report)
	if err != nil {
		return nil, err
	}
	return report, nil
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

func NewUsecase(es expenseStorage, us userStorage, cc currencyConverter, rc reportCache) *Usecase {
	return &Usecase{
		expenseStorage:    es,
		currencyConverter: cc,
		userStorage:       us,
		reportCache:       rc,
	}
}
