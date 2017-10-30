// SPDX License Identifier: MIT
package skullsup

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type uartDevice struct {
	name string
	port *serial.Port
}

const ack_timeout_s = 2500

func openUartDevice(name string) (*uartDevice, error) {
	var err error
	d := new(uartDevice)
	d.name = name

	c := &serial.Config{Name: d.name, Baud: 115200, ReadTimeout: time.Millisecond * ack_timeout_s}
	if d.port, err = serial.OpenPort(c); err != nil {
		if strings.Contains(err.Error(), "device or resource busy") {
			return nil, errors.New(NOT_READY)
		}
		return nil, err
	}

	return d, nil
}

func (d *uartDevice) write(payload []byte) error {
	var ack_exp byte
	var ack = make([]byte, 1)
	var err error
	var n int

	for _, b := range payload {
		if _, err = d.port.Write([]byte{b}); err != nil {
			return err
		}

		// The DigiCDC implementation can't handle data thrown
		// at it too quickly. Wait for an ACK before continuing
		if n, err = d.port.Read(ack); err != nil {
			return fmt.Errorf("Did not receive ACK for [%02x]", b)
		}

		ack_exp = ^b

		if n != 1 {
			return fmt.Errorf("Expect 1-byte ACK, got %d bytes", n)
		} else if ack[0] != ack_exp {
			return fmt.Errorf("Expected ACK=%02x, got %02x", ack_exp, ack)
		}
	}
	return nil
}

func (d *uartDevice) close() error {
	return d.port.Close()
}
