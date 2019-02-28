package z19

func crc(d []byte) byte {
	var res byte
	for i, p := range d {
		if i == 0 {
			continue
		}
		res += p
	}
	return 0xff - res + 1
}

func withCRC(d []byte) []byte {
	return append(d, crc(d))
}
