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

func (hex *hexDumper) write(payload []byte) error {
	hex.dump.Write(payload)
	fmt.Println()
	return nil
}

func (h *hexDumper) close() error {
	return nil
}
