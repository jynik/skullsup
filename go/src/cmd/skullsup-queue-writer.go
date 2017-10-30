// SPDX License Identifier: MIT
package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"

	"./client"
	c "./common"
	"./common/defaults"
)

var (
	period_ms = kingpin.
			Flag(c.FLAG_PERIOD, c.FLAG_PERIOD_CLIENT_DESC).
			Default(strconv.Itoa(defaults.PERIOD_MS)).
			Short(c.FLAG_PERIOD_SHORT).
			Uint16()

	color = kingpin.
		Command(c.CMD_COLOR, c.CMD_COLOR_DESC).
		Alias(c.CMD_COLOR_ALIAS)

	colorValue = color.
			Arg(c.ARG_COLORVAL, c.ARG_COLORVAL_DESC).
			Required().
			String()

	incantation = kingpin.
			Command(c.CMD_INCANT, c.CMD_INCANT_DESC).
			Alias(c.CMD_INCANT_ALIAS)

	psalm = incantation.
		Arg(c.ARG_PSALM, c.ARG_PSALM_DESC).
		Required().
		String()

	psalmArgs = incantation.
			Arg(c.ARG_PSALM_ARGS, c.ARG_PSALM_ARGS_DESC).
			Strings()

	reanimate = kingpin.
			Command(c.CMD_REANIM, c.CMD_REANIM_DESC).
			Alias(c.CMD_REANIM_ALIAS)

	frameStrs = reanimate.
			Arg(c.ARG_FRAMESTR, c.ARG_FRAMESTR_DESC).
			Required().
			Strings()
)

func main() {
	var args []string

	httpClient := client.New()
	cmd := kingpin.Parse()

	switch cmd {

	case c.CMD_COLOR:
		args = []string{*colorValue}

	case c.CMD_INCANT:
		args = append([]string{*psalm}, *psalmArgs...)

	case c.CMD_REANIM:
		args = *frameStrs

	case c.CMD_LIST:
		c.PrintPsalms()

	default:
		return
	}

	msg := c.Message{Command: cmd, Args: args, Period: strconv.Itoa(int(*period_ms))}
	_, err := httpClient.WriteMessage(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
