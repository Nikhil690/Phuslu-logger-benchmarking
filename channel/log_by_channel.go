package main

import (
	"github.com/Nikhil690/bench/logger"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// LogMessage struct (reusable via sync.Pool)
type LogMessage struct {
	Level string
	Text  string
	Extra string
}

// Create a sync.Pool to manage LogMessage objects
var logMessagePool = sync.Pool{
	New: func() any {
		return make([]*LogMessage, 0, 10) // Pre-allocate space for 10 logs per slice
	},
}

// Buffered log channel (size = 100)
var logCh = make(chan []*LogMessage, 64) // Enough to take 3500 logs at once

// logProcessor runs in a separate Goroutine, reading from logCh
func logProcessor() {
	for batch := range logCh { // Dequeue FIFO batch of logs
		for _, lm := range batch { // Process each log message in the batch
			switch lm.Level {
			case "INFO":
				if lm.Extra == "" {
					logger.Flash.Info().Str("", logger.MainLog).Msg(lm.Text)
				} else {
					logger.Flash.Info().Str("", logger.MainLog).Str("", lm.Extra).Msg(lm.Text)
				}
			case "ERROR":
				logger.Flash.Error().Msg(lm.Text)
			default:
				logger.Flash.Debug().Msg(lm.Text)
			}
		}
		// Return batch to sync.Pool for reuse
		logMessagePool.Put(batch) // Need to do something it is giving this error argument should be pointer-like to avoid allocations (SA6002)go-staticcheck
	}
}

// doSomeTask simulates work and logs using the pool
func doSomeTask(id int) time.Duration {
	start := time.Now()

	// Get batch from sync.Pool
	batch := logMessagePool.Get().([]*LogMessage)
	batch = batch[:0] // Reset slice length but keep capacity

	// Initial log entry
	batch = append(batch, &LogMessage{Level: logger.InfoLevel, Text: "Starting doSomeTask..." + fmt.Sprint(id), Extra: "func/doSomeTask"})

	// Loop for processing steps
	for i := 1; i <= 5; i++ {
		batch = append(batch, &LogMessage{
			Level: logger.InfoLevel,
			Text:  fmt.Sprintf("Processing step %d/5 in doSomeTask %s", i, start),
			Extra: "func/doSomeTask",
		})
	}

	// Final log entry
	batch = append(batch, &LogMessage{Level: logger.InfoLevel, Text: "doSomeTask execution time: ", Extra: "func/doSomeTask"})

	// Send batch to channel
	select {
	case logCh <- batch:
	default:
		logMessagePool.Put(batch) // this as well
	}
	elapsed := time.Since(start)
	return elapsed
}

func main() {
	logger.Ready() // Initialize logger
	id := 1

	// Start the logProcessor Goroutine
	go logProcessor()

	// Capture memory usage before execution
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore) // 📌 Read memory before execution

	// Start timing
	start := time.Now()
	var total time.Duration
	for range 500 {
		elapsed := doSomeTask(id)
		time.Sleep(1 * time.Microsecond)
		id++
		total += elapsed
	}
	avg := time.Duration(int64(total) / int64(50))
	time.Sleep(5 * time.Second)
	fmt.Println("Average execution time:", avg)
	elapsed := time.Since(start)
	time.Sleep(5 * time.Microsecond)

	// Capture memory usage after execution
	runtime.ReadMemStats(&memAfter) // 📌 Read memory after execution

	// Print elapsed time
	print(elapsed.String())

	// 📌 Print memory allocation details
	fmt.Println("\n📊 Memory Usage Report:")
	fmt.Printf("Heap Alloc: %v KB → %v KB\n", memBefore.HeapAlloc/1024, memAfter.HeapAlloc/1024)
	fmt.Printf("Total Alloc: %v KB → %v KB\n", memBefore.TotalAlloc/1024, memAfter.TotalAlloc/1024)
	fmt.Printf("Heap Objects: %v → %v\n", memBefore.HeapObjects, memAfter.HeapObjects)
	fmt.Printf("GC Cycles: %v → %v\n", memBefore.NumGC, memAfter.NumGC)
	fmt.Println("---------------------------")

	// Allow logProcessor to consume all messages before closing
	time.Sleep(200 * time.Millisecond)
	close(logCh) // Close channel to signal logProcessor to exit
}
