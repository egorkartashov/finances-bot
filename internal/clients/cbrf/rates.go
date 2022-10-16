package cbrf

import (
	"net/http"
	"time"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
)

type RatesApi struct {
}

func (r *RatesApi) FetchRatesToRub(currencies []currency.Currency, at time.Time) ([]currency.Rate, error) {
	dateStr := at.Format("02/01/2006")
	url := "https://www.cbr.ru/scripts/XML_daily.asp?date_req=" + dateStr

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	rates, err := Parse(res.Body, currencies, at)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
