// SPDX License Identifier: MIT
package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jynik/skullsup/go/src/color"
	"github.com/jynik/skullsup/go/src/device"
	"github.com/jynik/skullsup/go/src/psalm"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

const usageText = "Usage: %s [options] <command> [arguments]\n" +
	"Write a command to a locally-connected SkullsUp! device\n" +
	"\n" +
	"Commands:\n" +
	"  color [rrggbb]\n" +
	"    Cast colored light upon the Dark Realm, specified as a 3-byte hex string.\n" +
	"  incant [psalm] [args]\n" +
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
	deviceArg := flag.String("device", "", "Specifies the Skull to command.")
	periodArg := flag.Uint("period", 0, "Intra-frame period, in ms.")
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
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "No command provided.")
		os.Exit(1)
	}

	if *deviceArg == "" {
		fmt.Fprintln(os.Stderr, "No device specified.")
		os.Exit(2)
	}

	if strings.ToLower(args[0]) == "list" {
		psalms := psalm.List
		fmt.Println("\nPsalms and optional arguments")
		fmt.Println("------------------------------------------------------")
		fmt.Println("  random")
		for _, p := range psalms {
			line := fmt.Sprintf("  %-16s", p.Name)
			for _, a := range p.ArgNames {
				line += "[" + a + "] "
			}
			fmt.Println(line)
		}
		fmt.Println()
		return

	}

	skull, err := device.New(*deviceArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	defer skull.Close()

	rand.Seed(time.Now().UTC().UnixNano())

	switch strings.ToLower(args[0]) {
	case "color":
		if len(args) > 1 {
			err = skull.SetColor(args[1])
		} else {
			err = skull.SetColor(color.Random(16, 256).String())
		}
	case "incant":
		if len(args) < 2 || strings.ToLower(args[1]) == "random" {
			psalmArgs := psalm.Random("")
			err = skull.Incant(psalmArgs[0], psalmArgs[1:], uint16(*periodArg))
		} else {
			err = skull.Incant(args[1], args[2:], uint16(*periodArg))
		}

	case "reanimate":
		err = skull.Reanimate(args[1:], uint16(*periodArg))
	default:
		err = errors.New("Invalid command: " + args[0])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
	}
}
