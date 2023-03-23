package main

// Test suite for Merkle tree-based VPIR schemes. Only multi-bit schemes are
// implemented using this approach.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/si-co/vpir-code/lib/client"
	"github.com/si-co/vpir-code/lib/database"
	"github.com/si-co/vpir-code/lib/server"
	"github.com/si-co/vpir-code/lib/utils"
)

const (
	oneB            = 8
	oneKB           = 1024 * oneB
	oneMB           = 1024 * oneKB
	oneGB           = 1024 * oneMB
	testBlockLength = 16
)

var randomDB *database.DB

var DB_SIZE_EXPO = []uint{18, 20, 22, 24, 26}
var ITEM_SIZE_EXPO = []uint{4, 6, 8, 10}

var comm_file, _ = os.Create("./bench_comm.txt")
var mem_file, _ = os.Create("./bench_mem.txt")

func BenchmarkMerkle(b *testing.B) {
	for _, dbLenExpo := range DB_SIZE_EXPO {
		for _, itemLenExpo := range ITEM_SIZE_EXPO {
			runtime.GC()
			name := fmt.Sprintf("Merkle-2^%ddb-%db", dbLenExpo-itemLenExpo, itemLenExpo)
			comm_file.WriteString(name + " ")
			mem_file.WriteString(name + " ")
			b.Run(name, func(b *testing.B) {
				benchmarkMerkle(b, int(math.Pow(2, float64(dbLenExpo)))*8, int(math.Pow(2, float64(itemLenExpo))))
			})
			comm_file.WriteString("\n")
			mem_file.WriteString("\n")
		}
	}
}

func BenchmarkPIRPoint(b *testing.B) {
	for _, dbLenExpo := range DB_SIZE_EXPO {
		for _, itemLenExpo := range ITEM_SIZE_EXPO {
			runtime.GC()
			name := fmt.Sprintf("Normal-2^%ddb-%db", dbLenExpo-itemLenExpo, itemLenExpo)
			comm_file.WriteString(name + " ")
			mem_file.WriteString(name + " ")
			b.Run(name, func(b *testing.B) {
				benchmarkPIRPoint(b, int(math.Pow(2, float64(dbLenExpo)))*8, int(math.Pow(2, float64(itemLenExpo))))
			})
			comm_file.WriteString("\n")
			mem_file.WriteString("\n")
		}
	}
}

// This does not work for miraculous implementation of the merkle tree. It use the *32bit checksum* of a value as the key in a map to find its index. Too many items results in serious collision.
// And it will need more than 50G ram
// func BenchmarkMerkle28d16b(b *testing.B) {
// 	benchmarkMerkle(b, oneMB*256, 16)
// }

func benchmarkMerkle(b *testing.B, dbLen int, blockLen int) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	numServers := 2
	// since this scheme works on bytes, the bit size of one element is 8
	elemBitSize := 8
	numBlocks := dbLen / (elemBitSize * blockLen)
	nCols := int(math.Sqrt(float64(numBlocks)))
	nRows := numBlocks / nCols

	db := database.CreateRandomMerkle(utils.RandomPRG(), dbLen, nRows, blockLen)

	runtime.ReadMemStats(&m2)
	mem_file.WriteString(fmt.Sprintf("%dB ", (m2.Alloc - m1.Alloc)))
	retrieveBlocksMerkle(b, utils.RandomPRG(), db, numServers, numBlocks, "Merkle")
}

func retrieveBlocksMerkle(b *testing.B, rnd io.Reader, db *database.Bytes, numServers, numBlocks int, testName string) {
	c := client.NewPIR(rnd, &db.Info)
	servers := make([]*server.PIR, numServers)
	for i := range servers {
		servers[i] = server.NewPIR(db)
	}
	totalComm := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make([]byte, 4)
		binary.BigEndian.PutUint32(in, uint32(1))
		queries, _ := c.QueryBytes(in, numServers)
		answers := make([][]byte, numServers)
		for i, s := range servers {
			totalComm += len(queries[i])
			ans, _ := s.AnswerBytes(queries[i])
			totalComm += len(ans)
			answers[i] = ans
		}

		c.ReconstructBytes(answers)
	}
	totalComm /= b.N
	comm_file.WriteString(fmt.Sprintf("%dB ", totalComm))

}

// Test suite for classical PIR, used as baseline for the experiments.

func benchmarkPIRPoint(b *testing.B, dbLen int, blockLen int) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	elemBitSize := 8
	numBlocks := dbLen / (elemBitSize * blockLen)
	nCols := int(math.Sqrt(float64(numBlocks)))
	nRows := numBlocks / nCols

	// functions defined in vpir_test.go
	xofDB := utils.RandomPRG()
	xof := utils.RandomPRG()

	db := database.CreateRandomBytes(xofDB, dbLen, nRows, blockLen)

	runtime.ReadMemStats(&m2)
	mem_file.WriteString(fmt.Sprintf("%dB ", (m2.Alloc - m1.Alloc)))
	retrievePIRPoint(b, xof, db, numBlocks, "PIRPoint")
}

func retrievePIRPoint(b *testing.B, rnd io.Reader, db *database.Bytes, numBlocks int, testName string) {
	c := client.NewPIR(rnd, &db.Info)
	s0 := server.NewPIR(db)
	s1 := server.NewPIR(db)
	totalComm := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make([]byte, 4)
		binary.BigEndian.PutUint32(in, uint32(1))
		queries, _ := c.QueryBytes(in, 2)
		totalComm += len(queries[0]) + len(queries[1])
		a0, _ := s0.AnswerBytes(queries[0])
		a1, _ := s1.AnswerBytes(queries[1])
		totalComm += len(a0) + len(a1)
		answers := [][]byte{a0, a1}

		c.ReconstructBytes(answers)
	}
	totalComm /= b.N
	comm_file.WriteString(fmt.Sprintf("%dB ", totalComm))
}
