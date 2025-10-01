package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

const INFO_DEBUG_LEVEL string = "info"
const DEBUG_LEVEL string = "debug"
const TRACE_LEVEL string = "trace"

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

func (l *Logger) SetLogLevel(logLevel string) {
	// Parse log level from string
	switch logLevel {
	case "trace":
		l.logger.SetLevel(logrus.TraceLevel)
	case "debug":
		l.logger.SetLevel(logrus.DebugLevel)
	case "info":
		l.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		l.logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		l.logger.SetLevel(logrus.FatalLevel)
	case "panic":
		l.logger.SetLevel(logrus.PanicLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Tracef(msg string, args ...interface{}) {
	l.logger.Tracef(msg, args...)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l *Logger) Panicf(msg string, args ...interface{}) {
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
