package utils

import (
	"os"
	"log/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *slog.Logger
var lumberjackLogger *lumberjack.Logger

// InitLogger initialize a log instance
func InitLogger(name string) {

	// Configure lumberjack for log rotation
	lumberjackLogger = new(lumberjack.Logger)

	lumberjackLogger.Filename = name
	lumberjackLogger.MaxSize = 10   // megabytes
	lumberjackLogger.MaxBackups = 3 // number of old log files to retain

    // Create a new slog logger with lumberjack as the output
    handler := slog.NewJSONHandler(lumberjackLogger, nil)
    Log = slog.New(handler)
	Log = Log.With("pid", os.Getpid())
}

func CleanUp() {
	if lumberjackLogger != nil {
		lumberjackLogger.Close()
	}
}
