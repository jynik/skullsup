// SPDX License Identifier: MIT
package psalm

import (
	"github.com/jynik/skullsup/go/src/color"
	"github.com/jynik/skullsup/go/src/frame"
)

func pulse(args []string, ledCount uint) ([]frame.Frame, uint16, error) {
	var c color.Color
	var frames []frame.Frame
	var err error

	var brightness = [...]uint8{
		4, 8, 16, 32, 64, 128, 255,
		128, 64, 32,
		64, 128, 255,
		232, 200, 172, 128, 100, 70, 50, 30, 20, 10}

	if len(args) == 0 {
		c = color.MustCreate("ff0000")
	} else {
		c, err = color.New(args[0])
		if err != nil {
			return []frame.Frame{}, 0, err
		}
	}

	for _, val := range brightness {
		f := frame.NewColor(c.Scale(val, val, val))
		frames = append(frames, f)
	}

	return frames, 85, nil
}
