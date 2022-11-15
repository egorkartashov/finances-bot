package generator

import (
	"context"
	"sort"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
)

type Generator struct {
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	presenters        map[entities.ReportFormat]ReportPresenter
}

func New(es expenseStorage, cc currencyConverter, presenters []ReportPresenter) *Generator {
	return &Generator{
		expenseStorage:    es,
		currencyConverter: cc,
		presenters:        mustMakeMap(presenters),
	}
}

func mustMakeMap(presenters []ReportPresenter) map[entities.ReportFormat]ReportPresenter {
	m := make(map[entities.ReportFormat]ReportPresenter)
	for _, presenter := range presenters {
		format := presenter.Format()
		if _, ok := m[format]; ok {
			logger.Fatal("duplicated presenter: " + string(format))
		}
		m[format] = presenter
	}
	return m
}

func (g *Generator) Generate(ctx context.Context, req *reports.NewReportRequest) (
	*reports.GeneratedReportResponse, error,
) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "generate-report")
	defer span.Finish()

	presenter, ok := g.presenters[req.Format]
	if !ok {
		return nil, reports.NewErrUnsupportedFormat(string(req.Format))
	}

	startDate, err := req.Period.GetStartDate(time.Now().UTC())
	if err != nil {
		return nil, err
	}

	expenses, err := g.expenseStorage.GetExpenses(ctx, req.UserID, startDate)
	if err != nil {
		return nil, err
	}

	sumByCategory, err := g.calculateSumByCategoryInUserCurrency(ctx, expenses, req.Currency)
	if err != nil {
		return nil, err
	}
	entries := mapToReportEntries(sumByCategory)
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Category < entries[j].Category })
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].TotalSum.GreaterThan(entries[j].TotalSum) })

	report := &entities.Report{Cur: req.Currency, Entries: entries}

	payload, err := presenter.Present(ctx, report)
	if err != nil {
		return nil, err
	}

	return &reports.GeneratedReportResponse{Payload: payload}, nil
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

func (g *Generator) calculateSumByCategoryInUserCurrency(
	ctx context.Context, expenses []entities.Expense, userCurrency entities.Currency,
) (sumByCategory map[string]decimal.Decimal, err error) {
	sumByCategory = make(map[string]decimal.Decimal)
	for _, e := range expenses {
		var convertedSum decimal.Decimal
		convertedSum, err = g.currencyConverter.FromBase(ctx, userCurrency, e.Sum, e.Date)
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
