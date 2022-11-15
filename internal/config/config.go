package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/entities"
	"gopkg.in/yaml.v3"
)

const configFile = "configs/config.yaml"

type SendReportConfig struct {
	Host         string `env:"SEND_REPORT_HOST,required"`
	Port         int    `env:"SEND_REPORT_PORT,required"`
	ClientSecret string `env:"GRPC_CLIENT_SECRET,required"`
}

type Config struct {
	Token                string            `env:"TG_TOKEN,required"`
	RateFetchFreqMinutes int               `yaml:"rateFetchFreqMinutes"`
	DbDsn                string            `env:"FINANCES_DSN,required"`
	BaseCurrency         entities.Currency `yaml:"baseCurrency"`
	ServiceName          string            `yaml:"serviceName"`
	CacheURL             string            `env:"CACHE_URL,required"`
	KafkaBrokers         []string          `env:"KAFKA_BROKERS,required" envSeparator:";"`
	SendReportGrpc       SendReportConfig
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	s := &Service{
		config: Config{
			ServiceName:    "finances-bot",
			SendReportGrpc: SendReportConfig{},
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
	return s.config.DbDsn
}

func (s *Service) BaseCurrency() entities.Currency {
	return s.config.BaseCurrency
}

func (s *Service) ServiceName() string {
	return s.config.ServiceName
}

func (s *Service) CacheURL() string {
	return s.config.CacheURL
}

func (s *Service) KafkaBrokers() []string {
	return s.config.KafkaBrokers
}

func (s *Service) SendReportAddr() string {
	return fmt.Sprintf("%s:%d", s.config.SendReportGrpc.Host, s.config.SendReportGrpc.Port)
}

func (s *Service) SendReportPort() int {
	return s.config.SendReportGrpc.Port
}

func (s *Service) SendReportClientSecret() string {
	return s.config.SendReportGrpc.ClientSecret
}
