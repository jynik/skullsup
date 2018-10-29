// SPDX License Identifier: MIT
package color

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

type Color struct {
	Red, Green, Blue uint8
}

func New(s string) (Color, error) {
	s = strings.TrimSpace(s)
	if b, err := hex.DecodeString(s); err != nil || len(b) != 3 {
		return Color{}, fmt.Errorf("Invalid color: %s", s)
	} else {
		return Color{b[0], b[1], b[2]}, nil
	}
}

func Random(lumaMin, lumaMax int) Color {
	if lumaMax < 0 {
		lumaMax = 0
	} else if lumaMax > 256 {
		lumaMax = 256
	}

	if lumaMin < 0 {
		lumaMin = 0
	} else if lumaMin > 256 {
		lumaMin = 256
	} else if lumaMin > lumaMax {
		lumaMin = lumaMax
	}

	y := float32(rand.Intn(lumaMax-lumaMin) + lumaMin)
	u := float32(rand.Intn(256))
	v := float32(rand.Intn(256))

	r := 1.164*(y-16) + 1.596*(v-128)
	if r > 255 {
		r = 255
	} else if r < 0 {
		r = 0
	}

	g := 1.164*(y-16) - 0.813*(v-128) - 0.391*(u-128)
	if g > 255 {
		g = 255
	} else if g < 0 {
		g = 0
	}

	b := 1.164*(y-16) + 2.018*(u-128)
	if b > 255 {
		b = 255
	} else if b < 0 {
		b = 0
	}

	return Color{uint8(r), uint8(g), uint8(b)}
}

func MustCreate(s string) Color {
	ret, err := New(s)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c Color) String() string {
	return fmt.Sprintf("%02x%02x%02x", c.Red, c.Green, c.Blue)
}

func scale(a, b uint8) uint8 {
	result := float32(a)*float32(b)/255.0 + 0.5
	if result > 255.0 {
		result = 255
	}

	return uint8(result)
}

func (c Color) Scale(red, green, blue uint8) Color {
	return Color{
		scale(c.Red, red),
		scale(c.Green, green),
		scale(c.Blue, blue),
	}
}
