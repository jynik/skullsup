// SPDX License Identifier: MIT
package skullsup

import (
	"errors"
	"strings"
)

func getFgBg(args []string, defaultFg, defaultBg string) (Color, Color, error) {
	var fg, bg Color
	var err error

	if len(args) >= 1 {
		if fg, err = NewColor(args[0]); err != nil {
			return fg, bg, err
		}
	} else {
		fg = MustGetNewColor(defaultFg)
	}

	if len(args) >= 2 {
		if bg, err = NewColor(args[1]); err != nil {
			return fg, bg, err
		}
	} else {
		bg = MustGetNewColor(defaultBg)
	}

	return fg, bg, nil
}

func mustGetFgBg(args []string, defaultFg, defaultBg string) (Color, Color) {
	if fg, bg, err := getFgBg(args, defaultFg, defaultBg); err != nil {
		panic(err)
	} else {
		return fg, bg
	}
}

func aura(args []string) ([]Frame, uint16, error) {
	var frames []Frame

	var colors = [...]string{
		"400000", "802010", "c04020", "ff6030",
		"d08040", "c0a050", "a0c060", "80ff70",
		"40c080", "20a080", "0880a0", "106080",
		"184060", "202040", "280020",
	}

	// No user-defined fg
	bg, _, err := getFgBg(args, "401010", "000000")
	if err != nil {
		return []Frame{}, 0, err
	}

	for _, c := range colors {
		fg := MustGetNewColor(c)
		frames = append(frames, NewFrameLed(ALL_LEDS, bg, false))
		frames = append(frames, NewFrameLed(4, fg, false))
		frames = append(frames, NewFrameLed(5, fg, true))
	}

	return frames, 85, nil
}

func pulse(args []string) ([]Frame, uint16, error) {
	var color Color
	var frames []Frame
	var err error

	var brightness = [...]uint8{
		4, 8, 16, 32, 64, 128, 255,
		128, 64, 32,
		64, 128, 255,
		232, 200, 172, 128, 100, 70, 50, 30, 20, 10}

	if len(args) == 0 {
		color = MustGetNewColor("ff0000")
	} else {
		color, err = NewColor(args[0])
		if err != nil {
			return []Frame{}, 0, err
		}
	}

	for _, val := range brightness {
		f := NewFrameColor(color.Scale(val, val, val))
		frames = append(frames, f)
	}

	return frames, 85, nil
}

func hellivator(args []string) ([]Frame, uint16, error) {
	var frames []Frame
	var i uint8

	fg, bg, err := getFgBg(args, "008000", "000030")
	if err != nil {
		return []Frame{}, 0, err
	}

	for i = 0; i < 10; i += 2 {
		frames = append(frames, NewFrameLed(ALL_LEDS, bg, false))
		frames = append(frames, NewFrameLed(i, fg, false))
		frames = append(frames, NewFrameLed(i+1, fg, true))
	}

	for i = 7; i > 2; i -= 2 {
		frames = append(frames, NewFrameLed(ALL_LEDS, bg, false))
		frames = append(frames, NewFrameLed(i, fg, false))
		frames = append(frames, NewFrameLed(i-1, fg, true))
	}

	return frames, 85, nil
}

func vortex(args []string) ([]Frame, uint16, error) {
	var frames []Frame
	var i uint8

	fg, bg, err := getFgBg(args, "ff0000", "000505")
	if err != nil {
		return []Frame{}, 0, err
	}

	for i = 0; i < 10; i += 2 {
		frames = append(frames, NewFrameLed(ALL_LEDS, bg, false))
		frames = append(frames, NewFrameLed(i, fg, true))
	}

	for i = 9; i > 1; i -= 2 {
		frames = append(frames, NewFrameLed(ALL_LEDS, bg, false))
		frames = append(frames, NewFrameLed(i, fg, true))
	}

	return frames, 40, nil
}

type Range struct {
	Min, Max int
}

type Psalm struct {
	Name   string
	Args   Range
	Period Range
	Luma   []Range // Luma range per argument
	impl   func([]string) ([]Frame, uint16, error)
}

var psalms = []Psalm{
	{
		Name:   "aura",
		Args:   Range{0, 0},
		Period: Range{45, 100},
		Luma:   []Range{},
		impl:   aura,
	},

	{
		Name:   "hellivator",
		Args:   Range{0, 2},
		Period: Range{65, 150},
		Luma:   []Range{{128, 255}, {0, 32}},
		impl:   hellivator,
	},

	{
		Name:   "pulse",
		Args:   Range{0, 1},
		Period: Range{50, 85},
		Luma:   []Range{{32, 255}},
		impl:   pulse,
	},

	{
		Name:   "vortex",
		Args:   Range{0, 2},
		Period: Range{65, 125},
		Luma:   []Range{{64, 255}, {0, 64}},
		impl:   vortex,
	},
}

func loadPsalm(name string, args []string) ([]Frame, uint16, error) {
	for _, psalm := range psalms {
		if strings.ToLower(name) == psalm.Name {
			return psalm.impl(args)
		}
	}

	return []Frame{}, 0, errors.New("No such preset: " + name)
}

func Psalms() []Psalm {
	return psalms
}
