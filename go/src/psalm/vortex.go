// SPDX License Identifier: MIT
package psalm

import "github.com/jynik/skullsup/go/src/frame"

func vortex(args []string, ledCount uint) ([]frame.Frame, uint16, error) {
	var frames []frame.Frame

	fg, bg, err := getFgBg(args, "ff0000", "000505")
	if err != nil {
		return []frame.Frame{}, 0, err
	}

	for i := uint(0); i < ledCount; i++ {
		frames = append(frames, frame.Frame{frame.ALL_LEDS, bg, false})
		frames = append(frames, frame.Frame{uint8(i), fg, true})
	}

	return frames, 40, nil
}
