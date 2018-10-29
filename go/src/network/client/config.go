// SPDX License Identifier: MIT
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jynik/skullsup/go/src/file"
)

type Config struct {
	// Device to write to (only relevant for queue readers
	Device string `json:"device"`

	// Hostname or IP address of the SkullsUp! Server
	Host string `json:"host"`

	// Port number the server is running on
	Port uint16 `json:"port"`

	// Path to Certificate of CA that signed client cert
	CAPath string `json:"ca_path"`

	// Path to TLS Client Cerficate
	CertPath string `json:"cert_path"`

	// Path to client's private key
	KeyPath string `json:"key_path"`

	// Path to log file
	LogFilePath string `json:"log_path"`

	// Log level as one of: debug, info, error, silent
	LogLevel string `json:"log_level"`

	// Server polling period in seconds.
	PollPeriod int `json:"poll_period"`

	// Default frame period, in ms
	FramePeriod int `json:"frame_period"`

	// Queues available to read
	ReadQueues []string `json:"read_queues"`

	// Queues available to write
	WriteQueues []string `json:"write_queues"`

	// Disable verification of server certificate
	Insecure bool `json:"insecure"`
}

// Search default locations for a config
const FindDefaultConfig = ""

func loadConfig(filename string) (*Config, error) {

	config := new(Config)

	if filename == FindDefaultConfig {
		filename = "skullsup-client.cfg"
	}

	data, err := file.FindAndRead(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse configuration file: %s\n", err)
	}

	if len(config.Host) == 0 {
		return nil, errors.New("Configuration file is missing host definition.")
	}

	if config.Port == 0 {
		return nil, errors.New("Configuration file is missing valid port definition.")
	}

	config.CAPath = os.ExpandEnv(config.CAPath)
	if len(config.CAPath) == 0 {
		return nil, errors.New("Configuration file is missing ca_path definition.")
	}

	config.CertPath = os.ExpandEnv(config.CertPath)
	if len(config.CertPath) == 0 {
		return nil, errors.New("Configuration file is missing cert_path definition.")
	}

	config.KeyPath = os.ExpandEnv(config.KeyPath)
	if len(config.KeyPath) == 0 {
		return nil, errors.New("Configuration file is missing key_path definition.")
	}

	config.LogFilePath = os.ExpandEnv(config.LogFilePath)
	if len(config.LogFilePath) == 0 {
		config.LogFilePath = "stderr"
	}

	if config.PollPeriod <= 0 {
		config.PollPeriod = 15
	}

	if config.FramePeriod < 0 {
		// Use the incantation's default
		config.FramePeriod = 0
	}

	return config, err
}
