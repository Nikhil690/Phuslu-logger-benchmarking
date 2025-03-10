package main

import (
	"bench/logger"
	"time"
)

func main() {
	// 1) Initialize the logger
	// results := []result{}
	logger.Ready()
	var id int
	var avg time.Duration
	var total time.Duration
	for range 500 {
		elapsed := benchmarknormal(id)
		total += elapsed
		id++
	}
	avg = time.Duration(int64(total) / int64(50))
	time.Sleep(5 * time.Second)
	logger.Flash.Info().Msgf("Average execution time: %v", avg)
}

func benchmarknormal(id int) time.Duration {
	start := time.Now()
	logger.Flash.Info().Str("", logger.MainLog).Msgf("Starting doSomeTask... %d", id)
	for i := 0; i < 5; i++ {
		logger.Flash.Info().Str("", logger.MainLog).Msgf("Processing step %d/5 in doSomeTask %s", i, start.String())
	}
	logger.Flash.Info().Str("", logger.MainLog).Msg("doSomeTask execution time: ")
	elapsed := time.Since(start)
	return elapsed
}
