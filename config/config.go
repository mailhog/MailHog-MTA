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
	BindAddr  string          `json:",omitempty"`
	Hostname  string          `json:",omitempty"`
	PolicySet ServerPolicySet `json:",omitempty"`
	Backends  Backends        `json:",omitempty"`
	TLSConfig TLSConfig       `json:",omitempty"`
}

// TLSConfig holds a servers TLS config
type TLSConfig struct {
	CertFile string `json:",omitempty"`
	KeyFile  string `json:",omitempty"`
}

// ServerPolicySet defines the policies which can be applied per-server
type ServerPolicySet struct {
	// RequireAuthentication forces the server to require authentication before
	// any other commands (except STARTTLS) are accepted. Port 587 will typically
	// have this set to prevent abuse.
	RequireAuthentication bool
	// RequireLocalDelivery requires messages to be addressed to local domains
	// (primary or secondary local). E.g., port 25 will typically have this
	// set to avoid becoming an open relay.
	RequireLocalDelivery bool
	// MaximumRecipients is the maximum number of recipients accepted per-message.
	// Additional recipients will be rejected.
	MaximumRecipients int
	// DisableTLS disables the STARTTLS command.
	DisableTLS bool
	// RequireTLS requires all connections use TLS, disabling all commands except
	// STARTTLS until TLS negotiation is complete.
	RequireTLS bool
	// MaximumLineLength is the maximum length of a line in the SMTP conversation.
	MaximumLineLength int
	// MaximumConnections is the maximum number of concurrent connections the
	// server will accept.
	MaximumConnections int
	// RejectInvalidRecipients means invalid recipients at valid primary local domains
	// will be rejected at the 'RCPT TO' stage. The default behaviour is to accept
	// the message (and bounce it later) to minimise directory harvesting.
	RejectInvalidRecipients bool
}

// IdentityPolicySet defines the policies which can be applied per-user.
// If set, they take precedence over server policies
type IdentityPolicySet struct {
	RequireLocalDelivery    *bool
	MaximumRecipients       *int
	RejectInvalidRecipients *bool
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

// DefaultIdentityPolicySet defines a default policy set with no non-nil options
func DefaultIdentityPolicySet() IdentityPolicySet {
	return IdentityPolicySet{}
}

// DefaultSubmissionPolicySet defines the default ServerPolicySet for a submission server
func DefaultSubmissionPolicySet() ServerPolicySet {
	return ServerPolicySet{
		RequireAuthentication:   true,
		RequireLocalDelivery:    false,
		MaximumRecipients:       500,
		DisableTLS:              false,
		RequireTLS:              true,
		MaximumLineLength:       1024000,
		MaximumConnections:      1000,
		RejectInvalidRecipients: false,
	}
}

// DefaultSMTPPolicySet defines the default ServerPolicySet for an SMTP server
func DefaultSMTPPolicySet() ServerPolicySet {
	return ServerPolicySet{
		RequireAuthentication:   false,
		RequireLocalDelivery:    true,
		MaximumRecipients:       500,
		RequireTLS:              false,
		DisableTLS:              false,
		MaximumLineLength:       1024000,
		MaximumConnections:      1000,
		RejectInvalidRecipients: false,
	}
}

var cfg = DefaultConfig()

var configFile string

// Configure returns the configuration
func Configure() *Config {
	if len(configFile) > 0 {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Error reading %s: %s", configFile, err)
			os.Exit(1)
		}
		switch {
		case strings.HasSuffix(configFile, ".json"):
			err = json.Unmarshal(b, &cfg)
			if err != nil {
				fmt.Printf("Error parsing JSON in %s: %s", configFile, err)
				os.Exit(3)
			}
		default:
			fmt.Printf("Unsupported file type: %s\n", configFile)
			os.Exit(2)
		}

		cfg.relPath = filepath.Dir(configFile)
	}

	b, _ := json.MarshalIndent(&cfg, "", "  ")
	fmt.Println(string(b))

	return cfg
}

// RegisterFlags registers command line options
func RegisterFlags() {
	flag.StringVar(&configFile, "config-file", "", "Path to configuration file")
}
