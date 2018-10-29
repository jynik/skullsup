// SPDX License Identifier: MIT
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/jynik/skullsup/go/src/network/server"
	"github.com/jynik/skullsup/go/src/version"
)

const Version = "1.0.0"

const usageText = "Usage: %s [options]\n" +
	"Run a SkullUps! queue server\n\n" +
	"Options:\n"

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), usageText, path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	os.Exit(0)
}

func main() {
	cfgFile := flag.String("cfg", server.FindDefaultConfig, "Server configuration file")

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

	server, err := server.New(*cfgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = server.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
