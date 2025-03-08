package logger

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/phuslu/log"
)

// -------------------------------------------------------------
// 1) ANSI Color Codes & Minimal Level-To-Color Logic
// -------------------------------------------------------------
var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[97m"
)

// levelToColor maps a phuslu/log level to an ANSI color + short label
func levelToColor(level string) (colorCode, label string) {
	switch strings.ToLower(level) {
	case "debug":
		return colorCyan, "DEBUG"
	case "info":
		return colorGreen, "INFO "
	case "warn":
		return colorYellow, "WARN "
	case "error":
		return colorRed, "ERR  "
	case "fatal":
		// White text on a red background
		return "\033[41m\033[37m", "FATAL"
	default:
		return colorWhite, strings.ToUpper(level)
	}
}

// -------------------------------------------------------------
// 2) Custom Console Formatter (No extra time parsing or fields)
// -------------------------------------------------------------
func customConsoleFormatter(w io.Writer, args *log.FormatterArgs) (int, error) {
	// phuslu/log sets args.Time to an RFC3339Nano string by default (e.g. "2025-03-08T12:34:56.789Z").
	// We skip re-parsing for performance; we’ll just print it as-is.

	colorCode, levelLabel := levelToColor(args.Level)

	// Format: "timestamp | colored-level | message"
	// Example:
	// 2025-03-08T12:34:56.789Z | [GREEN]INFO  | Hello World
	return fmt.Fprintf(w, "%s | \033[0m%s%s\033[0m | %s\n",
		args.Time,    // e.g. "2025-03-08T12:34:56.789Z"
		colorCode,    // e.g. "\033[32m"
		levelLabel,   // e.g. "INFO "
		args.Message, // e.g. "Hello World"
	)
}

// -------------------------------------------------------------
// 3) Global Logger + Initialization
// -------------------------------------------------------------
var (
	globalLogger log.Logger
	initOnce     sync.Once
)

// Initialize sets up the global logger just once.
func Initialize(level log.Level) {
	initOnce.Do(func() {
		consoleWriter := &log.ConsoleWriter{
			ColorOutput: false, // We'll manually colorize in customConsoleFormatter
			Formatter:   customConsoleFormatter,
		}

		globalLogger = log.Logger{
			Level:  level,         // e.g. log.InfoLevel, log.DebugLevel, etc.
			Writer: consoleWriter, // Output to stdout in our custom format
			// Caller: 1,           // If you ever want file:line info in args.Caller
		}
	})
}

// Logger returns the global logger.
func Logger() log.Logger {
	return globalLogger
}

var (
	Lopu log.Logger
)

func Ready() {
	Initialize(log.InfoLevel)

	// 2) Use the global logger
	Lopu = Logger()
}
// ---------------------------------------------------------------------------------
// color codes for levels
// var (
// 	colorReset   = "\033[0m"
// 	colorRed     = "\033[31m"
// 	colorGreen   = "\033[32m"
// 	colorYellow  = "\033[33m"
// 	colorBlue    = "\033[34m"
// 	colorCyan    = "\033[36m"
// 	colorWhite   = "\033[97m"
// 	colorMagenta = "\033[35m"
// )

// // create a global logger + a custom console writer
// var (
// 	globalLogger log.Logger
// 	initOnce     sync.Once
// )

// // SubLogger wraps a phuslu/log Logger and automatically adds "component" to each entry.
// type SubLogger struct {
// 	baseLogger *log.Logger
// 	component  string
// }

// // Info() returns a new log.Entry with the component field set.
// func (s *SubLogger) Info() *log.Entry {
// 	return s.baseLogger.Info().Str("component", s.component)
// }

// func (s *SubLogger) Warn() *log.Entry {
// 	return s.baseLogger.Warn().Str("component", s.component)
// }

// func (s *SubLogger) Error() *log.Entry {
// 	return s.baseLogger.Error().Str("component", s.component)
// }

// func (s *SubLogger) Debug() *log.Entry {
// 	return s.baseLogger.Debug().Str("component", s.component)
// }

// func (s *SubLogger) Fatal() *log.Entry {
// 	return s.baseLogger.Fatal().Str("component", s.component)
// }

// // levelToColor maps a phuslu/log level to an ANSI color + short label (mimicking your Zap scheme).
// func levelToColor(level string) (colorCode, levelStr string) {
// 	switch strings.ToLower(level) {
// 	case "debug":
// 		return colorCyan, "DEBUG"
// 	case "info":
// 		return colorGreen, "INFO "
// 	case "warn":
// 		return colorYellow, "WARN "
// 	case "error":
// 		return colorRed, "ERR  "
// 	case "fatal":
// 		// white text on red background
// 		return "\033[41m\033[37m", "FATAL"
// 	default:
// 		return colorWhite, strings.ToUpper(level)
// 	}
// }

// // customConsoleFormatter replicates the line format:
// // timestamp | colored-level | caller | component | message
// func customConsoleFormatter(w io.Writer, args *log.FormatterArgs) (int, error) {
// 	// Reformat the timestamp from RFC3339Nano to "2006-01-02 | 15:04:05.000"
// 	parsed, err := time.Parse(time.RFC3339Nano, args.Time)
// 	if err == nil {
// 		args.Time = parsed.Format("2006-01-02 | 15:04:05.000")
// 	}

// 	// colorize the level
// 	colorCode, levelStr := levelToColor(args.Level)

// 	// Extract "component" from KeyValues
// 	var component string
// 	for _, kv := range args.KeyValues {
// 		if kv.Key == "component" {
// 			component = kv.Value
// 			break
// 		}
// 	}

// 	// If you want to display file:line, phuslu/log must be configured with logger.Caller=1 or so.
// 	// That would populate args.Caller. If not set, args.Caller might be empty.
// 	return fmt.Fprintf(w,
// 		"%s | \033[0m%s%s\033[0m | %s | %-5s | %s\n",
// 		args.Time,           // e.g. "2025-03-08 | 12:34:56.789"
// 		colorCode, levelStr, // e.g. "[green]INFO "
// 		args.Caller,  // e.g. "main.go:42"
// 		component,    // e.g. "MAIN"
// 		args.Message, // e.g. "Initialization complete"
// 	)
// }

// // Initialize sets up the global logger with your console format, only once.
// func Initialize(level log.Level) {
// 	initOnce.Do(func() {
// 		consoleWriter := &log.ConsoleWriter{
// 			ColorOutput: false, // We'll handle color in customConsoleFormatter
// 			Formatter:   customConsoleFormatter,
// 		}

// 		// Build the global logger
// 		globalLogger = log.Logger{
// 			Level:  level, // e.g. log.InfoLevel
// 			Writer: consoleWriter,
// 			// Caller:  1,  // uncomment if you want file:line in args.Caller
// 		}
// 	})
// }

// // Logger returns the *global* logger (in case you need it).
// func Logger() log.Logger {
// 	return globalLogger
// }

// // NewComponentLogger returns a SubLogger that automatically attaches "component" to each log entry.
// func NewComponentLogger(name string) *SubLogger {
// 	return &SubLogger{
// 		baseLogger: &globalLogger,
// 		component:  name,
// 	}
// }
