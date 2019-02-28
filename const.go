package z19

//BaudRate is the UART baud rate of the sensor
const BaudRate = 9600

// Command is a command sent to the sensor
type Command byte

const (
	// CmdGetReading takes a reading
	CmdGetReading Command = 0x86
	// CmdZeroPointCalibration performs a zero point calibration
	CmdZeroPointCalibration Command = 0x87
	// CmdSpanPointCalibration performs a span point calibration
	CmdSpanPointCalibration Command = 0x88
)
