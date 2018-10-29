// SPDX License Identifier: MIT
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jynik/skullsup/go/src/logger"
	"github.com/jynik/skullsup/go/src/network"
)

type Client struct {
	Cfg *Config
	Log *logger.Logger
	http.Client
}

func (c *Client) configureTransport() error {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(c.Cfg.CertPath, c.Cfg.KeyPath)
	if err != nil {
		return err
	}

	c.Log.Debug("Loaded Client Certificate: %s\n", c.Cfg.CertPath)
	c.Log.Debug("Loaded Client Key: %s\n", c.Cfg.KeyPath)

	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: c.Cfg.Insecure,
		},
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if c.Cfg.Insecure {
		c.Log.Info("Operating insecurely! Skipping server certificate verification.\n")
	}

	return nil
}

// Create a new SkullsUp! Client, provided a path to a configuration file
func New(filename string) (*Client, error) {
	var err error

	c := new(Client)
	c.Cfg, err = loadConfig(filename)
	if err != nil {
		return nil, err
	}

	c.Log, err = logger.New(c.Cfg.LogFilePath, c.Cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	if err := c.configureTransport(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) canAccess(queue string, queues []string) bool {
	for _, q := range queues {
		if queue == q {
			return true
		}
	}
	return false
}

func (c *Client) CanRead(queue string) bool {
	return c.canAccess(queue, c.Cfg.ReadQueues)
}

func (c *Client) CanWrite(queue string) bool {
	return c.canAccess(queue, c.Cfg.WriteQueues)
}

func responseAsError(r *http.Response) error {
	var statusText string
	if r.ContentLength >= 1 {
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)

		// Trim newline
		if buf[len(buf)-1] == '\n' {
			buf = buf[:len(buf)-1]
		}

		statusText = string(buf)

	} else {
		statusText = http.StatusText(r.StatusCode)
	}

	return fmt.Errorf("Received status %d: %s", r.StatusCode, statusText)
}

// Read a messsage from the specified queue
func (c *Client) Read(queue string) (*network.Message, error) {
	var msg network.Message

	if !c.CanRead(queue) {
		return nil, fmt.Errorf("Client is not configured to read from \"%s\"", queue)
	}

	url := network.QueueUrl(c.Cfg.Host, c.Cfg.Port, queue)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, responseAsError(resp)
	} else if resp.ContentLength < 1 {
		return nil, errors.New("Received empty response")
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

// Write a message to the specified queue
func (c *Client) Write(msg *network.Message, queue string) error {
	msg.Command = strings.ToLower(msg.Command)
	switch msg.Command {
	case network.CmdColor:
	case network.CmdIncant:
	case network.CmdReanimate:
	default:
		return fmt.Errorf("Invalid command: %s", msg.Command)
	}

	body, err := json.Marshal(*msg)
	if err != nil {
		return fmt.Errorf("Failed to construct message: %s\n", err)
	}

	url := network.QueueUrl(c.Cfg.Host, c.Cfg.Port, queue)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to construct request: %s\n", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		if resp.ContentLength >= 1 {
			defer resp.Body.Close()
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if len(buf) > 1 {
				// Trim newline
				buf = buf[:len(buf)-1]
				return fmt.Errorf("Status %d: %s", resp.StatusCode, buf)
			}
		}

		// At least return the status code if we're lost here...
		return fmt.Errorf("¯\\_(ツ)_/¯  Status %d", resp.StatusCode)
	}

	return nil
}
