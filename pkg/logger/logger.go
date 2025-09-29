package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

const INFO_DEBUG_LEVEL string = "info"

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(logLevel string) *Logger {
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
	return &Logger{
		logger: &logrus.Logger{
			Out:          output,
			Formatter:    formatter,
			Hooks:        make(logrus.LevelHooks),
			Level:        level,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
	}
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Tracef(msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l *Logger) Panic(msg string, args ...interface{}) {
	l.logger.Panicf(msg, args...)
}

// With returns a logger with the specified fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) Tracew(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Tracef(msg, args...)
}

func (l *Logger) Debugw(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Debugf(msg, args...)
}

func (l *Logger) Infow(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Infof(msg, args...)
}

func (l *Logger) Warnw(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Warnf(msg, args...)
}

func (l *Logger) Errorw(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Errorf(msg, args...)
}

func (l *Logger) Fatalw(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Fatalf(msg, args...)
}

func (l *Logger) Panicw(fields logrus.Fields, msg string, args ...interface{}) {
	l.logger.WithFields(fields).Panicf(msg, args...)
}
