package config

import (
	"flag"

	"github.com/ian-kent/envconf"
)

type DeliveryMode int

const (
	NoDelivery = DeliveryMode(iota)
	LocalDelivery
	RelayDelivery
)

type DeliveryPolicy struct {
	DeliveryMode DeliveryMode
}

type MTAConfig struct {
	gofigure interface{} `envPrefix:"MHMTA"`

	NoAuthPolicy DeliveryPolicy
	AuthPolicy   DeliveryPolicy
}

func DefaultConfig() *Config {
	if mtaCfg == nil {
		mtaCfg = &MTAConfig{
			NoAuthPolicy: DeliveryPolicy{LocalDelivery},
			AuthPolicy:   DeliveryPolicy{RelayDelivery},
		}
		//gofigure.Gofigure(mtaCfg)
	}

	return &Config{
		SMTPBindAddr: "0.0.0.0:25",
		Hostname:     "mailhog.example",
		MTAConfig:    mtaCfg,
	}
}

type Config struct {
	SMTPBindAddr string
	Hostname     string
	MTAConfig    *MTAConfig
}

var cfg = DefaultConfig()
var mtaCfg *MTAConfig

func Configure() *Config {
	return cfg
}

func RegisterFlags() {
	flag.StringVar(&cfg.SMTPBindAddr, "smtpbindaddr", envconf.FromEnvP("MHMTA_SMTP_BIND_ADDR", "0.0.0.0:25").(string), "SMTP bind interface and port, e.g. 0.0.0.0:25 or just :25")
	flag.StringVar(&cfg.Hostname, "hostname", envconf.FromEnvP("MHMTA_HOSTNAME", "mailhog.example").(string), "Hostname for EHLO/HELO response, e.g. mailhog.example")
}
