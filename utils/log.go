package utils

import (
	"os"
	"runtime"
	"log/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
}

var Log Logger
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
    Log.Logger = slog.New(handler)
	Log.Logger = Log.Logger.With("pid", os.Getpid())
}

func (l Logger) Error(msg string, args ...any){
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		l.Logger.Error(msg, "file", file, "line", line, "caller", funcName, args)
	} else {
		l.Logger.Error(msg, args)
	}
	
}

func CleanUp() {
	if lumberjackLogger != nil {
		lumberjackLogger.Close()
	}
}
