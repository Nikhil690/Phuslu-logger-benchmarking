package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global logger variables
var (
	Log         *zap.Logger
	NfLog         *zap.SugaredLogger
	MainLog       *zap.SugaredLogger
	InitLog       *zap.SugaredLogger
	CfgLog        *zap.SugaredLogger
	CtxLog        *zap.SugaredLogger
	GinLog        *zap.SugaredLogger
	SBILog        *zap.SugaredLogger
	AmPolicyLog   *zap.SugaredLogger
	BdtPolicyLog  *zap.SugaredLogger
	ConsLog       *zap.SugaredLogger
	CallLog       *zap.SugaredLogger
	OamLog        *zap.SugaredLogger
	PolicyAuthLog *zap.SugaredLogger
	ProcLog       *zap.SugaredLogger
	SmPolicyLog   *zap.SugaredLogger
	UtilLog       *zap.SugaredLogger
)

// customLevelEncoder formats log levels with colors and even spacing
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorStart, colorEnd, levelStr string

	switch level {
	case zapcore.DebugLevel:
		colorStart, colorEnd, levelStr = "\033[36m", "\033[0m", "DEBUG" // Cyan
	case zapcore.InfoLevel:
		colorStart, colorEnd, levelStr = "\033[32m", "\033[0m", "INFO " // Green
	case zapcore.WarnLevel:
		colorStart, colorEnd, levelStr = "\033[33m", "\033[0m", "WARN " // Yellow
	case zapcore.ErrorLevel:
		colorStart, colorEnd, levelStr = "\033[31m", "\033[0m", "ERR  " // Red
	case zapcore.FatalLevel:
		colorStart, colorEnd, levelStr = "\033[41m\033[37m", "\033[0m", "FATAL" // White on Red
	default:
		colorStart, colorEnd, levelStr = "", "", level.CapitalString()
	}

	enc.AppendString(colorStart + levelStr + colorEnd)
}

// customComponentEncoder formats the component field for logs
func customComponentEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-5s", fmt.Sprintf("%s", loggerName)))
}

// Initialize initializes the global logger and component loggers
func Initialize(logLevel zapcore.Level) {
	if Log != nil {
		return
	}

	customEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "timestamp",
		LevelKey:         "level",
		CallerKey:        "caller", // Shows file:line
		MessageKey:       "message",
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 | 15:04:05.000"),
		EncodeLevel:      customLevelEncoder,         // Add colors to log levels
		EncodeCaller:     zapcore.ShortCallerEncoder, // Shows file:line
		NameKey:          "component",
		EncodeName:       customComponentEncoder, // Add component field inline
		ConsoleSeparator: " | ",
	})

	output := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(customEncoder, output, logLevel)

	Log = zap.New(core)
	sugaredLogger := Log.Sugar()

	// Initialize component loggers
	MainLog = sugaredLogger.Named("MAIN")
	NfLog = sugaredLogger.Named("NF")
	InitLog = sugaredLogger.Named("INIT")
	CfgLog = sugaredLogger.Named("CFG")
	CtxLog = sugaredLogger.Named("CTX")
	GinLog = sugaredLogger.Named("GIN")
	SBILog = sugaredLogger.Named("SBI")
	ProcLog = sugaredLogger.Named("PROC")
	UtilLog = sugaredLogger.Named("UTIL")
	ConsLog = sugaredLogger.Named("CONS")
	CallLog = sugaredLogger.Named("CALL")
	OamLog = sugaredLogger.Named("OAM")
	PolicyAuthLog = sugaredLogger.Named("POLICY_AUTH")
	AmPolicyLog = sugaredLogger.Named("AM_POLICY")
	BdtPolicyLog = sugaredLogger.Named("BDT_POLICY")
	SmPolicyLog = sugaredLogger.Named("SM_POLICY")
}

// Logger provides access to the global structured logger
func Logger() *zap.Logger {
	if Log == nil {
		panic("Logger is not initialized. Call Initialize first.")
	}
	return Log
}

// Sugar provides access to the global sugared logger
func Sugar() *zap.SugaredLogger {
	if Log == nil {
		panic("Logger is not initialized. Call Initialize first.")
	}
	return Log.Sugar()
}