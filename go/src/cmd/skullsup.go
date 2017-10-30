// SPDX License Identifier: MIT
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"

	"../skullsup"
	c "./common"
	"./common/defaults"
)

var nonDefaultPeriod = false

func gotPeriodArg(*kingpin.ParseContext) error {
	nonDefaultPeriod = true
	return nil
}

var (
	device = kingpin.
		Flag(c.FLAG_DEVICE, c.FLAG_DEVICE_DESC).
		Default(defaults.DEVICE_NAME).
		Short(c.FLAG_DEVICE_SHORT).
		String()

	period_ms = kingpin.
			Flag(c.FLAG_PERIOD, c.FLAG_PERIOD_DEV_DESC).
			Action(gotPeriodArg).
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

	list = kingpin.
		Command(c.CMD_LIST, c.CMD_LIST_DESC).
		Alias(c.CMD_LIST_ALIAS)

	reanimate = kingpin.
			Command(c.CMD_REANIM, c.CMD_REANIM_DESC).
			Alias(c.CMD_REANIM_ALIAS)

	frameStrs = reanimate.
			Arg(c.ARG_FRAMESTR, c.ARG_FRAMESTR_DESC).
			Required().
			Strings()

)

func main() {
	var skull *skullsup.Skull
	var err error

	kingpin.Command(c.CMD_VERSION, c.CMD_VERSION_DESC)

	cmd := kingpin.Parse()
	// Test for a valid command before attempting to open the device
	switch cmd {
	case c.CMD_VERSION:
		fmt.Println(skullsup.Version())
		os.Exit(0)
	case c.CMD_COLOR:
	case c.CMD_LIST:
	case c.CMD_REANIM:
	case c.CMD_INCANT:
	default:
		fmt.Fprintln(os.Stderr, c.ERR_INVALID_CMD)
		return
	}

	skull, err = c.OpenDevice(*device)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer skull.Close()

	switch cmd {

	case c.CMD_COLOR:
		err = skull.SetColor(*colorValue)

	case c.CMD_LIST:
		c.PrintPsalms()

	case c.CMD_REANIM:
		err = skull.Reanimate(*frameStrs, *period_ms)

	case c.CMD_INCANT:
		period := uint16(0)
		if nonDefaultPeriod {
			period = *period_ms
		}
		err = skull.Incant(*psalm, *psalmArgs, period)

	default:
		panic(errors.New(c.ERR_BUG))
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
