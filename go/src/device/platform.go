// SPDX License Identifier: MIT
package device

type fwVersion struct {
	major uint // Major version - non-backwards compatible changes
	minor uint // Minor version - added features, backwards compatible
	patch uint // Patch version - bug fixes and non-functional changes
}

type platform struct {
	maxFrames uint
	numStrips uint
	stripLen  uint
	ledCount  uint
	fw        fwVersion
}

func (s *Skull) loadPlatformInfo() error {
	var err error
	var buf []byte

	cmd := []byte{CmdFwVersion, 0x00, 0x00, 0x00}
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

	cmd[0] = CmdMaxFrames
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.maxFrames = uint(buf[0])

	cmd[0] = CmdStripCount
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.numStrips = uint(buf[0])

	cmd[0] = CmdStripLen
	if _, err := s.dev.write(cmd, true); err != nil {
		return err
	}

	if buf, err = s.dev.read(1); err != nil {
		return err
	}
	s.plat.stripLen = uint(buf[0])

	s.plat.ledCount = s.plat.numStrips * s.plat.stripLen

	return nil
}
