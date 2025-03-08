package main

import (
	"bench/logger"
	"fmt"
	"time"

	"github.com/phuslu/log"
)

type result struct {
	name     string
	duration time.Duration
}

func main() {
	// 1) Initialize the logger
	results := []result{}
	logger.Ready()

	s, e := benchmarknormal(logger.Lopu)
	results = append(results, result{
		name:     "Benchmark",
		duration: e.Sub(s),
	})
	s, e = benchmarkformatted(logger.Lopu)
	results = append(results, result{
		name:     "Benchmarkformatted",
		duration: e.Sub(s),
	})

	fmt.Println("\nSummary of Logging Performance")
	fmt.Println("------------------------------------------------")
	fmt.Printf("| %-20s | %-15s |\n", "Writer", "Duration (s)")
	fmt.Println("------------------------------------------------")
	for _, r := range results {
		// Convert duration to seconds with 6 decimal places
		secs := float64(r.duration.Milliseconds()) / 1000.0
		fmt.Printf("| %-20s | %-15.6f |\n", r.name, secs)
	}
	fmt.Println("------------------------------------------------")

}

func benchmarknormal(logger log.Logger) (time.Time, time.Time) {
	start := time.Now()
	for i := 0; i < 500000; i++ {
		logger.Info().Msg("Packet processed successfully UE_ID 1001 " + "iteration")
	}
	end := time.Now()
	return start, end
}

func benchmarkformatted(logger log.Logger) (time.Time, time.Time) {
	start := time.Now()
	for i := 0; i < 500000; i++ {
		logger.Info().Msgf("Packet processed successfully UE_ID 1001 iteration: %d ", i)
	}
	end := time.Now()
	return start, end
}
// ---------------------------------------------------------------------

// func main() {
// 	var results []result
// 	logger.Initialize(log.InfoLevel)

// 	// 2. Create "sub-loggers" for each component
// 	mainLog := logger.NewComponentLogger("MAIN")
// 	// nfLog := logger.NewComponentLogger("NF")
// 	// initLog := logger.NewComponentLogger("INIT")

// 	// 3. Use them
// 	s, e := benchmarknormal(mainLog)
// 	results = append(results, result{
// 		name:     "Benchmarknormal",
// 		duration: e.Sub(s),
// 	})

// 	s, e = benchmarkformatted(mainLog)
// 	results = append(results, result{
// 		name:     "Benchmarkformatted",
// 		duration: e.Sub(s),
// 	})

// 	fmt.Println("\nSummary of Logging Performance")
// 	fmt.Println("------------------------------------------------")
// 	fmt.Printf("| %-20s | %-15s |\n", "Writer", "Duration (s)")
// 	fmt.Println("------------------------------------------------")
// 	for _, r := range results {
// 		// Convert duration to seconds with 6 decimal places
// 		secs := float64(r.duration.Milliseconds()) / 1000.0
// 		fmt.Printf("| %-20s | %-15.6f |\n", r.name, secs)
// 	}
// 	fmt.Println("------------------------------------------------")
// }

// func benchmarknormal(logger *logger.SubLogger) (time.Time, time.Time) {
// 	start := time.Now()
// 	for i := 0; i < 500000; i++ {
// 		logger.Info().Msg("Packet processed successfully UE_ID 1001 " + "iteration")
// 	}
// 	end := time.Now()
// 	return start, end
// }

// func benchmarkformatted(logger *logger.SubLogger) (time.Time, time.Time) {
// 	start := time.Now()
// 	for i := 0; i < 500000; i++ {
// 		logger.Info().Msgf("Packet processed successfully UE_ID 1001 iteration: %d ", i)
// 	}
// 	end := time.Now()
// 	return start, end
// }
