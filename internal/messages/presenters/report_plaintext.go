package presenters

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
)

type ReportPresenter struct {
}

func NewReportPresenter() *ReportPresenter {
	return &ReportPresenter{}
}

func (p *ReportPresenter) ReportToPlainText(report *expenses.Report) string {
	if len(report.Entries) == 0 {
		return "Трат нет"
	}

	var sb strings.Builder
	for i, e := range report.Entries {
		line := fmt.Sprintf("%v. %s: %v руб.\n", i+1, e.Category, e.TotalSumKop)
		sb.WriteString(line)
	}

	var reportStr = sb.String()
	reportStr = strings.Trim(reportStr, "\n")
	return reportStr
}
