package z19

func calcConcentration(high, low byte) uint16 {
	return uint16(high)<<8 + uint16(low)
}
