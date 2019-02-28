package z19

import (
	"github.com/tarm/serial"
)

// CreateSerialConfig creates a serial configuration with the baud rate of mh-z19
func CreateSerialConfig() *serial.Config {
	return &serial.Config{
		Baud: BaudRate,
	}
}
