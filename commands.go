package z19

import (
	"fmt"
	"io"
)

func writeCmd(w io.Writer, d []byte) error {
	d = withCRC(d)
	l, err := w.Write(d)
	if err != nil {
		return err
	} else if l != len(d) {
		return fmt.Errorf("expected to write %d bytes, however only %d is written", len(d), l)
	}
	return nil
}

func readRes(r io.Reader, length int) ([]byte, error) {
	length++ // add crc
	b := make([]byte, length)
	if l, err := io.ReadFull(r, b); err != nil {
		return nil, err
	} else if l != length {
		return nil, fmt.Errorf("expected to read %d bytes, however only %d is read", length, l)
	}
	if expectedCrc := crc(b[:len(b)-1]); expectedCrc != b[len(b)-1] {
		return nil, fmt.Errorf("crc8 mismatch, got %x, %x expected", b[len(b)-1], expectedCrc)
	}
	return b[:len(b)-1], nil
}

func buildRequest(id byte, cmd byte) []byte {
	return []byte{0xff, id, cmd, 0x00, 0x00, 0x00, 0x00, 0x00}
}

// TakeReading takes a reading from the sensor
func TakeReading(i io.ReadWriter) (uint16, error) {
	if err := writeCmd(i, buildRequest(0x01, byte(CmdGetReading))); err != nil {
		return 0, err
	}
	data, err := readRes(i, 8)
	if err != nil {
		return 0, err
	}
	if data[1] != byte(CmdGetReading) {
		return 0, fmt.Errorf("returned command mismatch, got 0x%x, 0x%x expected", data[1], CmdGetReading)
	}
	return calcConcentration(data[2], data[3]), nil
}

// ZeroPointCalibration sends a zero calibration command to the sensor
func ZeroPointCalibration(i io.Writer) error {
	return writeCmd(i, buildRequest(0x01, byte(CmdZeroPointCalibration)))
}

// SpanPointCalibration sends a span point calibration command to the sensor with a reference CO2 concentration
func SpanPointCalibration(i io.Writer, ref uint16) error {
	req := buildRequest(0x01, byte(CmdSpanPointCalibration))
	req[3] = byte(ref >> 8)
	req[4] = byte(ref & 0x00ff)
	return writeCmd(i, req)
}
