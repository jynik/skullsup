// SPDX License Identifier: MIT
// Package skullsup provides the Skulls Up! device API and its
// associated data types.
package skullsup

import (
	"errors"
	"strings"
	"time"
)

// Device abstraction interface.
type device interface {
	read(n uint) ([]byte, error)
	write(b []byte, check_ack bool) (byte, error)
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
	cmd_summon      = 0xff // Summon the device to accept commands
						   // This is required when the device is sleeping
						   // or if it's currently reanimated.
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
)

// Error messages
const (

	// Used to indicate device is not ready to accept commands
	ErrorNotReady = "The Dark Revenant is busy sowing seeds of chaos. Be patient."

	// Command timed out
	ErrorTimeout = "Our cries have gone unanswered and we've given up."
)


// Return a checksum for a payload sent to a device
func checksum(payload []byte) byte {
	ret := byte(0)
	for _, b := range payload {
		ret += b
	}
	return ret
}

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

	if err = s.summon(); err != nil {
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

// Ensure the device is ready to accept commands by sending the summon command
// and ensuring that we've gotten a valid ACK.
func (s *Skull) summon() error {
	var err error
	const maxRetries = 10

	for i := 0; i < 10; i++ {
		_, err = s.dev.write([]byte{cmd_summon, '1', '3', '8'}, true)
		if err == nil {
			return nil
		}

		time.Sleep(250 * time.Microsecond)
	}

	if strings.Contains(err.Error(), "EOF") {
		err = errors.New(ErrorTimeout)
	}

	return err
}

func (s *Skull) setColor(c Color) error {
	if err := s.summon(); err != nil {
		return err
	}
	_, err  := s.dev.write([]byte{cmd_set_color, c.red, c.green, c.blue}, true)
	return err
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
	_, err := s.dev.write([]byte{cmd, f.color.red, f.color.green, f.color.blue}, true)
	return err
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
	_, err :=  s.dev.write([]byte{cmd_reanimate, per_msb, per_lsb, 0x00}, true)
	return err
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

func (s *Skull) Reanimate(frameStrs []string, period uint16) error {
	frameStrs = sanitizeStrSlice(frameStrs)

	if err := s.summon(); err != nil {
		return err
	}

	frames := []Frame{}
	for _, frameStr := range frameStrs {
		if f, err := NewFrame(frameStr); err != nil {
			return err
		} else {
			frames = append(frames, f)
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

	if err := s.summon(); err != nil {
		return err
	}

	if err = s.loadFrames(frames); err != nil {
		return err
	}

	return s.reanimate(period)
}

func (s *Skull) Close() error {
	return s.dev.close()
}
