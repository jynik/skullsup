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
	"github.com/jynik/skullsup/go/src/network"
	"github.com/jynik/skullsup/go/src/network/client"
	"github.com/jynik/skullsup/go/src/psalm"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

const usageText = "Usage: %s [options]\n" +
	"Incant a random psalm with random arguments.\n" +
	"\n" +
	"Options:\n"

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), usageText, path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	os.Exit(0)
}

const psalmArgHelp = "Incant the specified psalm with random arguments. " +
	"If not specified, a random psalm will be used."

func main() {
	var flags cmdline.WriterFlags
	var psalmArg string

	flags.Init()
	flag.Usage = usage

	flag.StringVar(&psalmArg, "psalm", "", psalmArgHelp)
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

	rand.Seed(time.Now().UTC().UnixNano())

	client, err := client.New(flags.Cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if flags.Period > 0 {
		client.Cfg.FramePeriod = flags.Period
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

	msg := network.Message{
		Command: network.CmdIncant,
		Args:    psalm.Random(psalmArg),
		Period:  client.Cfg.FramePeriod,
	}

	client.Log.Debug("Writing to %s\n", queue)

	err = client.Write(&msg, queue)
	if err != nil {
		client.Log.Error("%s\n", err)
		os.Exit(-1)
	}
}
