package handlers

import (
	"context"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type ReportPresenter interface {
	ReportToPlainText(report *expenses.Report) string
}

type Report struct {
	expenses  *expenses.Usecase
	presenter ReportPresenter
	base
}

func NewReport(expenses *expenses.Usecase, presenter ReportPresenter, sender messages.MessageSender) *Report {
	return &Report{
		expenses:  expenses,
		presenter: presenter,
		base:      base{sender},
	}
}

const (
	ReportFormatMessage    = "отчет <период>, где период может быть одним из значений: неделя, месяц, год"
	ReportHelp             = "Чтобы получить отчет, введи команду: " + ReportFormatMessage
	IncorrectFormatMessage = "Неизвестный период отчета. Корректный формат: " + ReportFormatMessage
)

var (
	reportKeywords  = []string{"отчет", "отчёт"}
	keywordToPeriod = map[string]expenses.ReportPeriod{
		"неделя": expenses.ReportFor1Week,
		"месяц":  expenses.ReportFor1Month,
		"год":    expenses.ReportFor1Year,
	}
)

func (h *Report) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	foundKw := ""
	for _, kw := range reportKeywords {
		if strings.HasPrefix(msg.Text, kw) {
			foundKw = kw
		}
	}

	if foundKw == "" {
		return utils.HandleSkipped
	}

	params := strings.TrimPrefix(msg.Text, foundKw)
	params = strings.Trim(params, " ")
	reportPeriod, ok := keywordToPeriod[params]
	if !ok {
		err := h.MessageSender.SendText(IncorrectFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}

	report, err := h.expenses.GenerateReport(ctx, msg.UserID, reportPeriod)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	reportStr := h.presenter.ReportToPlainText(report)
	err = h.MessageSender.SendText(reportStr, msg.UserID)
	return utils.HandleWithErrorOrNil(err)
}

func (h *Report) Name() string {
	return "Report"
}
