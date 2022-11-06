package presenters

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
)

type Report struct {
}

func NewReport() *Report {
	return &Report{}
}

func (p *Report) ReportToPlainText(report *entities.Report) string {
	if len(report.Entries) == 0 {
		return "Трат нет"
	}

	var sb strings.Builder
	for i, e := range report.Entries {
		roundedSum := e.TotalSum.Round(2)
		line := fmt.Sprintf("%v. %s: %v %s \n", i+1, e.Category, roundedSum, report.Cur)
		sb.WriteString(line)
	}

	var reportStr = sb.String()
	reportStr = strings.Trim(reportStr, "\n")
	return reportStr
}
