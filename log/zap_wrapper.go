package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const FILENAME = "server.log"

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {}

	config.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)

	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugaredLogger = logger.Sugar()
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}

func Infof(format string, args ...any) {
	sugaredLogger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	sugaredLogger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	sugaredLogger.Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	sugaredLogger.Panicf(format, args...)
}

func Fatalf(format string, args ...any) {
	sugaredLogger.Fatalf(format, args...)
}

func Sync() {
	logger.Sync()
}

func Named(name string) *zap.Logger {
	return logger.Named(name)
}
