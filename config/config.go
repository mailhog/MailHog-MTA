package config

func DefaultConfig() *Config {
	return &Config{
		Servers: []*Server{
			&Server{
				BindAddr:  "0.0.0.0:25",
				Hostname:  "mailhog.example",
				PolicySet: DefaultSMTPPolicySet(),
			},
			&Server{
				BindAddr:  "0.0.0.0:587",
				Hostname:  "mailhog.example",
				PolicySet: DefaultSubmissionPolicySet(),
			},
		},
	}
}

type Config struct {
	Servers []*Server
}

type Server struct {
	BindAddr  string
	Hostname  string
	PolicySet PolicySet
}

type PolicySet struct {
	RequireAuthentication bool
	RequireLocalDelivery  bool
	MaximumRecipients     int
	RequireTLS            bool
}

func DefaultSubmissionPolicySet() PolicySet {
	return PolicySet{
		RequireAuthentication: true,
		RequireLocalDelivery:  false,
		MaximumRecipients:     500,
		RequireTLS:            true,
	}
}

func DefaultSMTPPolicySet() PolicySet {
	return PolicySet{
		RequireAuthentication: false,
		RequireLocalDelivery:  true,
		MaximumRecipients:     500,
		RequireTLS:            false,
	}
}

var cfg = DefaultConfig()

func Configure() *Config {
	return cfg
}

func RegisterFlags() {
	// TODO
}
