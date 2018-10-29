// SPDX License Identifier: MIT
package psalm

import (
	"fmt"
	"github.com/jynik/skullsup/go/src/frame"
)

func hellivator(args []string, ledCount uint) ([]frame.Frame, uint16, error) {
	var frames []frame.Frame

	if ledCount < 4 {
		return []frame.Frame{}, 0, fmt.Errorf("Animation requires more than 4 LEDs. Only %d specified.", ledCount)
	}

	fg, bg, err := getFgBg(args, "008000", "000030")
	if err != nil {
		return []frame.Frame{}, 0, err
	}

	for i := uint(0); i < ledCount/2; i++ {
		frames = append(frames, frame.Frame{frame.ALL_LEDS, bg, false})
		frames = append(frames, frame.Frame{uint8(i), fg, false})
		frames = append(frames, frame.Frame{uint8(ledCount-1-i), fg, true})
	}

	for j := ledCount/2 - 2; j > 0; j-- {
		frames = append(frames, frame.Frame{frame.ALL_LEDS, bg, false})
		frames = append(frames, frame.Frame{uint8(j), fg, false})
		frames = append(frames, frame.Frame{uint8(ledCount-1-j), fg, true})
	}

	return frames, 85, nil
}
