package dpf

import (
	"math/bits"
	"testing"

	"github.com/si-co/vpir-code/lib/field"
	"github.com/si-co/vpir-code/lib/utils"
)

func BenchmarkEvalFull(bench *testing.B) {
	// db parameters
	blockSize := 16
	numColumns := 200

	// generate random alpha
	alpha, err := new(field.Element).SetRandom(utils.RandomPRG())
	if err != nil {
		panic(err)
	}

	beta := make([]field.Element, blockSize+1)
	beta[0] = field.One()
	for i := 1; i < len(beta); i++ {
		beta[i].Mul(&beta[i-1], alpha)
	}

	// generate one key
	logN := uint64(bits.Len(uint(numColumns) - 1))
	key, _ := Gen(1, beta, logN)

	q := make([]field.Element, numColumns*(blockSize+1))

	bench.ResetTimer()
	bench.ReportAllocs()
	for i := 0; i < bench.N; i++ {
		EvalFullFlatten(key, logN, blockSize+1, q)
	}
}

/*
func BenchmarkXor16(bench *testing.B) {
	a := new(block)
	b := new(block)
	c := new(block)
	for i := 0; i < bench.N; i++ {
		xor16(&c[0], &b[0], &a[0])
	}
}

func TestEval(test *testing.T) {
	logN := uint64(8)
	alpha := uint64(123)
	beta := make([]field.Element, 2)
	beta[0].SetUint64(7613)
	beta[1].SetUint64(991)

	a, b := Gen(alpha, beta, logN)

	sum := make([]field.Element, 2)
	out0 := make([]field.Element, 2)
	out1 := make([]field.Element, 2)
	zero := field.Zero()
	for i := uint64(0); i < (uint64(1) << logN); i++ {
		Eval(a, i, logN, out0)
		Eval(b, i, logN, out1)

		for j := 0; j < 2; j++ {
			sum[j].Add(&out0[j], &out1[j])
		}

		//log.Printf("%v %v %v %v", i, alpha, beta.String(), sum.String())
		if i != alpha && (!sum[0].Equal(&zero) || !sum[1].Equal(&zero)) {
			test.Fail()
		}

		if i == alpha && (!sum[0].Equal(&beta[0]) || !sum[1].Equal(&beta[1])) {
			test.Fail()
		}
	}
}

func TestEvalFull(test *testing.T) {
	logN := uint64(9)
	alpha := uint64(123)
	beta := make([]field.Element, 2)
	beta[0].SetUint64(7613)
	beta[1].SetUint64(991)

	a, b := Gen(alpha, beta, logN)

	sum := make([]field.Element, 2)
	out0 := make([][]field.Element, 1<<logN)
	out1 := make([][]field.Element, 1<<logN)

	for i := 0; i < len(out0); i++ {
		out0[i] = make([]field.Element, 2)
		out1[i] = make([]field.Element, 2)
	}

	EvalFull(a, logN, out0)
	EvalFull(b, logN, out1)

	zero := field.Zero()
	for i := uint64(0); i < (uint64(1) << logN); i++ {
		for j := 0; j < 2; j++ {
			sum[j].Add(&out0[i][j], &out1[i][j])
		}

		//log.Printf("%v %v %v %v", i, alpha, beta.String(), sum.String())
		if i != alpha && (!sum[0].Equal(&zero) || !sum[1].Equal(&zero)) {
			test.Fail()
		}

		if i == alpha && (!sum[0].Equal(&beta[0]) || !sum[1].Equal(&beta[1])) {
			test.Fail()
		}
	}
}

func TestEvalFullShort(test *testing.T) {
	logN := uint64(2)
	alpha := uint64(2)
	beta := make([]field.Element, 2)
	beta[0].SetUint64(7613)
	beta[1].SetUint64(991)

	a, b := Gen(alpha, beta, logN)

	sum := make([]field.Element, 2)
	out0 := make([][]field.Element, 1<<logN)
	out1 := make([][]field.Element, 1<<logN)

	for i := 0; i < len(out0); i++ {
		out0[i] = make([]field.Element, 2)
		out1[i] = make([]field.Element, 2)
	}

	EvalFull(a, logN, out0)
	EvalFull(b, logN, out1)

	zero := field.Zero()
	for i := uint64(0); i < (uint64(1) << logN); i++ {
		for j := 0; j < 2; j++ {
			sum[j].Add(&out0[i][j], &out1[i][j])
		}

		//log.Printf("%v %v %v %v", i, alpha, beta.String(), sum.String())
		if i != alpha && (!sum[0].Equal(&zero) || !sum[1].Equal(&zero)) {
			test.Fail()
		}

		if i == alpha && (!sum[0].Equal(&beta[0]) || !sum[1].Equal(&beta[1])) {
			test.Fail()
		}
	}
}

func TestEvalFullPartial(test *testing.T) {
	logN := uint64(9)
	alpha := uint64(123)
	beta := make([]field.Element, 2)
	beta[0].SetUint64(7613)
	beta[1].SetUint64(991)

	a, b := Gen(alpha, beta, logN)

	sum := make([]field.Element, 2)

	outlen := 278
	out0 := make([][]field.Element, outlen)
	out1 := make([][]field.Element, outlen)

	for i := 0; i < len(out0); i++ {
		out0[i] = make([]field.Element, 2)
		out1[i] = make([]field.Element, 2)
	}

	EvalFull(a, logN, out0)
	EvalFull(b, logN, out1)

	zero := field.Zero()
	for i := uint64(0); i < uint64(outlen); i++ {
		for j := 0; j < 2; j++ {
			sum[j].Add(&out0[i][j], &out1[i][j])
		}

		//log.Printf("%v %v %v %v", i, alpha, beta.String(), sum.String())
		if i != alpha && (!sum[0].Equal(&zero) || !sum[1].Equal(&zero)) {
			test.Fail()
		}

		if i == alpha && (!sum[0].Equal(&beta[0]) || !sum[1].Equal(&beta[1])) {
			test.Fail()
		}
	}
}
*/
