package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var Logger zerolog.Logger

func InitLog() {
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s |", i)
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}).
		Level(getLoggerLevel(os.Getenv("LOG_LEVEL"))).
		With().
		Timestamp().
		Caller().
		Str("go_version", runtime.Version()).
		Int("pid", os.Getpid()).
		Logger()
}

func getLoggerLevel(value string) zerolog.Level {
	switch value {
	case "DEBUG":
		return zerolog.DebugLevel
	case "TRACE":
		return zerolog.TraceLevel
	case "INFO":
		return zerolog.InfoLevel
	default:
		return zerolog.InfoLevel
	}
}
