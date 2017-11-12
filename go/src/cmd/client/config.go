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
	"strings"

	c "../common"
)

type clientInfo struct {
	username string
	secret   string
	queues   []string
}

type clientDaemon struct {
	configPath string // Path to client configuration
	pollRate   uint   // Poll rate in seconds
}

func getFields(s *bufio.Scanner) ([]string) {
	line := strings.Replace(s.Text(), "\r", "", -1)
	line = strings.Replace(line, "\t", "", -1)
	return strings.Split(line, " ")
}

func (c *Client) LoadDefaultServer(filename string) error {
	file, err := os.Open(filename)
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
			return fmt.Errorf("Line %d in %s contains to few fields.", lineNo, filename)
		}

		client.Host = fields[1]
		client.Port =
		return nil
	}
}

func loadClientInfo(filename string, write bool) (*clientInfo, error) {
	var client clientInfo

	file, err := os.Open(filename)
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
			return nil, fmt.Errorf("Line %d in %s contains too few fields.", lineNo, filename)
		}

		switch strings.ToLower(fields[3]) {
		case "read":
			if write {
				continue
			}
		case "write":
			if !write {
				continue
			}
		default:
			return nil, fmt.Errorf("Line %d in %s specifies an invalid permission (%s).", lineNo, filename)
		}

		client.username = fields[1]
		client.secret = fields[2]

		for _, queue := range fields[4:] {
			if !c.IsValidQueueName(queue) {
				return nil, fmt.Errorf("Line %d in %s contains an invalid queue.", lineNo, filename)
			}

			client.queues = append(client.queues, queue)
		}

		return &client, nil
	}

	return nil, errors.New("No client information found.")
}
