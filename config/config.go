package config

import (
	"flag"

	"github.com/ian-kent/envconf"
)

func DefaultConfig() *Config {
	return &Config{
		SMTPBindAddr: "0.0.0.0:25",
		Hostname:     "mailhog.example",
	}
}

type Config struct {
	SMTPBindAddr string
	Hostname     string
}

var cfg = DefaultConfig()

func Configure() *Config {
	return cfg
}

func RegisterFlags() {
	flag.StringVar(&cfg.SMTPBindAddr, "smtpbindaddr", envconf.FromEnvP("MHMTA_SMTP_BIND_ADDR", "0.0.0.0:25").(string), "SMTP bind interface and port, e.g. 0.0.0.0:25 or just :25")
	flag.StringVar(&cfg.Hostname, "hostname", envconf.FromEnvP("MHMTA_HOSTNAME", "mailhog.example").(string), "Hostname for EHLO/HELO response, e.g. mailhog.example")
}
