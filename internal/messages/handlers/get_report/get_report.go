package get_report

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type ReportPresenter interface {
	ReportToPlainText(report *entities.Report) string
}

type Report struct {
	usecase   usecase
	presenter ReportPresenter
	sender    messages.MessageSender
}

func New(uc usecase, p ReportPresenter, s messages.MessageSender) *Report {
	return &Report{
		usecase:   uc,
		presenter: p,
		sender:    s,
	}
}

const (
	ReportFormatMessage    = "отчет <период>, где период может быть одним из значений: неделя, месяц, год"
	ReportHelp             = "Чтобы получить отчет, введи команду: " + ReportFormatMessage
	IncorrectFormatMessage = "Неизвестный период отчета. Корректный формат: " + ReportFormatMessage
)

var (
	reportKeywords  = []string{"отчет", "отчёт"}
	keywordToPeriod = map[string]entities.ReportPeriod{
		"неделя": entities.ReportFor1Week,
		"месяц":  entities.ReportFor1Month,
		"год":    entities.ReportFor1Year,
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

	span, ctx := opentracing.StartSpanFromContext(ctx, "report")
	defer span.Finish()

	params := strings.TrimPrefix(msg.Text, foundKw)
	params = strings.Trim(params, " ")
	reportPeriod, ok := keywordToPeriod[params]
	span.SetTag("report-period", reportPeriod)
	if !ok {
		err := h.sender.SendText(IncorrectFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}

	report, err := h.usecase.GenerateReport(ctx, msg.UserID, reportPeriod)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	reportStr := h.presenter.ReportToPlainText(report)
	err = h.sender.SendText(reportStr, msg.UserID)
	return utils.HandleWithErrorOrNil(err)
}

func (h *Report) Name() string {
	return "Report"
}
