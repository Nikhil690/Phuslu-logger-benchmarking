package main

import (
	"github.com/Nikhil690/Phuslu-logger-benchmarking/logger"
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

// Buffered log channel (adjust size as needed)
var logCh = make(chan *LogMessage, 1000)

// logProcessor runs in a separate Goroutine, reading from logCh
func logProcessor(wg *sync.WaitGroup) {
	defer wg.Done() // Ensure proper cleanup

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
func doSomeTask(taskID int) {
	start := time.Now()

	// Log before starting
	lm := logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = fmt.Sprintf("Starting Task %d...", taskID)
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}

	// Simulate intermediate processing steps
	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Microsecond)
		lm := logMessagePool.Get().(*LogMessage)
		lm.Level = "INFO"
		lm.Text = fmt.Sprintf("Task %d: Processing step %d...", taskID, i)
		select {
		case logCh <- lm:
		default:
			logMessagePool.Put(lm)
		}
	}

	// Simulate task execution
	// time.Sleep(2 * time.Second)

	// Log completion
	lm = logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = fmt.Sprintf("Task %d finished successfully", taskID)
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}

	// Log execution time at exit
	elapsed := time.Since(start)
	lm = logMessagePool.Get().(*LogMessage)
	lm.Level = "INFO"
	lm.Text = fmt.Sprintf("Task %d execution time: %v", taskID, elapsed)
	select {
	case logCh <- lm:
	default:
		logMessagePool.Put(lm)
	}
}

// TestMultipleTasks runs 50 instances of doSomeTask in parallel
func TestMultipleTasks(t *testing.T) {
	logger.Ready() // Initialize logger

	// WaitGroups to track Goroutines
	var wg sync.WaitGroup

	// Start the log processor
	wg.Add(1)
	go logProcessor(&wg)

	// Measure memory usage while running 50 tasks
	MeasureMemory(func() {
		var taskWg sync.WaitGroup

		// Run 50 tasks in parallel
		taskCount := 50
		taskWg.Add(taskCount)
		for i := 1; i <= taskCount; i++ {
			go func(taskID int) {
				defer taskWg.Done()
				doSomeTask(taskID)
			}(i)
		}

		// Wait for all tasks to finish
		taskWg.Wait()
		fmt.Println("All tasks completed.")
	})

	// Allow logs to be processed before shutting down
	time.Sleep(1 * time.Second)

	// Close log channel and wait for logProcessor to exit
	close(logCh)
	wg.Wait()
}
