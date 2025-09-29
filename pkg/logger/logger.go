package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

const INFO_DEBUG_LEVEL string = "info"

func init() {
	logger = NewLogger(INFO_DEBUG_LEVEL)
}

func GetLogger() *logrus.Logger {
	return logger
}

func NewLogger(logLevel string) *logrus.Logger {
	var level logrus.Level

	// Parse log level from string
	switch logLevel {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel
	}

	// For development, use console formatter and write to stderr
	// For production, use JSON formatter
	var formatter logrus.Formatter
	var output io.Writer
	if os.Getenv("APP_ENV") != "development" {
		formatter = &logrus.JSONFormatter{}
		output = os.Stderr
	} else {
		formatter = &logrus.TextFormatter{}
		output = os.Stderr
	}

	// Create logger
	logger := &logrus.Logger{
		Out:          output,
		Formatter:    formatter,
		Hooks:        make(logrus.LevelHooks),
		Level:        level,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	return logger
}

func Trace(msg string, args ...interface{}) {
	logger.Tracef(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

// With returns a logger with the specified fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func Tracew(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Tracef(msg, args...)
}

func Debugw(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Debugf(msg, args...)
}

func Infow(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Infof(msg, args...)
}

func Warnw(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Warnf(msg, args...)
}

func Errorw(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Errorf(msg, args...)
}

func Fatalw(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Fatalf(msg, args...)
}

func Panicw(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Panicf(msg, args...)
}
