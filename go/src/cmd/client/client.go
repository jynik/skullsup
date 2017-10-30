// SPDX License Identifier: MIT
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	c "../common"
	"../common/defaults"

	"../../skullsup"
)

const DummyConfig = defaults.DUMMY_PREFIX + defaults.WRITER_CONFIG

type Client struct {
	Host       string // Queue Server hostname or IP
	Port       uint16 // Queue Server port
	ConfigPath string // Path to client configuration

	httpClient http.Client
}

func (c *Client) disableCertVerification(*kingpin.ParseContext) error {
	c.httpClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return nil
}

func printVersion(*kingpin.ParseContext) error {
	fmt.Println(skullsup.Version())
	os.Exit(0)
	return nil
}

func New() *Client {
	var client Client

	kingpin.Flag(c.FLAG_HOST, c.FLAG_HOST_DESC).
		Short(c.FLAG_HOST_SHORT).
		Default(defaults.REMOTE).
		StringVar(&client.Host)

	kingpin.Flag(c.FLAG_PORT, c.FLAG_PORT_DESC).
		Short(c.FLAG_PORT_SHORT).
		Default(strconv.Itoa(defaults.PORT)).
		Uint16Var(&client.Port)

	kingpin.Flag(c.FLAG_CLIENT_CONFIG, c.FLAG_CLIENT_CONFIG_DESC).
		Default(DummyConfig).
		Short(c.FLAG_CLIENT_CONFIG_SHORT).
		StringVar(&client.ConfigPath)

	kingpin.
		Flag(c.FLAG_TLS_FOOTGUN, c.FLAG_TLS_FOOTGUN_DESC).
		Hidden().
		Action(client.disableCertVerification).
		Bool()

	kingpin.
		Flag(c.FLAG_VERSION, c.FLAG_VERSION_DESC).
		Action(printVersion).
		Bool()

	return &client
}

func (client *Client) addHeaders(r *http.Request, writer bool) error {
	info, err := loadClientInfo(client.ConfigPath, writer)
	if err != nil {
		return err
	}

	r.SetBasicAuth(info.username, info.secret)
	r.Header.Add(c.HEADER_QUEUE, strings.Join(info.queues, c.HEADER_QUEUE_SEP))

	return nil
}

func responseError(r *http.Response) error {
		if r.ContentLength >= 1 {
			buf := make([]byte, r.ContentLength)
			r.Body.Read(buf)
			if len(buf) > 1 {
				buf = buf[:len(buf) - 1] // Trim trailing newline
			}
			return fmt.Errorf("Status %d: %s", r.StatusCode, buf)
		}

		return fmt.Errorf("Status %d: <No extra data>", r.StatusCode)
}

func (client *Client) WriteMessage(msg c.Message) (*http.Response, error) {

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.EndpointURL(client.Host, client.Port), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	err = client.addHeaders(req, true)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, responseError(resp)
	}

	return resp, nil
}

func (client *Client) ReadMessage() (c.Message, *http.Response, error) {
	var msg c.Message

	req, err := http.NewRequest(http.MethodGet, c.EndpointURL(client.Host, client.Port), nil)
	if err != nil {
		return msg, nil, err
	}

	err = client.addHeaders(req, false)
	if err != nil {
		return msg, nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return msg, resp, err
	} else if resp.StatusCode != 200 {
		return msg, nil, responseError(resp)
	}

	if resp.ContentLength < 1 {
		return msg, nil, fmt.Errorf("%d: <No data>\n", resp.StatusCode)
	}

	buf := make([]byte, resp.ContentLength)
	resp.Body.Read(buf)

	err = json.Unmarshal(buf, &msg)
	return msg, resp, err
}
