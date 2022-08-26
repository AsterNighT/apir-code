package utils

import (
	"crypto/aes"
	"math"
)

// WARNING: DO NOT USE THESE KEYS IN PRODUCTION!
var SeedMatrixA = [aes.BlockSize]byte{19, 177, 222, 148, 155, 239, 159, 227, 155, 99, 246, 214, 220, 162, 30, 66}

type ParamsLWE struct {
	P     uint32  // plaintext modulus
	N     int     // lattice/secret dimension
	Sigma float64 // Error parameter

	L int    // number of rows of database
	M int    // number of columns of database
	B uint64 // bound used in reconstruction

	SeedA    *PRGKey // matrix  used to generate digest
	BytesMod int     // bytes of the modulo, either 4 for 32 bits or 8 for 64
}

func ParamsDefault() *ParamsLWE {
	return &ParamsLWE{
		P:        2,
		N:        2300,
		Sigma:    6.4,
		L:        512,
		M:        128,
		B:        1000,
		SeedA:    GetDefaultSeedMatrixA(),
		BytesMod: 8,
	}
}

func ParamsWithDatabaseSize(rows, columns int) *ParamsLWE {
	p := ParamsDefault()
	p.L = rows
	p.M = columns
	p.B = uint64(rows * 12 * int(math.Ceil(p.Sigma))) // rows is equal to sqrt(\ell), 12 is ~ sqrt(128)

	return p
}

func GetDefaultSeedMatrixA() *PRGKey {
	key := PRGKey(SeedMatrixA)
	return &key
}

// TODO: remove if we go with the 32-bits version
func ParamsDefault128() *ParamsLWE {
	p := ParamsDefault()
	p.BytesMod = 16

	return p
}

// TODO: remove if we go with the 32-bits version
func ParamsWithDatabaseSize128(rows, columns int) *ParamsLWE {
	p := ParamsDefault128()
	p.L = rows
	p.M = columns

	return p
}
