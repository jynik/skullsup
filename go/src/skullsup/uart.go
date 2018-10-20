// SPDX License Identifier: MIT
package skullsup

import (
	//"encoding/hex"
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

const baudrate = 9600
const ack_timeout_ms = 500

func openUartDevice(name string) (*uartDevice, error) {
	var err error
	d := new(uartDevice)
	d.name = name

	c := &serial.Config{Name: d.name, Baud: baudrate, ReadTimeout: time.Millisecond * ack_timeout_ms}
	if d.port, err = serial.OpenPort(c); err != nil {
		if strings.Contains(err.Error(), "device or resource busy") {
			return nil, errors.New(ErrorNotReady)
		}
		return nil, err
	}

	return d, nil
}

func (d *uartDevice) read(n uint) ([]byte, error) {
	var buf []byte = make([]byte, n)
	_, err := d.port.Read(buf)
	//fmt.Printf("Read %s\n", hex.Dump(buf))
	return buf, err
}

func (d *uartDevice) write(payload []byte, check_ack bool) (byte, error) {
	if _, err := d.port.Write(payload); err != nil {
		return 0, err
	} else if !check_ack {
		return 0, nil
	}

	ack_exp := checksum(payload)
	//fmt.Printf("Sent <%02x> %s", ack_exp, hex.Dump(payload))

	ack, err := d.read(1)
	if err != nil {
		return 0, err
	}

	if ack[0] != ack_exp {
		return ack[0], fmt.Errorf("Expected ACK = 0x%02x. Got 0x%02x.", ack_exp, ack)
	}

	return ack[0], nil
}

func (d *uartDevice) close() error {
	return d.port.Close()
}
