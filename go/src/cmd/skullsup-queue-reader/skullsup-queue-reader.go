// SPDX License Identifier: MIT
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jynik/skullsup/go/src/device"
	"github.com/jynik/skullsup/go/src/network"
	"github.com/jynik/skullsup/go/src/network/client"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

func update(c *client.Client, queue string, s *device.Skull) error {
	msg, err := c.Read(queue)
	if err != nil {
		return err
	}

	if len(msg.Args) == 0 {
		return fmt.Errorf("Received \"%s\" with no arguments", msg.Command)
	}

	switch msg.Command {
	case network.CmdColor:
		color := msg.Args[0]
		err := s.SetColor(color)
		if err != nil {
			return err
		}
		c.Log.Debug("Set color: %s\n", color)

	case network.CmdReanimate:
		err := s.Reanimate(msg.Args, uint16(msg.Period))
		if err != nil {
			return err
		}
		c.Log.Debug("Reanimating with period=%d, args=%s\n", msg.Period, msg.Args)

	case network.CmdIncant:
		var args []string
		psalm := msg.Args[0]

		if len(msg.Args) > 1 {
			args = msg.Args[1:]
		}

		err := s.Incant(psalm, args, uint16(msg.Period))
		if err != nil {
			return err
		}

		c.Log.Debug("Incanting %s with period=%d, args=%s\n", psalm, msg.Period, args)

	default:
		c.Log.Error("Burning unknown command: %s\n", msg.Command)
	}

	return nil
}

func main() {
	deviceArg := flag.String("device", "", "Device to connect to")
	queueArg := flag.String("queue", "", "Read from a specific queue. Only used by -once.")
	onceArg := flag.Bool("once", false, "Perform a single read and exit.")
	cfgFileArg := flag.String("cfg", client.FindDefaultConfig, "Configuration file to use")
	versionArg := flag.Bool("version", false, "Display program version and exit")
	apiVersionArg := flag.Bool("api-version", false, "Display SkullsUp! API version and exit")

	flag.Parse()

	if *versionArg {
		fmt.Println(Version)
		return
	} else if *apiVersionArg {
		fmt.Println(version.String)
		return
	}

	client, err := client.New(*cfgFileArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *deviceArg != "" {
		client.Cfg.Device = *deviceArg
	}

	if client.Cfg.Device == "" {
		fmt.Fprintf(os.Stderr, "No device specified in configuration or via command line.\n")
		os.Exit(2)
	}

	device, err := device.New(client.Cfg.Device)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	defer device.Close()

	numQueues := len(client.Cfg.ReadQueues)
	if numQueues < 1 {
		fmt.Fprintf(os.Stderr, "Client does not have any read queues configured.\n")
		os.Exit(4)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	if *onceArg {
		readFrom := *queueArg
		if *queueArg == "" {
			readFrom = client.Cfg.ReadQueues[rand.Intn(numQueues)]
		}

		client.Log.Debug("Reading from queue: %s\n", readFrom)
		err = update(client, readFrom, device)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
		return
	}

	loggedEmpty := 0
	i := 0
	for /*ever!*/ {
		readFrom := *queueArg
		if *queueArg == "" {
			readFrom = client.Cfg.ReadQueues[i]
		}

		client.Log.Debug("Reading from queue: %s\n", readFrom)

		err = update(client, readFrom, device)
		if err != nil {
			if strings.Contains(err.Error(), network.ErrorQueueEmpty) {
				// Avoid filling the logs with this
				if loggedEmpty == 0 {
					client.Log.Error("%s\n", err)
				} else {
					client.Log.Debug("%s\n", err)
				}
				loggedEmpty = (loggedEmpty + 1) % 60
			} else {
				client.Log.Error("%s\n", err)
			}
		}

		i = (i + 1) % numQueues
		time.Sleep(time.Duration(client.Cfg.PollPeriod) * time.Second)
	}
}
