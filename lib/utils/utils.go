package utils

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	"github.com/si-co/vpir-code/lib/constants"
)

// MaxBytesLength get maximal []byte length in map[int][]byte
func MaxBytesLength(in map[int][]byte) int {
	max := 0
	for _, v := range in {
		if len(v) > max {
			max = len(v)
		}
	}

	return max
}

// Divides dividend by divisor and rounds up the result to the nearest multiple
func DivideAndRoundUpToMultiple(dividend, divisor, multiple int) int {
	return int(math.Ceil(float64(dividend)/float64(divisor*multiple))) * multiple
}

// Increase num to the next perfect square.
// If the square root is a whole number, do not modify anything.
// Otherwise, return the square of the square root + 1.
func IncreaseToNextSquare(num *int) {
	i, f := math.Modf(math.Sqrt(float64(*num)))
	if f == 0 {
		return
	}
	*num = int(math.Pow(i+1, 2))
}

func RandUint32() (uint32, error) {
	var buf [4]byte
	_, err := rand.Read(buf[:])
	return binary.BigEndian.Uint32(buf[:]) % constants.ModP, err
}
