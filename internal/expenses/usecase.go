//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_mocks -destination mocks/${GOPACKAGE}.go
package expenses

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type (
	expenseStorage interface {
		AddExpense(userID int64, exp Expense) error
		GetExpenses(userID int64, minTime time.Time) ([]Expense, error)
	}
	userUc interface {
		GetOrRegister(id int64) (users.User, error)
	}
	currencyConverter interface {
		Convert(sum decimal.Decimal, from, to currency.Currency, date time.Time) (decimal.Decimal, error)
	}
)

type Usecase struct {
	l                 *log.Logger
	baseCurr          currency.Currency
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userUc            userUc
}

func NewUsecase(
	expenseStorage expenseStorage,
	baseCurr currency.Currency,
	userUc userUc,
	currencyProvider currencyConverter,
	l *log.Logger,
) *Usecase {
	return &Usecase{
		l:                 l,
		baseCurr:          baseCurr,
		expenseStorage:    expenseStorage,
		userUc:            userUc,
		currencyConverter: currencyProvider,
	}
}

type Expense struct {
	Category string
	Sum      decimal.Decimal
	Date     time.Time
}

type AddExpenseDto struct {
	Category string
	Sum      int32
	Date     time.Time
}

type AddExpenseResult struct {
	Category string
	Sum      int32
	Cur      currency.Currency
	Date     time.Time
}

func (u *Usecase) AddExpense(userID int64, exp AddExpenseDto) (res AddExpenseResult, err error) {
	user, err := u.userUc.GetOrRegister(userID)
	if err != nil {
		u.l.Println(fmt.Errorf("AddExpense: %w", err))
		return
	}

	sumDecimal := decimal.NewFromInt32(exp.Sum)
	subRub, err := u.currencyConverter.Convert(sumDecimal, user.Currency, u.baseCurr, exp.Date)
	if err != nil {
		return
	}

	storedExp := Expense{
		Category: exp.Category,
		Sum:      subRub,
		Date:     exp.Date,
	}
	if err = u.expenseStorage.AddExpense(userID, storedExp); err != nil {
		return
	}

	res = AddExpenseResult{
		Category: exp.Category,
		Sum:      exp.Sum,
		Cur:      user.Currency,
		Date:     exp.Date,
	}
	return
}

type ReportPeriod byte

const (
	ReportFor1Week  ReportPeriod = 1
	ReportFor1Month ReportPeriod = 2
	ReportFor1Year  ReportPeriod = 3
)

type Report struct {
	Cur     currency.Currency
	Entries []ReportEntry
}

type ReportEntry struct {
	Category string
	TotalSum decimal.Decimal
}

func (u *Usecase) GenerateReport(userID int64, reportPeriod ReportPeriod) (report *Report, err error) {
	minTime, err := getReportMinTime(reportPeriod)
	if err != nil {
		return
	}

	expenses, err := u.expenseStorage.GetExpenses(userID, minTime)
	if err != nil {
		return
	}

	user, err := u.userUc.GetOrRegister(userID)
	if err != nil {
		return
	}

	sumByCategory, err := u.getSumByCategory(expenses, user)
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

func (u *Usecase) getSumByCategory(expenses []Expense, user users.User) (sumByCategory map[string]decimal.Decimal, err error) {
	sumByCategory = make(map[string]decimal.Decimal)
	for _, e := range expenses {
		var convertedSum decimal.Decimal
		convertedSum, err = u.currencyConverter.Convert(e.Sum, u.baseCurr, user.Currency, e.Date)
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
