package main

import (
	"github.com/Nikhil690/Phuslu-logger-benchmarking/logger"
	"fmt"
	"sync"
	"testing"
	"time"
)

// LogMessage struct (reusable via sync.Pool)
type LogMessage_test struct {
	Level string
	Text  string
	Extra string
}

// Create a sync.Pool to manage LogMessage objects
var logMessagePool_test = sync.Pool{
	New: func() any {
		return &LogMessage_test{} // Allocate a new object when needed
	},
}

// Buffered log channel (size = 100)
var logCh_test = make(chan *LogMessage_test, 10000)

// logProcessor runs in a separate Goroutine, reading from logCh_test
func logProcessor_test() {
	for lm := range logCh_test {
		switch lm.Level {
		case "INFO":
			if lm.Extra == "" {
				logger.Flash.Info().Msg(lm.Text)
			} else {
				logger.Flash.Info().Str("", lm.Extra).Msg(lm.Text)
			}
		case "ERROR":
			logger.Flash.Error().Msg(lm.Text)
		default:
			logger.Flash.Debug().Msg(lm.Text)
		}
		// Return the log message to the pool for reuse
		logMessagePool_test.Put(lm)
	}
}

// doSomeTask simulates work and logs using the pool
func doSomeTask_test() (time.Duration, time.Duration) {
	start := time.Now()

	// Log before starting
	lm := logMessagePool_test.Get().(*LogMessage_test)
	lm.Level = "INFO"
	lm.Text = "Starting doSomeTask..."
	lm.Extra = "func/doSomeTask"
	select {
	case logCh_test <- lm:
	default: // Drop the message if the channel is full
		logMessagePool_test.Put(lm) // Reuse it
	}

	// Simulate intermediate processing steps
	loop_time := time.Now()
	for i := 1; i <= 100000; i++ {
		time.Sleep(1000 * time.Nanosecond)

		lm := logMessagePool_test.Get().(*LogMessage_test)
		lm.Level = "INFO"
		lm.Text = fmt.Sprintf("Processing step %d/3 in doSomeTask...", i)
		lm.Extra = "func/doSomeTask"
		select {
		case logCh_test <- lm:
		default:
			logMessagePool_test.Put(lm)
		}
	}
	loop_end := time.Since(loop_time)

	// Simulate task execution

	// Log completion
	lm = logMessagePool_test.Get().(*LogMessage_test)
	lm.Level = "INFO"
	lm.Text = "doSomeTask finished successfully"
	lm.Extra = "func/doSomeTask"
	select {
	case logCh_test <- lm:
	default:
		logMessagePool_test.Put(lm)
	}

	// Log execution time at exit
	lm = logMessagePool_test.Get().(*LogMessage_test)
	lm.Level = "INFO"
	lm.Text = "doSomeTask execution finish"
	lm.Extra = "func/doSomeTask"
	select {
	case logCh_test <- lm:
	default:
		logMessagePool_test.Put(lm)
	}
	return time.Since(start), loop_end
}

// TestDoSomeTaskChannel sets up the channel-based logging
// and verifies doSomeTask takes ~2s, logging results via the channel.
func TestDoSomeTaskChannel(t *testing.T) {
	logger.Ready() // Initialize logger

	// Start the logProcessor Goroutine
	go logProcessor_test()

	// Start timing
	elapsed_time, loop_time := doSomeTask_test()
	time.Sleep(5 * time.Second)

	total := elapsed_time - loop_time

	// Get a pooled log message and send it to the channel
	lm := logMessagePool_test.Get().(*LogMessage_test)
	lm.Level = "INFO"
	lm.Text = "doSomeTask took " + total.String()
	lm.Extra = ""
	logCh_test <- lm

	// Allow logProcessor to consume all messages before closing
	time.Sleep(200 * time.Millisecond)
	close(logCh_test) // Close channel to signal logProcessor to exit
}
