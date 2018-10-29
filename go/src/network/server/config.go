// SPDX License Identifier: MIT
package server

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jynik/skullsup/go/src/file"
	"github.com/jynik/skullsup/go/src/network"
)

type Config struct {
	// Port number the server listens on
	Port uint16 `json:"port"`

	// Address to bind the server listening socket to
	Address string `json:"addr"`

	// Path to server's TLS certificate (public PEM)
	CertPath string `json:"cert_path"`

	// Path to private key corresponding to TLS certificate
	KeyPath string `json:"key_path"`

	// Path to the trusted root CA certificate for client-provided certificates
	// We will trust clients whose certificates are verified against this CA
	// and whose certificate hashes are not in the revocation blacklist.
	CAPath string `json:"ca_path"`

	// Path to log file. May be a a filename, "stdout", or "stderr"
	LogPath string `json:"log_path"`

	// Log level as one of: debug, info, error, silent
	LogLevel string `json:"log_level"`

	// Blacklisted client certificates, identified by hex-encoded serial number
	// strings This is intended to serve as a quick and dirty revocation list.
	Blacklist []string `json:"blacklist"`

	// User configuration
	Users network.UserList `json:"users"`
}

// Indicates that the config file search path should be found
const FindDefaultConfig = ""

func validateSerial(serial *string) error {
	if len(*serial) == 0 {
		return fmt.Errorf("Encountered empty serial number string")
	}

	if _, err := hex.DecodeString(*serial); err != nil {
		return fmt.Errorf("Invalid hex string: \"%s\"", serial)
	} else {
		*serial = strings.ToLower(*serial)
	}
	return nil
}

// Load a server configuration from the specified filename.
// Returns a valid *Config on success, and a non-nil error otherwise.
func loadConfig(filename string) (*Config, error) {
	cfg := new(Config)

	if filename == FindDefaultConfig {
		filename = "skullsup-queue-server.cfg"
	}

	data, err := file.FindAndRead(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.CertPath = os.ExpandEnv(cfg.CertPath)
	cfg.KeyPath = os.ExpandEnv(cfg.KeyPath)
	cfg.CAPath = os.ExpandEnv(cfg.CAPath)
	cfg.LogPath = os.ExpandEnv(cfg.LogPath)

	for i := range cfg.Blacklist {
		if err := validateSerial(&cfg.Blacklist[i]); err != nil {
			return nil, err
		}
	}

	// Ensure there are no duplicate or blacklisted serial numbers in the users list
	serialMap := map[string]bool{}
	nameMap := map[string]bool{}
	for i := range cfg.Users {
		name := cfg.Users[i].Name
		if _, exists := nameMap[name]; exists {
			return nil, fmt.Errorf("Duplicate user name detected: %s", name)
		}
		nameMap[name] = true

		if err := validateSerial(&cfg.Users[i].CertSerial); err != nil {
			return nil, err
		}

		serial := cfg.Users[i].CertSerial

		if _, exists := serialMap[serial]; exists {
			return nil, fmt.Errorf("Duplicate certificate serial number detected: %s", serial)
		}
		serialMap[serial] = true

		for _, blacklisted := range cfg.Blacklist {
			if blacklisted == serial {
				return nil, fmt.Errorf("Certificate for user \"%s\" is blacklisted: %s", cfg.Users[i].Name, serial)
			}
		}
	}

	return cfg, nil
}

// Look up a user based upon their certificate.
//
// It is assumed that this certificate has already been validated and
// that is is signed by a trusted CA!
//
func (c *Config) LookupUser(cert *x509.Certificate) (*network.User, error) {

	serial := hex.EncodeToString(cert.SerialNumber.Bytes())

	for _, b := range c.Blacklist {
		if b == serial {
			return nil, fmt.Errorf("Certificate is blacklisted (serial=%s)", serial)
		}
	}

	for i, u := range c.Users {
		if u.CertSerial == serial && u.Name == cert.Subject.CommonName {
			return &c.Users[i], nil
		}
	}

	return nil, nil
}
