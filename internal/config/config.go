package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type Config struct {
	Token                string `yaml:"token"`
	RateFetchFreqMinutes int    `yaml:"rateFetchFreqMinutes"`
	Dsn                  string `env:"FINANCES_DSN"`
	BaseCurrency         entities.Currency
}

type Service struct {
	config Config
}

func New(baseCurr entities.Currency) (*Service, error) {
	s := &Service{
		config: Config{
			BaseCurrency: baseCurr,
		},
	}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing yaml")
	}

	err = env.Parse(&s.config)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing env")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) RateFetchFreqMinutes() int {
	return s.config.RateFetchFreqMinutes
}

func (s *Service) Dsn() string {
	return s.config.Dsn
}

func (s *Service) BaseCurrency() entities.Currency {
	return s.config.BaseCurrency
}
