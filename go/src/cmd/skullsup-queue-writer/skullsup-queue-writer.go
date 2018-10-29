// SPDX License Identifier: MIT
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jynik/skullsup/go/src/cmdline"
	"github.com/jynik/skullsup/go/src/network"
	"github.com/jynik/skullsup/go/src/network/client"
	"github.com/jynik/skullsup/go/src/psalm"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

const usageText = "Usage: %s [options] <command> [arguments]\n" +
	"\n" +
	"Write a custom command to a SkullsUp! server" +
	"\n" +
	"Commands:\n" +
	"  color <rrggbb>\n" +
	"    Cast colored light upon the Dark Realm, specified as a 3-byte hex string.\n" +
	"  incant <psalm> [args]\n" +
	"    Incant an unholy psalm, with optional changes to its common utterance.\n" +
	"  list\n" +
	"    List available psalms.\n" +
	"  reanimate <frame> [frame] ...\n" +
	"    Reanimate the undead in a manner of your choosing.\n" +
	"\n" +
	"Options:\n"

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), usageText, path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	os.Exit(0)
}

func main() {
	var flags cmdline.WriterFlags

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

	args := flag.Args()

	if len(args) == 1 && strings.ToLower(args[0]) == "list" {
		psalms := psalm.List
		fmt.Println("\nPsalms and optional arguments")
		fmt.Println("------------------------------------------------------")
		for _, p := range psalms {
			line := fmt.Sprintf("  %-16s", p.Name)
			for _, a := range p.ArgNames {
				line += "[" + a + "] "
			}
			fmt.Println(line)
		}
		fmt.Println()
		os.Exit(0)
	} else if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "A command and associated argument are required.")
		os.Exit(1)
	}

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
		rand.Seed(time.Now().UTC().UnixNano())
		queue = client.Cfg.WriteQueues[rand.Intn(numQueues)]
	}
	client.Log.Debug("Writing to %s\n", queue)

	msg := network.Message{
		Command: args[0],
		Args:    args[1:],
		Period:  client.Cfg.FramePeriod,
	}

	err = client.Write(&msg, queue)
	if err != nil {
		client.Log.Error("%s\n", err)
		os.Exit(-1)
	}
}
