package main

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
	"bench/logger"
)

type result struct {
	name     string
	duration time.Duration
}

func main() {
	// Store benchmark results for each writer
	var results []result
	initLogger()

	s ,e := benchmarknormal()
	results = append(results, result{
		name:     "Benchmark",
		duration: e.Sub(s),
	})
	s, e = benchmarkformatted()
	results = append(results, result{
		name:     "Benchmarkformatted",
		duration: e.Sub(s),
	})
	// s, e := consoleWriter()
	// results = append(results, result{
	// 	name:     "ConsoleWriter",
	// 	duration: e.Sub(s),
	// })

	// s, e = ioWriter()
	// results = append(results, result{
	// 	name:     "IOWriter",
	// 	duration: e.Sub(s),
	// })

	// s, e = fileWriter()
	// results = append(results, result{
	// 	name:     "FileWriter",
	// 	duration: e.Sub(s),
	// })

	// s, e = asyncWriter()
	// results = append(results, result{
	// 	name:     "AsyncWriter",
	// 	duration: e.Sub(s),
	// })

	// s, e = syslogWriter()
	// results = append(results, result{
	// 	name:     "SyslogWriter",
	// 	duration: e.Sub(s),
	// })

	// s, e = linuxJournalWriter()
	// results = append(results, result{
	// 	name:     "LinuxJournalWriter",
	// 	duration: e.Sub(s),
	// })

	// Print summary table
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

func initLogger() {
	logger.Initialize(zapcore.InfoLevel)
}

func benchmarkformatted() (time.Time, time.Time) {
	start := time.Now()
	for i := 0; i < 500000; i++ {
		logger.MainLog.Infof("Packet processed successfully", "UE_ID", 1001, "iteration", i)
	}
	end := time.Now()
	return start, end
}

func benchmarknormal() (time.Time, time.Time) {
	start := time.Now()
	for i := 0; i < 500000; i++ {
		logger.MainLog.Info("Packet processed successfully", "UE_ID", 1001, "iteration", i)
	}
	end := time.Now()
	return start, end
}
// 1) ConsoleWriter
// func consoleWriter() (time.Time, time.Time) {
// 	logger := log.Logger{
// 		Writer: &log.ConsoleWriter{
// 			ColorOutput: true, // Colorized console output
// 		},
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()
// 	return start, end
// }

// // 2) IOWriter (writes to Stdout)
// func ioWriter() (time.Time, time.Time) {
// 	logger := log.Logger{
// 		Writer: &log.IOWriter{
// 			Writer: os.Stdout,
// 		},
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()
// 	return start, end
// }

// // 3) FileWriter (writes to a file)
// func fileWriter() (time.Time, time.Time) {
// 	logger := log.Logger{
// 		Writer: &log.FileWriter{
// 			Filename: "test.log",
// 		},
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()
// 	return start, end
// }

// // 4) AsyncWriter (non-blocking logging to file)
// func asyncWriter() (time.Time, time.Time) {
// 	// Wrap a FileWriter in an AsyncWriter
// 	asyncW := &log.AsyncWriter{
// 		Writer: &log.FileWriter{
// 			Filename: "test_async.log",
// 		},
// 		// Note: Old fields like QueueSize, FlushInterval are no longer present in phuslu/log.
// 	}

// 	logger := log.Logger{
// 		Writer: asyncW,
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()

// 	// Give the async logger some time to flush (important for huge loops)
// 	time.Sleep(200 * time.Millisecond)

// 	// If you need an explicit flush/close, do it here:
// 	// _ = logger.Close() // or asyncW.Close() if implemented in newer phuslu/log

// 	return start, end
// }

// // 5) SyslogWriter (centralized system logging)
// func syslogWriter() (time.Time, time.Time) {
// 	// Priority can be LOG_DEBUG, LOG_INFO, LOG_ERR, etc.
// 	// Tag is a string to identify your application in the syslog
// 	sysWriter := &log.SyslogWriter{
// 		// Priority: syslog.LOG_INFO,
// 		Tag: "5G-Core",
// 	}

// 	logger := log.Logger{
// 		Writer: sysWriter,
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()

// 	return start, end
// }

// // 6) Linux JournalWriter (for journald)
// func linuxJournalWriter() (time.Time, time.Time) {
// 	journalW := &log.JournalWriter{}

// 	logger := log.Logger{
// 		Writer: journalW,
// 	}

// 	start := time.Now()
// 	for i := 0; i < 1_000_000; i++ {
// 		logger.Info().
// 			Str("component", "5G-Core").
// 			Int("UE_ID", 1001).
// 			Int("iteration", i).
// 			Dur("processing_time", time.Since(start)).
// 			Msg("Packet processed successfully")
// 	}
// 	end := time.Now()

// 	return start, end
// }
