package main

import (
	"bench/logger"
	"fmt"
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

// Buffered log channel (size = 100)
var logCh = make(chan *LogMessage, 10000000)

// logProcessor runs in a separate Goroutine, reading from logCh
func logProcessor() {
	for lm := range logCh {
		switch lm.Level {
		case "INFO":
			logger.Lopu.Info().Msg(lm.Text)
		case "ERROR":
			logger.Lopu.Error().Msg(lm.Text)
		default:
			logger.Lopu.Debug().Msg(lm.Text)
		}
		// Return the log message to the pool for reuse
		logMessagePool.Put(lm)
	}
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
	default: // Drop the message if the channel is full
		logMessagePool.Put(lm) // Reuse it
	}

	// Simulate intermediate processing steps
	for i := 1; i <= 100000; i++ {
		// time.Sleep(500 * time.Millisecond)

		lm := logMessagePool.Get().(*LogMessage)
		lm.Level = "INFO"
		lm.Text = fmt.Sprintf("Processing step %d/3 in doSomeTask...", i)
		select {
		case logCh <- lm:
		default:
			logMessagePool.Put(lm)
		}
	}

	// Simulate task execution
	time.Sleep(2 * time.Second)

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


// TestDoSomeTaskChannel sets up the channel-based logging
// and verifies doSomeTask takes ~2s, logging results via the channel.
func TestDoSomeTaskChannel(t *testing.T) {
	logger.Ready() // Initialize logger

	// Start the logProcessor Goroutine
	go logProcessor()

	// Start timing
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
	close(logCh) // Close channel to signal logProcessor to exit
}
