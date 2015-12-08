package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TODO: make TLSConfig and PolicySet 'ref'able

// DefaultConfig provides a default (but relatively useless) configuration
func DefaultConfig() *Config {
	return &Config{
		Backends: map[string]BackendConfig{
			"local_auth": BackendConfig{
				Type: "local",
				Data: map[string]interface{}{
					"config": "auth.json",
				},
			},
			"local_resolver": BackendConfig{
				Type: "local",
				Data: map[string]interface{}{
					"config": "resolve.json",
				},
			},
			"local_delivery": BackendConfig{
				Type: "local",
				Data: map[string]interface{}{},
			},
		},
		Servers: []*Server{
			&Server{
				BindAddr:  "0.0.0.0:25",
				Hostname:  "mailhog.example",
				PolicySet: DefaultSMTPPolicySet(),
				Backends: Backends{
					Auth: &BackendConfig{
						Ref: "local_auth",
					},
					Resolver: &BackendConfig{
						Ref: "local_resolver",
					},
					Delivery: &BackendConfig{
						Ref: "local_delivery",
					},
				},
			},
			&Server{
				BindAddr:  "0.0.0.0:587",
				Hostname:  "mailhog.example",
				PolicySet: DefaultSubmissionPolicySet(),
				Backends: Backends{
					Auth: &BackendConfig{
						Ref: "local_auth",
					},
					Resolver: &BackendConfig{
						Ref: "local_resolver",
					},
					Delivery: &BackendConfig{
						Ref: "local_delivery",
					},
				},
			},
		},
	}
}

// Config defines the top-level application configuration
type Config struct {
	relPath string

	Servers  []*Server                `json:",omitempty"`
	Backends map[string]BackendConfig `json:",omitempty"`
}

// RelPath returns the path to the configuration file directory,
// used when loading external files using relative paths
func (c Config) RelPath() string {
	return c.relPath
}

// Server defines the configuration of an individual bind address
type Server struct {
	BindAddr  string    `json:",omitempty"`
	Hostname  string    `json:",omitempty"`
	PolicySet PolicySet `json:",omitempty"`
	Backends  Backends  `json:",omitempty"`
}

// PolicySet defines the policies which can be applied per-server
type PolicySet struct {
	RequireAuthentication bool
	RequireLocalDelivery  bool
	MaximumRecipients     int
	EnableTLS             bool
	RequireTLS            bool
	MaximumLineLength     int
	MaximumConnections    int
}

// Backends defines the backend configurations for a server
type Backends struct {
	Auth     *BackendConfig `json:",omitempty"`
	Resolver *BackendConfig `json:",omitempty"`
	Delivery *BackendConfig `json:",omitempty"`
}

// BackendConfig defines an individual backend configuration
type BackendConfig struct {
	Type string                 `json:",omitempty"`
	Ref  string                 `json:",omitempty"`
	Data map[string]interface{} `json:",omitempty"`
}

// DefaultSubmissionPolicySet defines the default PolicySet for a submission server
func DefaultSubmissionPolicySet() PolicySet {
	return PolicySet{
		RequireAuthentication: true,
		RequireLocalDelivery:  false,
		MaximumRecipients:     500,
		RequireTLS:            true,
		EnableTLS:             true,
		MaximumLineLength:     1024000,
		MaximumConnections:    1000,
	}
}

// DefaultSMTPPolicySet defines the default PolicySet for an SMTP server
func DefaultSMTPPolicySet() PolicySet {
	return PolicySet{
		RequireAuthentication: false,
		RequireLocalDelivery:  true,
		MaximumRecipients:     500,
		RequireTLS:            false,
		EnableTLS:             true,
		MaximumLineLength:     1024000,
		MaximumConnections:    1000,
	}
}

var cfg = DefaultConfig()

var configPath string

// Configure returns the configuration
func Configure() *Config {
	if len(configPath) > 0 {
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading %s: %s", configPath, err)
			os.Exit(1)
		}
		switch {
		case strings.HasSuffix(configPath, ".json"):
			err = json.Unmarshal(b, &cfg)
			if err != nil {
				fmt.Printf("Error parsing JSON in %s: %s", configPath, err)
				os.Exit(3)
			}
		default:
			fmt.Printf("Unsupported file type: %s\n", configPath)
			os.Exit(2)
		}

		cfg.relPath = filepath.Dir(configPath)
	}

	return cfg
}

// RegisterFlags registers command line options
func RegisterFlags() {
	flag.StringVar(&configPath, "config-path", "", "Path to configuration file")
}
