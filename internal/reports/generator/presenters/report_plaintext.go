package presenters

import (
	"context"
	"fmt"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type FormatMessage struct {
}

func NewFormatMessage() *FormatMessage {
	return &FormatMessage{}
}

func (p *FormatMessage) Format() entities.ReportFormat {
	return entities.ReportFormatMessage
}

func (p *FormatMessage) Present(_ context.Context, report *entities.Report) (string, error) {
	if len(report.Entries) == 0 {
		return "Трат нет", nil
	}

	var sb strings.Builder
	for i, e := range report.Entries {
		roundedSum := e.TotalSum.Round(2)
		line := fmt.Sprintf("%v. %s: %v %s \n", i+1, e.Category, roundedSum, report.Cur)
		sb.WriteString(line)
	}

	var reportStr = sb.String()
	reportStr = strings.Trim(reportStr, "\n")
	return reportStr, nil
}
