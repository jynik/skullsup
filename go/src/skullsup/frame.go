// SPDX License Identifier: MIT
package skullsup

import (
	"fmt"
	"strconv"
	"strings"
)

type Frame struct {
	led   uint8 // LED # or skullsup.ALL_LEDS
	color Color // RGB color value
	delay bool  // Include an intra-frame delay after updating LED(s)
}

func validDuration(x uint64) bool {
	return x < 4096
}

// color[:led[:options]]
func NewFrame(s string) (Frame, error) {
	var frame Frame
	var err error
	var tmp uint64

	fields := strings.Split(s, ":")
	if len(fields) > 3 {
		return Frame{}, fmt.Errorf("Invalid frame: %s", s)
	}

	frame.color, err = NewColor(fields[0])
	if err != nil {
		return Frame{}, err
	}

	if len(fields) >= 2 {
		tmp, err = strconv.ParseUint(fields[1], 10, 8)
		if err != nil {
			return Frame{}, fmt.Errorf("Invalid LED number: %s", fields[2])
		} else if tmp >= ALL_LEDS {
			return Frame{}, fmt.Errorf("Invalid LED number: %d", tmp)
		}
		frame.led = uint8(tmp)
	} else {
		frame.led = ALL_LEDS
	}

	frame.delay = true
	if len(fields) >= 3 {
		if fields[2] == "n" || fields[3] == "N" {
			frame.delay = false
		} else {
			return Frame{}, fmt.Errorf("Invalid option flag: %s", fields[3])
		}
	}

	return frame, nil
}

func MustCreateFrame(s string) Frame {
	if frame, err := NewFrame(s); err != nil {
		panic(err)
	} else {
		return frame
	}
}

func NewFrameColor(c Color) Frame {
	return Frame{led: ALL_LEDS, color: c, delay: true}
}

func NewFrameLed(l uint8, c Color, d bool) Frame {
	return Frame{led: l, color: c, delay: d}
}

func (f Frame) Color() Color {
	return f.color
}
