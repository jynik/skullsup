// SPDX License Identifier: MIT
package frame

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jynik/skullsup/go/src/color"
)

const ALL_LEDS = 0x3f

type Frame struct {
	Led   uint8       // LED # or skullsup.ALL_LEDS
	Color color.Color // RGB color value
	Delay bool        // Include an intra-frame delay after updating LED(s)
}

func validDuration(x uint64) bool {
	return x < 4096
}

// color[:led[:options]]
func New(s string) (Frame, error) {
	var frame Frame
	var err error
	var tmp uint64

	fields := strings.Split(s, ":")
	if len(fields) > 3 {
		return Frame{}, fmt.Errorf("Invalid frame: %s", s)
	}

	frame.Color, err = color.New(fields[0])
	if err != nil {
		return Frame{}, err
	}

	if len(fields) >= 2 && strings.ToLower(fields[1]) != "all" {
		tmp, err = strconv.ParseUint(fields[1], 10, 8)
		if err != nil {
			return Frame{}, fmt.Errorf("Invalid LED number: %s", fields[1])
		} else if tmp >= ALL_LEDS {
			return Frame{}, fmt.Errorf("Invalid LED number: %d", tmp)
		}
		frame.Led = uint8(tmp)
	} else {
		frame.Led = ALL_LEDS
	}

	frame.Delay = true
	if len(fields) >= 3 {
		if fields[2] == "n" || fields[2] == "N" {
			frame.Delay = false
		} else {
			return Frame{}, fmt.Errorf("Invalid option flag: %s", fields[3])
		}
	}

	return frame, nil
}

func MustCreate(s string) Frame {
	if frame, err := New(s); err != nil {
		panic(err)
	} else {
		return frame
	}
}

func NewColor(c color.Color) Frame {
	return Frame{ALL_LEDS, c, true}
}

func (f *Frame) String() string {
	c := f.Color.String()
	options := ""
	if !f.Delay {
		options = ":N"
	}
	return fmt.Sprintf("%s:%d%s\n", c, f.Led, options)
}
