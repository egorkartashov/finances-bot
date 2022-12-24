package cbrf

import (
	"context"
	"net/http"
	"time"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
)

type RatesApi struct {
}

func (r *RatesApi) FetchRatesToRub(ctx context.Context, currencies []entities.Currency, at time.Time) (
	[]entities.Rate, error,
) {
	dateStr := at.Format("02/01/2006")
	url := "https://www.cbr.ru/scripts/XML_daily.asp?date_req=" + dateStr

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	rates, err := Parse(res.Body, currencies, at)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
