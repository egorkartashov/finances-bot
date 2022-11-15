package get_report

import (
	"context"
	"errors"
	"strings"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/usecases/get_report"
)

type ReportPresenter interface {
	ReportToPlainText(report *entities.Report) string
}

type GetReport struct {
	usecase usecase
	sender  messages.MessageSender
}

func New(uc usecase, s messages.MessageSender) *GetReport {
	return &GetReport{
		usecase: uc,
		sender:  s,
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

func (g *GetReport) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
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
		err := g.sender.SendText(IncorrectFormatMessage, msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}

	req := get_report.ReportRequest{
		UserID: msg.UserID,
		Period: reportPeriod,
	}
	response, err := g.usecase.GenerateReport(ctx, req)
	if err != nil {
		return utils.HandleWithErrorOrNil(err)
	}

	if !response.GeneratingInBackground {
		err = g.sendFormattedReport(response.CachedReport)
		return utils.HandleWithErrorOrNil(err)
	} else {
		err = g.sender.SendText("Начинаем построение отчета...", msg.UserID)
		return utils.HandleWithErrorOrNil(err)
	}
}

func (g *GetReport) Name() string {
	return "GetReport"
}

func (g *GetReport) SendFinishedReport(ctx context.Context, report *reports.FormattedReport) error {
	if err := g.usecase.ReportFinished(ctx, report); err != nil {
		return err
	}
	return g.sendFormattedReport(report)
}

func (g *GetReport) sendFormattedReport(report *reports.FormattedReport) error {
	msg, err := g.toBotMessage(report)
	if err != nil {
		return err
	}
	return g.sender.SendMessage(report.UserID, msg)
}

func (g *GetReport) toBotMessage(report *reports.FormattedReport) (*messages.Message, error) {
	switch report.Format {
	case entities.ReportFormatMessage:
		return &messages.Message{Text: report.Payload}, nil
	default:
		return nil, errors.New("unsupported report format: " + string(report.Format))
	}
}
