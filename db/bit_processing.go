package db

// what's a varint?
// a base 128 integer.
// 128 = 128 ^ 1 + 128 ^ 0
// how is 80 represented?
// 80 = 128 ^ 0 + payloadBits = 0 + 80 = 80
func ProcessVarint(b []byte) (int64, int) {
	var x int64
	for i := range b {
		if i < 8 {
			// keep adding the last 7 bits of the current byte unless the MSB is 1
			x = (x << 7) | int64(b[i]&0x7f)
			if b[i]&0x80 == 0 {
				return x, i + 1
			}
		}
		// 9th bit reached, take it as is
		if i == 8 {
			x = (x << 8) | int64(b[i])
			return x, i + 1
		}
	}
	// should not reach here
	return 0, -1
}

func ReadTwos24Bit(b []byte) int64 {
	n := int64(uint64(b[0])<<16 | uint64(b[1])<<8 | uint64(b[2]))
	if n&(1<<23) != 0 {
		n -= (1 << 24)
	}
	return n
}

// 48 bits
// 8 x 6
// we expect 6 bits
// 128 + 64 + 32 + 16 + 8 + 0
func ReadTwos48Bit(b []byte) int64 {
	n := int64(uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 + uint64(b[4])<<8 + uint64(b[5])))
	if n&(1<<47) != 0 {
		n -= (1 << 48)
	}
	return n
}
