package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/si-co/vpir-code/lib/client"
	"github.com/si-co/vpir-code/lib/database"
	"github.com/si-co/vpir-code/lib/monitor"
	"github.com/si-co/vpir-code/lib/server"
	"github.com/si-co/vpir-code/lib/utils"
	"github.com/stretchr/testify/require"
)

func TestAmplifyOneMbEncoding(t *testing.T) {
	threshold := 8
	dbLen := 1024 * 1024 // dbLen is specified in bits
	db := database.CreateRandomBinaryLWEWithLength(utils.RandomPRG(), dbLen)
	p := utils.ParamsWithDatabaseSize(db.Info.NumRows, db.Info.NumColumns)

	retrieveBlocksAmplify(t, db, p, threshold, "TestAmplifyOneMb")
}

func TestAmplifyOneMb(t *testing.T) {
	threshold := 8
	dbLen := 1024 * 1024 // dbLen is specified in bits
	db := database.CreateRandomBinaryLWEWithLength(utils.RandomPRG(), dbLen)
	p := utils.ParamsWithDatabaseSize(db.Info.NumRows, db.Info.NumColumns)

	retrieveBlocksAmplify(t, db, p, threshold, "TestAmplifyOneMb")
}

func retrieveBlocksAmplify(t *testing.T, db *database.LWE, params *utils.ParamsLWE, threshold int, testName string) {
	c := client.NewAmplify(utils.RandomPRG(), &db.Info, params, threshold)
	s := server.NewAmplify(db)

	totalTimer := monitor.NewMonitor()
	repetitions := 10
	for k := 0; k < repetitions; k++ {
		i := rand.Intn(params.L * params.M)
		query, err := c.QueryBytes(i)
		require.NoError(t, err)

		a, err := s.AnswerBytes(query)
		require.NoError(t, err)

		res, err := c.ReconstructBytes(a)
		require.NoError(t, err)
		require.Equal(t, db.Matrix.Get(utils.VectorToMatrixIndices(i, db.Info.NumColumns)), res)
	}
	fmt.Printf("TotalCPU time %s: %.1fms\n", testName, totalTimer.Record())

}

func retrieveBlocksAmplifyNoBytes(t *testing.T, db *database.LWE, params *utils.ParamsLWE, threshold int, testName string) {
	c := client.NewAmplify(utils.RandomPRG(), &db.Info, params, threshold)
	s := server.NewAmplify(db)

	totalTimer := monitor.NewMonitor()
	repetitions := 10
	for k := 0; k < repetitions; k++ {
		i := rand.Intn(params.L * params.M)
		ii, jj := utils.VectorToMatrixIndices(i, db.Info.NumColumns)
		query := c.Query(ii, jj)

		a := s.Answer(query)

		res, err := c.Reconstruct(a)
		require.NoError(t, err)
		require.Equal(t, db.Matrix.Get(ii, jj), res)
	}
	fmt.Printf("TotalCPU time %s: %.1fms\n", testName, totalTimer.Record())

}
