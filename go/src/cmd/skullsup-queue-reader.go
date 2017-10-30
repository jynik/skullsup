// SPDX License Identifier: MIT
package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"./client"
	c "./common"
	"./common/defaults"
	"./common/logger"
)

var (
)

type reader struct {
	httpClient  *client.Client
	skull       string
	log			*logger.Logger

	period      int64	// Sever polling period, in seconds. <= 0 implies no polling, run once.
	logFilePath string	// Path to write log files, stdout, or stderr
	quiet       bool	// Supress all log output
	verbose		bool	// Enable verbose output
}

func (r *reader) update() error {
	r.log.VPrintf("Querying %s:%d\n", r.httpClient.Host, r.httpClient.Port)

	msg, _, err := r.httpClient.ReadMessage()
	if err != nil {
		return err
	}

	r.log.VPrintf("Dequeued message: %s\n", msg)

	numArgs := len(msg.Args)
	if numArgs < 1 {
		return errors.New("Message contained too few arguments.")
	}

	period, err := strconv.ParseInt(msg.Period, 10, 16)
	if err != nil || period < 0 {
		return fmt.Errorf("Invalid period (ms) provided: %s", msg.Period)
	}

	skull, err := c.OpenDevice(r.skull)
	if err != nil {
		return err
	}
	defer skull.Close()

	switch msg.Command {
	case c.CMD_COLOR:
		r.log.VPrintf("Setting color to %s\n", msg.Args[0])
		err = skull.SetColor(msg.Args[0])
	case c.CMD_REANIM:
		r.log.VPrintf("Reanimating with period=%d, args=%s\n", period, msg.Args)
		err = skull.Reanimate(msg.Args, uint16(period))
	case c.CMD_INCANT:
		r.log.VPrintf("Incanting %s with args=%s and period=%d\n", msg.Args[0], msg.Args[1:], period)
		err = skull.Incant(msg.Args[0], msg.Args[1:], uint16(period))
	default:
		r.log.Printf("Invalid command: %s\n", msg.Command)
	}

	return err
}

func handleCmdline(r *reader) {
	var err error

	r.httpClient = client.New()

	kingpin.Flag(c.FLAG_PERIOD, c.FLAG_PERIOD_CLIENT_DESC).
		Short(c.FLAG_PERIOD_SHORT).
		Int64Var(&r.period)

	kingpin.Flag(c.FLAG_LOGFILE, c.FLAG_LOGFILE_DESC).
		Default(defaults.LOGFILE).
		Short(c.FLAG_LOGFILE_SHORT).
		StringVar(&r.logFilePath)

	kingpin.Flag(c.FLAG_VERBOSE, c.FLAG_VERBOSE_DESC).
		Short(c.FLAG_VERBOSE_SHORT).
		BoolVar(&r.verbose)

	kingpin.Flag(c.FLAG_QUIET, c.FLAG_QUIET_DESC).
		Short(c.FLAG_QUIET_SHORT).
		BoolVar(&r.quiet)

	kingpin.Flag(c.FLAG_DEVICE, c.FLAG_DEVICE_DESC).
		Default(defaults.DEVICE_NAME).
		Short(c.FLAG_DEVICE_SHORT).
		StringVar(&r.skull)

	kingpin.Parse()

	if !r.quiet {
		r.log, err = logger.New(r.logFilePath, r.verbose)
	} else {
		r.log, err = logger.New("", false)
	}

	if err != nil {
		logger.Fatal(err)
	}

	// Just fall back to "read once" operation for an invalid value
	if r.period < 0 {
		r.period = 0
	}

	skull, err := c.OpenDevice(r.skull)
	skull.Reset()
	defer skull.Close()
	if err != nil {
		r.log.Fatal(err)
	}
}

func main() {
	var reader reader

	handleCmdline(&reader)

	runOnce := reader.period == 0
	if !runOnce {
		reader.log.Printf("Updating every %d seconds...", reader.period)
	}

	for {
		err := reader.update()
		if err != nil {
			reader.log.Print(err)
		}
		if !runOnce {
			time.Sleep(time.Duration(reader.period) * time.Second)
		} else {
			break
		}
	}
}
