package cbrf

import (
	"encoding/xml"
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/currency"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/entities"
	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

func Parse(xmlStream io.Reader, currencies []entities.Currency, at time.Time) ([]entities.Rate, error) {
	currSet := make(map[string]struct{})
	for _, cur := range currencies {
		currSet[string(cur)] = struct{}{}
	}

	var valCurs ValCurs
	dec := xml.NewDecoder(xmlStream)
	dec.CharsetReader = charset.NewReaderLabel
	if err := dec.Decode(&valCurs); err != nil {
		return nil, err
	}

	rates := make([]entities.Rate, 0, len(currencies))
	for _, valute := range valCurs.Valutes {
		if _, ok := currSet[valute.CharCode]; !ok {
			continue
		}

		valueStr := strings.Replace(valute.Value, ",", ".", 1)
		value, err := decimal.NewFromString(valueStr)
		if err != nil {
			return nil, err
		}

		rate := entities.Rate{
			From:  entities.Currency(valute.CharCode),
			To:    currency.RUB,
			Value: value,
			Date:  at,
		}
		rates = append(rates, rate)
	}

	return rates, nil
}
