package util

func Int16ToBytes(v int16) []byte {
	return []byte{byte(v >> 8), byte(v)}
}

func Int32ToBytes(v int32) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

func Int64ToBytes(v int64) []byte {
	return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

func Uint16ToBytes(v uint16) []byte {
	return []byte{byte(v >> 8), byte(v)}
}

func Uint32ToBytes(v uint32) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

func Uint64ToBytes(v uint64) []byte {
	return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

func Int16FromBytes(b []byte) int16 {
	return int16(b[0])<<8 | int16(b[1])
}

func Int32FromBytes(b []byte) int32 {
	return int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
}

func Int64FromBytes(b []byte) int64 {
	return int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
}

func Uint16FromBytes(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func Uint32FromBytes(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func Uint64FromBytes(b []byte) uint64 {
	return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
}

func Log2(n int) int {
	if n <= 0 {
		return 0
	}
	var r int
	for n >>= 1; n != 0; n >>= 1 {
		r++
	}
	return r
}
