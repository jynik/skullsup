// SPDX License Identifier: MIT
// Package skullsup provides the Skulls Up! device API and its
// associated data types.
package skullsup

import "strings"

// Device abstraction interface.
type device interface {
	read(n uint) ([]byte, error)
	write(b []byte) error
	close() error
}

// Platform IDs
const (
	SKULL = iota // 10 NeoPixels, alternating left-right, top-to-bottom
	BULB         // 16 NeoPixels, top-to-bottom (left side), then bottom-to-top (right sidde)
	SIM   = 0xf
)

// Opaque
type Skull struct {
	dev  device   // Device handle
	plat platform // Platform attributes
}

const (
	cmd_wake        = 0xff // Wake command after power-up
	cmd_reset       = 0xfe // Clear frames and set to specified color
	cmd_reanimate   = 0xfd // Start (re)animation
	cmd_set_color   = 0xfc // Reset and set fixed color
	cmd_fw_ver      = 0xfb // Retrive firmware version
	cmd_strip_count = 0xfa // Retrieve LED strip count
	cmd_strip_len   = 0xf9 // Retrieve number of LEDs per strip
	cmd_layout      = 0xf8 // Retrieve physical LED layout
	cmd_max_frames  = 0xf7 // Max supported frame count
	// f6 - 0x80 are reserved for future commands

	// Do not include a frame delay, just update LEDs. OR this with LED "ID"
	no_frame_delay = 0x40

	// Addresses all LEDs when loading a frame. 0x3e - 0x00 address single LEDs
	ALL_LEDS = 0x3f

	NOT_READY = "The Dark Revenant is busy sowing seeds of chaos"
)

func New(name string) (*Skull, error) {
	var err error

	s := new(Skull)

	if name == "hexdump" {
		s.dev, err = openHexDumper(name)
	} else {
		s.dev, err = openUartDevice(name)
	}

	if err != nil {
		return nil, err
	}

	// Wake the device up, if it's not alive already...
	if err = s.Wake(); err != nil {
		return nil, err
	}

	if err = s.loadPlatformInfo(); err != nil {
		return nil, err
	}

	return s, err
}

// Skullsup! Library Version
func Version() string {
	return version
}

func (s *Skull) Wake() error {
	return s.dev.write([]byte{cmd_wake, '1', '3', '8'})
}

func (s *Skull) Reset() error {
	return s.dev.write([]byte{cmd_reset, 0x00, 0x00, 0x00})
}

func (s *Skull) setColor(c Color) error {
	return s.dev.write([]byte{cmd_set_color, c.red, c.green, c.blue})
}

func (s *Skull) SetColor(colorStr string) error {
	colorStr = strings.TrimSpace(colorStr)
	if color, err := NewColor(colorStr); err != nil {
		return err
	} else {
		return s.setColor(color)
	}
}

func (s *Skull) loadFrame(f Frame) error {
	var cmd uint8 = f.led
	if !f.delay {
		cmd |= no_frame_delay
	}
	return s.dev.write([]byte{cmd, f.color.red, f.color.green, f.color.blue})
}

func (s *Skull) loadFrames(frames []Frame) error {
	for _, f := range frames {
		if err := s.loadFrame(f); err != nil {
			return err
		}
	}
	return nil
}

func (s *Skull) reanimate(period uint16) error {
	var per_msb uint8 = uint8(period >> 8)
	var per_lsb uint8 = uint8(period & 0xff)
	return s.dev.write([]byte{cmd_reanimate, per_msb, per_lsb, 0x00})
}

// User input may contain empty (whitespace) strings - drop them
func sanitizeStrSlice(strs []string) []string {
	ret := []string{}
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s != "" {
			ret = append(ret, s)
		}
	}

	return ret
}

func (s *Skull) Reanimate(frames []string, period uint16) error {
	frames = sanitizeStrSlice(frames)

	// Clear out any previous animations
	if err := s.Reset(); err != nil {
		return err
	}

	for _, frameStr := range frames {
		if f, err := NewFrame(frameStr); err != nil {
			return err
		} else {
			if err = s.loadFrame(f); err != nil {
				s.Reset()
				return err
			}
		}
	}

	return s.reanimate(period)
}

func (s *Skull) Incant(psalm string, args []string, period uint16) error {
	args = sanitizeStrSlice(args)

	frames, defaultPeriod, err := s.getPsalm(psalm, args)
	if err != nil {
		return err
	}

	if period == 0 {
		period = defaultPeriod
	}

	// Clear out any previous animations
	if err = s.Reset(); err != nil {
		return err
	}

	if err = s.loadFrames(frames); err != nil {
		s.Reset()
		return err
	}

	return s.reanimate(period)
}

func (s *Skull) Close() error {
	return s.dev.close()
}
