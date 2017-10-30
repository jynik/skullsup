// SPDX License Identifier: MIT
package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"

	c "./common"
	"./common/defaults"
	"./common/logger"
	s "./server"
	"../skullsup"
)

func printVersion(*kingpin.ParseContext) error {
	fmt.Println(skullsup.Version())
	os.Exit(0)
	return nil
}

func processCommandLine(config *s.Config) {
	kingpin.Flag(c.FLAG_CLIENT_CONFIG, c.FLAG_CLIENT_CONFIG_DESC).
		Default(defaults.SERVER_CONFIG).
		Short(c.FLAG_CLIENT_CONFIG_SHORT).
		StringVar(&config.ClientConfigPath)

	kingpin.Flag(c.FLAG_PORT, c.FLAG_PORT_DESC).
		Default(strconv.Itoa(defaults.PORT)).
		Short(c.FLAG_PORT_SHORT).
		Uint16Var(&config.Port)

	kingpin.Flag(c.FLAG_TLS_CERT, c.FLAG_TLS_CERT_DESC).
		Default(defaults.TLS_CERT).
		Short(c.FLAG_TLS_CERT_SHORT).
		StringVar(&config.TlsCertPath)

	kingpin.Flag(c.FLAG_PRIVATE_KEY, c.FLAG_PRIVATE_KEY_DESC).
		Default(defaults.PRIVATE_KEY).
		Short(c.FLAG_PRIVATE_KEY_SHORT).
		StringVar(&config.PrivateKeyPath)

	kingpin.Flag(c.FLAG_LADDRESS, c.FLAG_LADDRESS_DESC).
		Default(defaults.LADDRESS).
		Short(c.FLAG_LADDRESS_SHORT).
		StringVar(&config.ListenAddr)

	kingpin.Flag(c.FLAG_LOGFILE, c.FLAG_LOGFILE_DESC).
		Default(defaults.LOGFILE).
		Short(c.FLAG_LOGFILE_SHORT).
		StringVar(&config.LogFilePath)

	kingpin.Flag(c.FLAG_VERBOSE, c.FLAG_VERBOSE_DESC).
		Default(defaults.VERBOSE).
		Short(c.FLAG_VERBOSE_SHORT).
		BoolVar(&config.Verbose)

	kingpin.Flag(c.FLAG_QUIET, c.FLAG_QUIET_DESC).
		Default("false").
		Short(c.FLAG_QUIET_SHORT).
		BoolVar(&config.Quiet)

	kingpin.Flag(c.FLAG_VERSION, c.FLAG_VERSION_DESC).
		Action(printVersion).
		Bool()

	kingpin.Parse()
}

func main() {
	var server s.Server
	var config s.Config

	processCommandLine(&config)
	if err := server.Run(config); err != nil {
		logger.Fatal("Error: " + err.Error())
	}
}
