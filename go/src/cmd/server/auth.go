/*
 * SPDX License Identifier: MIT
 *
 * Server user configuration file parsing
 *
 * This file grants clients, identified by a unique name (and authenticated by
 * a password), the ability to write data to one command queue or read from
 * multiple command queues. Queues are identified by a version 4 UUID.
 *
 * Permissions are specified on a per-line basis, as follows:
 *
 * # Grant <username> the ability to write to the queue referenced by <uuid>
 * user <username> <hash> write <uuid> [uuid ... [uuid]]
 *
 * # Grant <username> the ability to read from one or more queues
 * user <username> <hash> read  <uuid> [uuid ... [uuid]]
 *
 */
package server

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	c "../common"
)

const bcryptCost = 11

type user struct {
	write  bool     // True if user can write to queue, false otherwise
	name   string   // Username
	hash   string   // Hash of per-user secret
	queues []string // Queues user can accesss
}

func (u *user) canAccessQueue(queue string, write bool) bool {
	if u.write != write {
		return false
	}

	for _, q := range u.queues {
		if strings.Compare(strings.ToLower(queue), strings.ToLower(q)) == 0 {
			return true
		}
	}

	return false
}

func (s *Server) lookupUser(name string) (*user, error) {
	var file *os.File
	var err error
	var user user

	filename := s.config.ClientConfigPath
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
			// Skip empty, commented, or irrelevant lines
			continue
		} else if len(fields) < 5 {
			return nil, fmt.Errorf("Line %s in %s contains too few fields.", strconv.Itoa(lineNo), filename)
		} else if fields[1] != name {
			// Not the user we're interested in
			continue
		}

		user.name = fields[1]

		switch strings.ToLower(fields[3]) {
		case "read":
			user.write = false
		case "write":
			user.write = true
		default:
			return nil, fmt.Errorf("Line %s in %s specifies an invalid permission (%s).", strconv.Itoa(lineNo))
		}

		user.hash = fields[2]

		for _, queue := range fields[4:] {
			if !c.IsValidQueueName(queue) {
				return nil, fmt.Errorf("Line %s in %s contains an invalid queue ID (%s).", strconv.Itoa(lineNo), filename, line)
			}

			user.queues = append(user.queues, queue)
		}

		return &user, nil
	}

	return nil, errors.New(c.ERR_INVAL)
}

func (s *Server) Authenticate(queues []string, username string, password string, write bool) error {
	user, err := s.lookupUser(username)
	if err != nil {
		s.log.VPrintf("Failed to look up user=%s: %s ", username, err)
		// Kill a little time. TODO: Compiled out?
		if _, err = bcrypt.GenerateFromPassword([]byte("-={dummy}=-"), bcryptCost); err != nil {
			s.log.Println(err)
		}
		return errors.New(c.ERR_AUTH)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.hash), []byte(password))
	if err != nil {
		s.log.VPrintf("Invalid password provided: %s\n", err)
		return errors.New(c.ERR_AUTH)
	}

	for _, queue := range queues {
		if !c.IsValidQueueName(queue) {
			s.log.VPrintln("Auth failure due to invalid queue ID")
			return errors.New(c.ERR_AUTH)
		}

		if !user.canAccessQueue(queue, write) {
			s.log.VPrintln("User not authorized to access queue")
			return errors.New(c.ERR_AUTH)
		}
	}

	return nil
}
