package config

func DefaultConfig() *Config {
	return &Config{
		Servers: []*Server{
			&Server{
				BindAddr: "0.0.0.0:25",
				Hostname: "mailhog.example",
			},
			&Server{
				BindAddr: "0.0.0.0:587",
				Hostname: "mailhog.example",
			},
		},
	}
}

type Config struct {
	Servers []*Server
}

type Server struct {
	BindAddr string
	Hostname string
	// PolicySet
}

var cfg = DefaultConfig()

func Configure() *Config {
	return cfg
}

func RegisterFlags() {
	// TODO
}
