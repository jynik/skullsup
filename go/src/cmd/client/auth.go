/*
 * SPDX License Identifier: MIT
 *
 * Load client authentication information from a file formatted as follows:
 *
 * # skullsup-client: Write requested command to all of the specified queues
 * user <username> <secret> write <uuid> [uuid ... [uuid]]
 *
 * # skullsup-daemon: Read from the specified queues in round-robin fashion
 * user <username> <hash> read <uuid> [uuid ... [uuid]]
 */

package client

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func loadClientInfo(filename string, write bool) (*clientInfo, error) {
	var client clientInfo
	var file *os.File
	var err error

	if file, err = os.Open(filename); err != nil {
		return nil, err
	}
	defer file.Close()

	lineNo := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo += 1

		line := strings.Replace(scanner.Text(), "\r", "", -1)
		line = strings.Replace(line, "\t", "", -1)
		fields := strings.Split(line, " ")

		if len(fields) == 0 || fields[0] != "user" {
			// Skip empty, commented lines, or irrelevant lines
			continue
		} else if len(fields) < 4 {
			return nil, fmt.Errorf("Line %s in %s contains too few fields.", strconv.Itoa(lineNo), filename)
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
			return nil, fmt.Errorf("Line %s in %s specifies an invalid permission (%s).", strconv.Itoa(lineNo))
		}

		client.username = fields[1]
		client.secret = fields[2]

		for _, queue := range fields[4:] {
			if !c.IsValidQueueName(queue) {
				return nil, fmt.Errorf("Line %s in %s contains an invalid queue.", strconv.Itoa(lineNo), filename)
			}

			client.queues = append(client.queues, queue)
		}

		return &client, nil
	}

	return nil, errors.New("No client information found.")
}
