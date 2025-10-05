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
	l := NewLogger("debug")
	ll := l.GetLogger()

	// level
	assert.Equal(t, logrus.DebugLevel, ll.GetLevel())

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

	l.SetLogLevel("warn")
	assert.Equal(t, logrus.WarnLevel, l.GetLogger().GetLevel())
}

func TestPassthroughMethodsWriteEntries(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("debug")
	hook := test.NewLocal(l.GetLogger())

	l.Debug("hello")
	l.Infof("world: %d", 42)
	l.WithFields(logrus.Fields{"k": "v"}).Info("with fields")
	l.Warnw(logrus.Fields{"a": 1}, "warn: %s", "x")

	// minimal sanity checks
	assert.GreaterOrEqual(t, len(hook.AllEntries()), 3)
	levels := []logrus.Level{hook.AllEntries()[0].Level, hook.AllEntries()[1].Level}
	assert.Contains(t, levels, logrus.DebugLevel)
	assert.Contains(t, levels, logrus.InfoLevel)
}

func TestFatalDoesNotExitWhenExitFuncStubbed(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	l := NewLogger("fatal")
	ll := l.GetLogger()

	called := false
	ll.ExitFunc = func(int) { called = true } // stub os.Exit

	hook := test.NewLocal(ll) //nolint:staticcheck
	l.Fatalf("fatal msg: %s", "boom")
	entries := hook.AllEntries()
	assert.True(t, called, "ExitFunc should be called")
	assert.GreaterOrEqual(t, len(entries), 1)
	assert.Equal(t, logrus.FatalLevel, entries[len(entries)-1].Level)
}
