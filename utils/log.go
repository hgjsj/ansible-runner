package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger
var lumberjackLogger *lumberjack.Logger

// InitLogger initialize a log instance
func InitLogger(name string) {

	// Configure lumberjack for log rotation
	lumberjackLogger = new(lumberjack.Logger)

	lumberjackLogger.Filename = "logs/" + name
	lumberjackLogger.MaxSize = 10   // megabytes
	lumberjackLogger.MaxBackups = 3 // number of old log files to retain

	// Configure the file write syncer
	writeSyncer := zapcore.AddSync(lumberjackLogger)

	// Configure the encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Configure the core
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	// Create the logger
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer Log.Sync() // flushes buffer, if any
}

func CleanUp() {
	if lumberjackLogger != nil {
		lumberjackLogger.Close()
	}
}
