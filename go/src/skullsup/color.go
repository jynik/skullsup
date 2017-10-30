// SPDX License Identifier: MIT
package skullsup

import (
	"encoding/hex"
	"fmt"
)

type Color struct {
	red, green, blue uint8
}

func NewColor(s string) (Color, error) {
	if b, err := hex.DecodeString(s); err != nil || len(b) != 3 {
		return Color{}, fmt.Errorf("Invalid color: %s", s)
	} else {
		return Color{red: b[0], green: b[1], blue: b[2]}, nil
	}
}

func MustGetNewColor(s string) Color {
	ret, err := NewColor(s)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c Color) String() string {
	return fmt.Sprintf("%02x%02x%02x", c.red, c.green, c.blue)
}

func subtract(a, b uint8) uint8 {
	if b > a {
		return 0
	} else {
		return a - b
	}
}

func scale(a, b uint8) uint8 {
	result := float32(a)*float32(b)/255.0 + 0.5
	if result > 255.0 {
		result = 255
	}

	return uint8(result)
}

func (c Color) AddColorWithOverflow(o Color) Color {
	return Color{
		red:   c.red + o.red,
		green: c.green + o.green,
		blue:  c.blue + o.blue,
	}
}

func (c Color) Subtract(red, green, blue uint8) Color {
	return Color{
		red:   subtract(c.red, red),
		green: subtract(c.green, green),
		blue:  subtract(c.blue, blue),
	}
}

func (c Color) Scale(red, green, blue uint8) Color {
	return Color{
		red:   scale(c.red, red),
		green: scale(c.green, green),
		blue:  scale(c.blue, blue),
	}
}
