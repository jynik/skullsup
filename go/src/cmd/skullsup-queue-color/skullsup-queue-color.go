// SPDX License Identifier: MIT
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/jynik/skullsup/go/src/cmdline"
	"github.com/jynik/skullsup/go/src/color"
	"github.com/jynik/skullsup/go/src/network"
	"github.com/jynik/skullsup/go/src/network/client"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

const usageText = "Usage: %s [options]\n" +
	"Write a random color to a SkullsUp! server queue\n\n" +
	"Options:\n"

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), usageText, path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	os.Exit(0)
}

func main() {
	var flags cmdline.WriterFlags

	flags.Period = -1 // Don't use this option
	flags.Init()

	versionArg := flag.Bool("version", false, "Display program version and exit")
	apiVersionArg := flag.Bool("api-version", false, "Display SkullsUp! API version and exit")

	flag.Usage = usage
	flag.Parse()

	if *versionArg {
		fmt.Println(Version)
		return
	} else if *apiVersionArg {
		fmt.Println(version.String)
		return
	}

	client, err := client.New(flags.Cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	numQueues := len(client.Cfg.WriteQueues)
	if numQueues < 1 {
		fmt.Fprint(os.Stderr, "Client does not have any write queues configured.\n")
		os.Exit(3)
	}

	queue := flags.Queue
	if queue == "" {
		queue = client.Cfg.WriteQueues[rand.Intn(numQueues)]
	}
	client.Log.Debug("Writing to %s\n", queue)

	rand.Seed(time.Now().UTC().UnixNano())
	msg := network.Message{
		Command: network.CmdColor,
		Args:    []string{color.Random(10, 256).String()},
		Period:  client.Cfg.FramePeriod,
	}

	err = client.Write(&msg, queue)
	if err != nil {
		client.Log.Error("%s\n", err)
		os.Exit(-1)
	}
}
