package main

import (
	"bench/logger"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// LogMessage struct (reusable via sync.Pool)
type LogMessage struct {
	Level string
	Text  string
}

// Create a sync.Pool to manage LogMessage objects
var logMessagePool = sync.Pool{
	New: func() any {
		return &LogMessage{} // Allocate a new object when needed
	},
}

// Buffered log channel (size = 10M)
var logCh = make(chan *LogMessage, 100)

// logProcessor runs in a separate Goroutine, reading from logCh
func logProcessor() {
	for lm := range logCh {
		switch lm.Level {
		case "INFO":
			logger.Flash.Info().Msg(lm.Text)
		case "ERROR":
			logger.Flash.Error().Msg(lm.Text)
		default:
			logger.Flash.Debug().Msg(lm.Text)
		}
		// Return the log message to the pool for reuse
		logMessagePool.Put(lm)
	}
}

// MeasureMemory prints memory usage before & after execution
func MeasureMemory(action func()) {
	var mBefore, mAfter runtime.MemStats

	// Capture memory usage before running the function
	runtime.ReadMemStats(&mBefore)

	// Execute the function
	action()

	// Capture memory usage after running the function
	runtime.ReadMemStats(&mAfter)

	// Print memory usage results
	fmt.Printf("\nMemory Usage Report:\n")
	fmt.Printf("Heap Allocated: %v KB → %v KB\n", mBefore.HeapAlloc/1024, mAfter.HeapAlloc/1024)
	fmt.Printf("Total Alloc: %v KB → %v KB\n", mBefore.TotalAlloc/1024, mAfter.TotalAlloc/1024)
	fmt.Printf("Heap Objects: %v → %v\n", mBefore.HeapObjects, mAfter.HeapObjects)
	fmt.Printf("GC Runs: %v → %v\n", mBefore.NumGC, mAfter.NumGC)
	fmt.Println("---------------------------")
}

// doSomeTask simulates work and logs using the pool
func doSomeTask() {
	start := time.Now()

	// Log before starting
	lm := logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = "Starting doSomeTask..."
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}

	// Simulate intermediate processing steps
	for i := 1; i <= 20; i++ {
		time.Sleep(1 * time.Microsecond)
		lm := logMessagePool.Get().(*LogMessage)
		lm.Level = "INFO"
		lm.Text = fmt.Sprintf("Processing step %d in doSomeTask...", i)
		select {
		case logCh <- lm:
		default:
			logMessagePool.Put(lm)
		}
	}

	// Simulate task execution
	// Log completion
	lm = logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = "doSomeTask finished successfully"
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}

	// Log execution time at exit
	elapsed := time.Since(start)
	lm = logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = fmt.Sprintf("doSomeTask execution time: %v", elapsed)
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}
}

// TestDoSomeTaskChannel with Memory Profiling
func TestDoSomeTaskChannel(t *testing.T) {
	logger.Ready() // Initialize logger

	// Start the logProcessor Goroutine
	go logProcessor()

	// Measure memory usage during execution
	MeasureMemory(func() {
		start := time.Now()
		doSomeTask()
		elapsed := time.Since(start)
		time.Sleep(10 * time.Second)

		// Get a pooled log message and send it to the channel
		lm := logMessagePool.Get().(*LogMessage)
		lm.Level = "INFO"
		lm.Text = "doSomeTask took " + elapsed.String()
		logCh <- lm

		// Check if task was too fast
		if elapsed < 2*time.Second {
			lm := logMessagePool.Get().(*LogMessage)
			lm.Level = "ERROR"
			lm.Text = "Task finished too quickly: " + elapsed.String()
			logCh <- lm
		}

		// Allow logProcessor to consume all messages before closing
		time.Sleep(200 * time.Millisecond)
		close(logCh)
	})
}
