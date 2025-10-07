// logger/logger_test.go
package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	test "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerDevUsesTextFormatter(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("info")
	ll := l.GetLogger()

	// level
	assert.Equal(t, logrus.InfoLevel, ll.GetLevel())

	// formatter
	_, isText := ll.Formatter.(*logrus.TextFormatter)
	assert.True(t, isText, "development should use TextFormatter")

	// output (stderr)
	assert.Equal(t, os.Stderr, ll.Out)
}

func TestNewLoggerProdUsesJSONFormatter(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	l := NewLogger("info")
	ll := l.GetLogger()

	assert.Equal(t, logrus.InfoLevel, ll.GetLevel())

	_, isJSON := ll.Formatter.(*logrus.JSONFormatter)
	assert.True(t, isJSON, "non-development should use JSONFormatter")
	assert.Equal(t, os.Stderr, ll.Out)
}

func TestNewLoggerSetLevelTrace(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("trace")
	assert.Equal(t, logrus.TraceLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelDebug(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("debug")
	assert.Equal(t, logrus.DebugLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelInfo(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("info")
	assert.Equal(t, logrus.InfoLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelWarn(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("warn")
	assert.Equal(t, logrus.WarnLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelError(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("error")
	assert.Equal(t, logrus.ErrorLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelFatal(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("fatal")
	assert.Equal(t, logrus.FatalLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerSetLevelPanic(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("panic")
	assert.Equal(t, logrus.PanicLevel, l.GetLogger().GetLevel())
}

func TestNewLoggerDefaultLevelOnUnknownString(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("unknown-level")
	assert.Equal(t, logrus.InfoLevel, l.GetLogger().GetLevel())
}

func TestSetLogLevelChangesLevel(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("info")
	l.SetLogLevel("trace")
	assert.Equal(t, logrus.TraceLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("debug")
	assert.Equal(t, logrus.DebugLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("info")
	assert.Equal(t, logrus.InfoLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("warn")
	assert.Equal(t, logrus.WarnLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("error")
	assert.Equal(t, logrus.ErrorLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("fatal")
	assert.Equal(t, logrus.FatalLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("panic")
	assert.Equal(t, logrus.PanicLevel, l.GetLogger().GetLevel())
	l.SetLogLevel("unknown-level")
	assert.Equal(t, logrus.InfoLevel, l.GetLogger().GetLevel())
}

func TestPassthroughMethodsWriteEntries(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	// NewLogger should set level based on "trace" string.
	l := NewLogger("trace")
	hook := test.NewLocal(l.GetLogger())

	// Trace
	l.Trace("hello")
	l.Tracef("world: %d", 42)
	l.Tracew(logrus.Fields{"a": 1}, "trace: %s", "x")
	l.WithFields(logrus.Fields{"k": "v"}).Trace("with fields")

	// Debug
	l.Debug("hello")
	l.Debugf("world: %d", 42)
	l.Debugw(logrus.Fields{"a": 1}, "debug: %s", "x")
	l.WithFields(logrus.Fields{"k": "v"}).Debug("with fields")

	// Info
	l.Info("hello")
	l.Infof("world: %d", 42)
	l.Infow(logrus.Fields{"a": 1}, "info: %s", "x")
	l.WithFields(logrus.Fields{"k": "v"}).Info("with fields")

	// Warn
	l.Warn("hello")
	l.Warnf("world: %d", 42)
	l.Warnw(logrus.Fields{"a": 1}, "warn: %s", "x")
	l.WithFields(logrus.Fields{"k": "v"}).Warn("with fields")

	// Error
	l.Error("hello")
	l.Errorf("world: %d", 42)
	l.Errorw(logrus.Fields{"a": 1}, "error: %s", "x")
	l.WithFields(logrus.Fields{"k": "v"}).Error("with fields")

	entries := hook.AllEntries()
	assert.Len(t, entries, 20, "should write 20 log entries (5 levels × 4 calls)")

	type check struct {
		idx    int
		level  logrus.Level
		msg    string
		fields map[string]any
	}

	expect := []check{
		// Trace (0..3)
		{0, logrus.TraceLevel, "hello", nil},
		{1, logrus.TraceLevel, "world: 42", nil},
		{2, logrus.TraceLevel, "trace: x", map[string]any{"a": 1}},
		{3, logrus.TraceLevel, "with fields", map[string]any{"k": "v"}},

		// Debug (4..7)
		{4, logrus.DebugLevel, "hello", nil},
		{5, logrus.DebugLevel, "world: 42", nil},
		{6, logrus.DebugLevel, "debug: x", map[string]any{"a": 1}},
		{7, logrus.DebugLevel, "with fields", map[string]any{"k": "v"}},

		// Info (8..11)
		{8, logrus.InfoLevel, "hello", nil},
		{9, logrus.InfoLevel, "world: 42", nil},
		{10, logrus.InfoLevel, "info: x", map[string]any{"a": 1}},
		{11, logrus.InfoLevel, "with fields", map[string]any{"k": "v"}},

		// Warn (12..15)
		{12, logrus.WarnLevel, "hello", nil},
		{13, logrus.WarnLevel, "world: 42", nil},
		{14, logrus.WarnLevel, "warn: x", map[string]any{"a": 1}},
		{15, logrus.WarnLevel, "with fields", map[string]any{"k": "v"}},

		// Error (16..19)
		{16, logrus.ErrorLevel, "hello", nil},
		{17, logrus.ErrorLevel, "world: 42", nil},
		{18, logrus.ErrorLevel, "error: x", map[string]any{"a": 1}},
		{19, logrus.ErrorLevel, "with fields", map[string]any{"k": "v"}},
	}

	for _, e := range expect {
		entry := entries[e.idx]
		assert.Equal(t, e.level, entry.Level, "level mismatch at index %d", e.idx)
		assert.Equal(t, e.msg, entry.Message, "message mismatch at index %d", e.idx)

		if e.fields != nil {
			for k, v := range e.fields {
				assert.Contains(t, entry.Data, k, "missing field %q at index %d", k, e.idx)
				assert.Equal(t, v, entry.Data[k], "field %q mismatch at index %d", k, e.idx)
			}
		} else {
			// when not expecting fields, ensure none of our test fields leaked in
			assert.NotContains(t, entry.Data, "a", "unexpected field at index %d", e.idx)
			assert.NotContains(t, entry.Data, "k", "unexpected field at index %d", e.idx)
		}
	}
}

func TestFatalDoesNotExitWhenExitFuncStubbed(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	l := NewLogger("fatal")
	ll := l.GetLogger()

	exitCount := 0
	ll.ExitFunc = func(code int) {
		exitCount++
		assert.Equal(t, 1, code, "Fatal* should exit with code 1")
	}

	hook := test.NewLocal(ll)

	l.Fatal("ini fatal")
	l.Fatalf("fatal msg: %s", "boom")
	l.Fatalw(logrus.Fields{"k": "v"}, "fatal")
	l.WithFields(logrus.Fields{"k": "v"}).Fatal("with fields")

	entries := hook.AllEntries()

	// ExitFunc should be called once per Fatal* call
	assert.Equal(t, 4, exitCount, "ExitFunc should be called 4 times")

	// We should have 4 fatal entries
	assert.Len(t, entries, 4)
	for i, e := range entries {
		assert.Equal(t, logrus.FatalLevel, e.Level, "entry %d not Fatal level", i)
	}

	// Spot-check messages and fields
	assert.Equal(t, "ini fatal", entries[0].Message)
	assert.Equal(t, "fatal msg: boom", entries[1].Message)

	assert.Equal(t, "fatal", entries[2].Message)
	assert.Equal(t, "v", entries[2].Data["k"])

	assert.Equal(t, "with fields", entries[3].Message)
	assert.Equal(t, "v", entries[3].Data["k"])
}

func TestPanicDoesNotExitWhenExitFuncStubbed(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	l := NewLogger("panic")
	ll := l.GetLogger()

	called := false
	ll.ExitFunc = func(int) { called = true } // stub os.Exit (used by Fatal*, not Panic*)

	hook := test.NewLocal(ll) // capture entries

	// Each Panic* must be wrapped so the test continues after the panic.
	assert.Panics(t, func() { l.Panic("ini panic") })
	assert.Panics(t, func() { l.Panicf("panic msg: %s", "boom") })
	assert.Panics(t, func() { l.Panicw(logrus.Fields{"k": "v"}, "panic") })
	assert.Panics(t, func() { l.WithFields(logrus.Fields{"k": "v"}).Panic("with fields") })

	entries := hook.AllEntries()

	// Panic* should NOT call ExitFunc; only Fatal* does.
	assert.False(t, called, "ExitFunc should NOT be called for Panic*")

	// We logged 4 entries, all at Panic level
	assert.Len(t, entries, 4)
	for i, e := range entries {
		assert.Equal(t, logrus.PanicLevel, e.Level, "entry %d not Panic level", i)
	}

	// Quick spot-check of messages/fields
	assert.Equal(t, "ini panic", entries[0].Message)
	assert.Equal(t, "panic msg: boom", entries[1].Message)
	assert.Equal(t, "panic", entries[2].Message)
	assert.Equal(t, "v", entries[2].Data["k"])
	assert.Equal(t, "with fields", entries[3].Message)
	assert.Equal(t, "v", entries[3].Data["k"])
}
