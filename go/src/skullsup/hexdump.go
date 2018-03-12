// SPDX License Identifier: MIT
package skullsup

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type hexDumper struct {
	dump io.WriteCloser
}

func openHexDumper(filename string) (*hexDumper, error) {
	d := new(hexDumper)
	d.dump = hex.Dumper(os.Stdout)
	return d, nil
}

func (h *hexDumper) info() (uint, fwVersion) {
	return SIM, fwVersion{0, 1, 0}
}

func (h *hexDumper) read(n uint) ([]byte, error) {
	buf := make([]byte, n)
	for i := uint(0); i < n; i++ {
		buf[i] = 0xff
	}
	return buf, nil
}

func (hex *hexDumper) write(payload []byte) error {
	hex.dump.Write(payload)
	fmt.Println()
	return nil
}

func (h *hexDumper) close() error {
	return nil
}
