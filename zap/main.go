package main

import (
	"bench/logger"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1) Initialize the logger
	// results := []result{}
	logger.Initialize(zapcore.InfoLevel)
	var id int
	var avg time.Duration
	var total time.Duration
	for range 500 {
		elapsed := benchmarknormal(logger.MainLog, id)
		total += elapsed
		id++
	}
	avg = time.Duration(int64(total) / int64(50))
	time.Sleep(5 * time.Second)
	logger.MainLog.Infof("Average execution time: %v", avg)
}

func benchmarknormal(logger *zap.SugaredLogger, id int) time.Duration {
	start := time.Now()
	logger.Infof("Starting doSomeTask... %d", id)
	for i := 0; i < 5; i++ {
		logger.Infof("Processing step %d/5 in doSomeTask %s", i, start.String())
	}
	logger.Info("doSomeTask execution time: ")
	elapsed := time.Since(start)
	return elapsed
}
