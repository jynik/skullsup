// SPDX License Identifier: MIT
package device

import (
	"errors"
	"strings"
	"time"

	"github.com/jynik/skullsup/go/src/color"
	"github.com/jynik/skullsup/go/src/frame"
	"github.com/jynik/skullsup/go/src/psalm"
)

// Device abstraction interface.
type device interface {
	read(n uint) ([]byte, error)
	write(b []byte, check_ack bool) (byte, error)
	close() error
}

// Platform IDs
const (
	SKULL = iota // 2 Adafruit NeoPixel Sticks (2x8 LEDs)
	SIM   = 0xf
)

// Opaque
type Skull struct {
	dev  device   // Device handle
	plat platform // Platform attributes
}

const (
	// Summon the device to accept commands. This is required when the device
	// is sleeping or if it's currently reanimated.
	CmdSummon = 0xff

	// Clear frames and set to specified color
	CmdReset = 0xfe

	// Start (re)animation
	CmdReanimate = 0xfd

	// Reset and set fixed color
	CmdSetColor = 0xfc

	// Get firmware version
	CmdFwVersion = 0xfb

	// Retrieve LED strip count
	CmdStripCount = 0xfa

	// Retrieve number of LEDs per strip
	CmdStripLen = 0xf9

	// Retrieve physical LED layout
	CmdLayout = 0xf8

	// Max supported frame count
	CmdMaxFrames = 0xf7

	// f6 - 0x80 are reserved for future commands

	// Do not include a frame delay, just update LEDs. OR this with LED "ID"
	NoFrameDelay = 0x40

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

// Ensure the device is ready to accept commands by sending the summon command
// and ensuring that we've gotten a valid ACK.
func (s *Skull) summon() error {
	var err error
	const maxRetries = 10

	for i := 0; i < 10; i++ {
		_, err = s.dev.write([]byte{CmdSummon, '1', '3', '8'}, true)
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

func (s *Skull) setColor(c color.Color) error {
	if err := s.summon(); err != nil {
		return err
	}
	_, err := s.dev.write([]byte{CmdSetColor, c.Red, c.Green, c.Blue}, true)
	return err
}

func (s *Skull) SetColor(colorStr string) error {
	if color, err := color.New(colorStr); err != nil {
		return err
	} else {
		return s.setColor(color)
	}
}

func (s *Skull) loadFrame(f frame.Frame) error {
	var cmd uint8 = f.Led
	if !f.Delay {
		cmd |= NoFrameDelay
	}
	_, err := s.dev.write([]byte{cmd, f.Color.Red, f.Color.Green, f.Color.Blue}, true)
	return err
}

func (s *Skull) loadFrames(frames []frame.Frame) error {
	for _, f := range frames {
		if err := s.loadFrame(f); err != nil {
			return err
		}
	}
	return nil
}

func (s *Skull) reanimate(period uint16) error {
	var msb uint8 = uint8(period >> 8)
	var lsb uint8 = uint8(period & 0xff)
	_, err := s.dev.write([]byte{CmdReanimate, msb, lsb, 0x00}, true)
	return err
}

func (s *Skull) Reanimate(frameStrs []string, period uint16) error {
	if err := s.summon(); err != nil {
		return err
	}

	if period == 0 {
		period = 100
	}

	frames := []frame.Frame{}
	for _, frameStr := range frameStrs {
		if f, err := frame.New(frameStr); err != nil {
			return err
		} else {
			frames = append(frames, f)
		}
	}

	if err := s.loadFrames(frames); err != nil {
		return err
	}

	return s.reanimate(period)
}

func (s *Skull) Incant(psalmName string, args []string, period uint16) error {
	frames, defaultPeriod, err := psalm.Lookup(psalmName, args, s.plat.ledCount)
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
