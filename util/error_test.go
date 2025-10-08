package util

import (
	"errors"
	"testing"
)

// --- LogPanicError ---

func TestLogPanicError_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic, got: %v", r)
		}
	}()
	LogPanicError(errors.New("should just log, not panic"))
}

// --- PanicIfError ---

func TestPanicIfError_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	PanicIfError(errors.New("boom"))
}

func TestPanicIfError_NoPanicOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic, got: %v", r)
		}
	}()
	PanicIfError(nil)
}
