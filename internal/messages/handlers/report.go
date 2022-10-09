package handlers

import (
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
)

type ReportPresenter interface {
	ReportToPlainText(report *expenses.Report) string
}

type Report struct {
	expenses  *expenses.Model
	presenter ReportPresenter
	base
}

func NewReport(expenses *expenses.Model, presenter ReportPresenter, sender messages.MessageSender) *Report {
	return &Report{
		expenses:  expenses,
		presenter: presenter,
		base:      base{sender},
	}
}

const (
	ReportFormatMessage    = "отчет <период>, где период может быть одним из значений: неделя, месяц, год"
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

func (h *Report) Handle(msg messages.Message) messages.MessageHandleResult {
	foundKw := ""
	for _, kw := range reportKeywords {
		if strings.HasPrefix(msg.Text, kw) {
			foundKw = kw
		}
	}

	if foundKw == "" {
		return messages.MessageHandleResult{Skipped: true, Err: nil}
	}

	params := strings.TrimPrefix(msg.Text, foundKw)
	params = strings.Trim(params, " ")
	reportPeriod, ok := keywordToPeriod[params]
	if !ok {
		err := h.messageSender.SendMessage(IncorrectFormatMessage, msg.UserID)
		return messages.MessageHandleResult{Skipped: false, Err: err}
	}

	report, err := h.expenses.GenerateReport(msg.UserID, reportPeriod)
	if err != nil {
		err := h.messageSender.SendMessage("Что-то пошло не так", msg.UserID)
		return messages.MessageHandleResult{Skipped: false, Err: err}
	}

	reportStr := h.presenter.ReportToPlainText(report)
	err = h.messageSender.SendMessage(reportStr, msg.UserID)
	return messages.MessageHandleResult{Skipped: false, Err: err}
}

func (h *Report) Name() string {
	return "ReportHandler"
}
