package utils

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
}

var Log Logger
var lumberjackLogger *lumberjack.Logger

// InitLogger initialize a log instance
func InitLogger() {
	logFile := viper.GetString("log.file")

	if logFile != "" {
		// Configure lumberjack for log rotation
		lumberjackLogger = new(lumberjack.Logger)

		lumberjackLogger.Filename = logFile
		lumberjackLogger.MaxSize = 10   // megabytes
		lumberjackLogger.MaxBackups = 3 // number of old log files to retain
		handler := slog.NewJSONHandler(lumberjackLogger, nil) // Create a new slog logger with lumberjack as the output
		Log.Logger = slog.New(handler)
	} else {
		// Default to stdout
		Log.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
    
    
    
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
