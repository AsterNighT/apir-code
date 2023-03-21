package main

// Test suite for Merkle tree-based VPIR schemes. Only multi-bit schemes are
// implemented using this approach.

import (
	"encoding/binary"
	"io"
	"math"
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

func BenchmarkMerkle(b *testing.B) {
	numServers := 2
	dbLen := oneMB * 128
	blockLen := testBlockLength
	// since this scheme works on bytes, the bit size of one element is 8
	elemBitSize := 8
	numBlocks := dbLen / (elemBitSize * blockLen)
	nCols := int(math.Sqrt(float64(numBlocks)))
	nRows := numBlocks / nCols

	db := database.CreateRandomMerkle(utils.RandomPRG(), dbLen, nRows, blockLen)

	retrieveBlocksMerkle(b, utils.RandomPRG(), db, numServers, numBlocks, "Merkle")
}

func retrieveBlocksMerkle(b *testing.B, rnd io.Reader, db *database.Bytes, numServers, numBlocks int, testName string) {
	c := client.NewPIR(rnd, &db.Info)
	servers := make([]*server.PIR, numServers)
	for i := range servers {
		servers[i] = server.NewPIR(db)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make([]byte, 4)
		binary.BigEndian.PutUint32(in, uint32(0))
		queries, _ := c.QueryBytes(in, numServers)

		answers := make([][]byte, numServers)
		for i, s := range servers {
			s.AnswerBytes(queries[i])
		}

		c.ReconstructBytes(answers)
	}
}

// Test suite for classical PIR, used as baseline for the experiments.

func BenchmarkPIRPoint(b *testing.B) {
	dbLen := oneMB
	blockLen := testBlockLength
	elemBitSize := 8
	numBlocks := dbLen / (elemBitSize * blockLen)
	nCols := int(math.Sqrt(float64(numBlocks)))
	nRows := numBlocks / nCols

	// functions defined in vpir_test.go
	xofDB := utils.RandomPRG()
	xof := utils.RandomPRG()

	db := database.CreateRandomBytes(xofDB, dbLen, nRows, blockLen)

	retrievePIRPoint(b, xof, db, numBlocks, "PIRPoint")
}

func retrievePIRPoint(b *testing.B, rnd io.Reader, db *database.Bytes, numBlocks int, testName string) {
	c := client.NewPIR(rnd, &db.Info)
	s0 := server.NewPIR(db)
	s1 := server.NewPIR(db)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make([]byte, 4)
		binary.BigEndian.PutUint32(in, uint32(0))
		queries, _ := c.QueryBytes(in, 2)

		a0, _ := s0.AnswerBytes(queries[0])
		a1, _ := s1.AnswerBytes(queries[1])

		answers := [][]byte{a0, a1}

		c.ReconstructBytes(answers)
	}
}
