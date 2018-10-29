// SPDX License Identifier: MIT
package psalm

import "github.com/jynik/skullsup/go/src/color"

type Range struct {
	Min, Max int
}

func getFgBg(args []string, defaultFg, defaultBg string) (color.Color, color.Color, error) {
	var fg, bg color.Color
	var err error

	if len(args) >= 1 {
		if fg, err = color.New(args[0]); err != nil {
			return fg, bg, err
		}
	} else {
		fg = color.MustCreate(defaultFg)
	}

	if len(args) >= 2 {
		if bg, err = color.New(args[1]); err != nil {
			return fg, bg, err
		}
	} else {
		bg = color.MustCreate(defaultBg)
	}

	return fg, bg, nil
}

func mustGetFgBg(args []string, defaultFg, defaultBg string) (color.Color, color.Color) {
	if fg, bg, err := getFgBg(args, defaultFg, defaultBg); err != nil {
		panic(err)
	} else {
		return fg, bg
	}
}
