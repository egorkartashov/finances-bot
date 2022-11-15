package get_report

import (
	"context"
	"time"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/reports"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type Usecase struct {
	userStorage     userStorage
	reportCache     reportCache
	reportRequester reportRequester
}

type ReportRequest struct {
	UserID   int64
	Period   entities.ReportPeriod
	Currency *entities.Currency
}

type ReportResponse struct {
	GeneratingInBackground bool
	CachedReport           *reports.FormattedReport
}

const defaultFormat = entities.ReportFormatMessage

func (u *Usecase) GenerateReport(ctx context.Context, req ReportRequest) (
	resp ReportResponse, err error,
) {
	user, ok, err := u.userStorage.Get(ctx, req.UserID)
	if err != nil {
		return
	}
	if !ok {
		err = users.NewUserNotFoundErr(req.UserID)
		return
	}

	newReportReq := toNewReportRequest(req, user)
	formattedReport, err := u.reportCache.Get(ctx, newReportReq)
	if err != nil {
		return ReportResponse{}, err
	}
	if formattedReport != nil {
		return ReportResponse{CachedReport: formattedReport, GeneratingInBackground: false}, nil
	}

	if err = u.reportRequester.Send(ctx, newReportReq); err != nil {
		return ReportResponse{}, err
	}
	return ReportResponse{GeneratingInBackground: true}, nil
}

func (u *Usecase) ReportFinished(ctx context.Context, report *reports.FormattedReport) error {
	return u.reportCache.Save(ctx, report)
}

func toNewReportRequest(req ReportRequest, user entities.User) *reports.NewReportRequest {
	newReportReq := &reports.NewReportRequest{
		UserID:   req.UserID,
		Currency: user.Currency,
		Period:   req.Period,
		Date:     time.Now(),
		Format:   defaultFormat,
	}
	if req.Currency != nil {
		newReportReq.Currency = *req.Currency
	}
	return newReportReq
}

func NewUsecase(us userStorage, rc reportCache, rr reportRequester) *Usecase {
	return &Usecase{
		userStorage:     us,
		reportCache:     rc,
		reportRequester: rr,
	}
}
