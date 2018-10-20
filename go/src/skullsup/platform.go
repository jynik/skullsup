package skullsup

import (
	"fmt"
)

type fwVersion struct {
	major uint // Major version - non-backwards compatible changes
	minor uint // Minor version - added features, backwards compatible
	patch uint // Patch version - bug fixes and non-functional changes
}

type platform struct {
	maxFrames uint
	numStrips uint
	stripLen  uint
	layout    uint
	ledCount  uint
	ledXform  []uint8
	fw        fwVersion
}

const (
	LAYOUT_INCREMENTING = (0 << 0)
	LAYOUT_ALTERNATING  = (1 << 0)
	LAYOUT_WRAP_NORMAL  = (0 << 1)
	LAYOUT_WRAP_INVERT  = (1 << 1)
)

func xformIdentity(count uint) []uint8 {
	xform := make([]uint8, count)
	for i := uint(0); i < count; i++ {
		xform[i] = uint8(i)
	}
	return xform
}

func xformDeinterleave(count uint) []uint8 {
	xform := make([]uint8, count)
	for i := uint(0); i < count; i++ {
		if i&0x1 == 0 {
			xform[i] = uint8(i / 2)
		} else {
			xform[i] = uint8(count - 1 - i/2)
		}
	}
	return xform
}

func (s *Skull) loadPlatformInfo() error {
	var err error
	var buf []byte

	cmd := []byte{cmd_fw_ver, 0x00, 0x00, 0x00}
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(2); err != nil {
		return err
	}
	packed := uint16(buf[0]) | (uint16(buf[1]) << 8)
	s.plat.fw.major = uint(packed >> 11)
	s.plat.fw.minor = uint((packed >> 6) & 0x1f)
	s.plat.fw.patch = uint(packed & 0x2f)

	cmd[0] = cmd_max_frames
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.maxFrames = uint(buf[0])

	cmd[0] = cmd_strip_count
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.numStrips = uint(buf[0])

	cmd[0] = cmd_strip_len
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.stripLen = uint(buf[0])

	s.plat.ledCount = s.plat.numStrips * s.plat.stripLen

	cmd[0] = cmd_layout
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err := s.dev.read(1); err != nil {
		return err
	} else {
		s.plat.layout = uint(buf[0])
	}

	switch s.plat.layout {
	case (LAYOUT_ALTERNATING | LAYOUT_WRAP_NORMAL):
		s.plat.ledXform = xformIdentity(s.plat.ledCount)
		fmt.Println("Layout:          Incrementing, Wrap-Invert")
	case (LAYOUT_INCREMENTING | LAYOUT_WRAP_INVERT):
		s.plat.ledXform = xformDeinterleave(s.plat.ledCount)
	default:
		return fmt.Errorf("Unimplemented LED address layout: 0x%08x\n", s.plat.layout)
	}

	return nil
}
