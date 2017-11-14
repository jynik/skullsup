/*
 * SPDX License Identifier: MIT
 *
 * Load client configuration information from a file formatted as follows:
 *
 * Use the following remote server and port as the defaults. These can be
 * overridden with their associated command line arguments.
 *	server <host> <port>
 *
 * Write requested command to all of the specified queues:
 *	user <username> <secret> write <uuid> [uuid ... [uuid]]
 *
 * Read from the specified queues in round-robin fashion:
 *	user <username> <secret> read <uuid> [uuid ... [uuid]]
 *
 */

package client

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../common"
	"../common/defaults"
)

type userInfo struct {
	username string
	secret   string
	queues   []string
}

type clientDaemon struct {
	configPath string // Path to client configuration
	pollRate   uint   // Poll rate in seconds
}

// Expand the config path if it contains our dummy location
func (client *Client) expandConfigPath() error {
	if strings.Contains(client.ConfigPath, defaults.DUMMY_PREFIX) {
		re := regexp.MustCompile("[^a-zA-Z0-9_]")
		envvar := string(re.ReplaceAll([]byte(defaults.DUMMY_PREFIX), []byte{}))
		if location, exists := os.LookupEnv(envvar); !exists {
			return errors.New("No such environment variable: " + envvar)
		} else {
			client.ConfigPath = location + client.ConfigPath[len(defaults.DUMMY_PREFIX):]
		}
	}

	return nil
}

func getFields(s *bufio.Scanner) []string {
	line := strings.Replace(s.Text(), "\r", "", -1)
	line = strings.Replace(line, "\t", "", -1)
	return strings.Split(line, " ")
}

func (c *Client) loadRemoteHostDefaults() error {
	c.Host = defaults.REMOTE
	c.Port = defaults.PORT

	file, err := os.Open(c.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	defer file.Close()

	lineNo := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineNo += 1
		fields := getFields(scanner)

		if len(fields) == 0 || strings.ToLower(fields[0]) != "server" {
			// Skip empty, commented lines, or irrelevant lines
			continue
		}

		if len(fields) < 3 {
			return fmt.Errorf("Line %d in %s contains to few fields.", lineNo, c.ConfigPath)
		}

		c.Host = fields[1]

		if port, err := strconv.ParseUint(fields[2], 10, 16); err != nil || port < 0 {
			return fmt.Errorf("Invalid port number (%s) on line %d of %s.", fields[2], lineNo, c.ConfigPath)
		} else {
			c.Port = uint16(port)
		}
	}

	return nil
}

func (c *Client) loadUserInfo(writer bool) (*userInfo, error) {
	var ret userInfo

	file, err := os.Open(c.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineNo := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo += 1
		fields := getFields(scanner)

		if len(fields) == 0 || strings.ToLower(fields[0]) != "user" {
			// Skip empty, commented lines, or irrelevant lines
			continue
		} else if len(fields) < 4 {
			return nil, fmt.Errorf("Line %d in %s contains too few fields.", lineNo, c.ConfigPath)
		}

		switch strings.ToLower(fields[3]) {
		case "read":
			if writer {
				continue
			}
		case "write":
			if !writer {
				continue
			}
		default:
			return nil, fmt.Errorf("Line %d in %s specifies an invalid permission (%s).", lineNo, c.ConfigPath)
		}

		ret.username = fields[1]
		ret.secret = fields[2]

		for _, queue := range fields[4:] {
			if !common.IsValidQueueName(queue) {
				return nil, fmt.Errorf("Line %d in %s contains an invalid queue.", lineNo, c.ConfigPath)
			}

			ret.queues = append(ret.queues, queue)
		}

		return &ret, nil
	}

	return nil, errors.New("No client information found.")
}
