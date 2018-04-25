package sanny

import (
	"encoding/binary"
	"math"
)

const (
	SIZE_FLOAT32 = 4
	SIZE_INT32   = 4
)

func IntsToBytes(ints []int) []byte {
	bytes := make([]byte, len(ints)*SIZE_INT32)
	for i, _ := range ints {
		binary.BigEndian.PutUint32(bytes[SIZE_INT32*i:SIZE_INT32*i+SIZE_INT32], uint32(ints[i]))
	}
	return bytes
}

func BytesToInts(bytes []byte) []int {
	size := len(bytes) / SIZE_INT32
	ints := make([]int, size)
	for i := 0; i < size; i++ {
		ints[i] = int(int32(binary.BigEndian.Uint32(bytes[0:SIZE_INT32])))
		bytes = bytes[SIZE_INT32:]
	}
	return ints
}

func Float32sToBytes(floats []float32) []byte {
	bytes := make([]byte, len(floats)*SIZE_FLOAT32)
	for i, _ := range floats {
		binary.BigEndian.PutUint32(bytes[SIZE_FLOAT32*i:SIZE_FLOAT32*i+SIZE_FLOAT32], math.Float32bits(floats[i]))
	}
	return bytes
}

func BytesToFloat32s(bytes []byte) []float32 {
	size := len(bytes) / SIZE_FLOAT32
	floats := make([]float32, size)
	for i := 0; i < size; i++ {
		floats[i] = math.Float32frombits(binary.BigEndian.Uint32(bytes[0:SIZE_FLOAT32]))
		bytes = bytes[SIZE_FLOAT32:]
	}
	return floats
}
