package util_test

import (
	"errors"
	"flight-api/util"
	"testing"
)

// --- LogPanicError ---

func TestLogPanicError_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic, got: %v", r)
		}
	}()
	util.LogPanicError(errors.New("should just log, not panic"))
}

// --- PanicIfError ---

func TestPanicIfError_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	util.PanicIfError(errors.New("boom"))
}

func TestPanicIfError_NoPanicOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic, got: %v", r)
		}
	}()
	util.PanicIfError(nil)
}

// --- RecoverPanic ---

func TestRecoverPanic_CapturesString(t *testing.T) {
	var err error
	func() {
		defer util.RecoverPanic(&err)
		panic("string panic")
	}()

	if err == nil || err.Error() != "string panic" {
		t.Fatalf("expected error 'string panic', got: %v", err)
	}
}

func TestRecoverPanic_CapturesError(t *testing.T) {
	var err error

	someErr := errors.New("wrapped")

	func() {
		defer util.RecoverPanic(&err)
		panic(someErr)
	}()

	if err == nil || err.Error() != "wrapped" {
		t.Fatalf("expected error 'wrapped', got: %v", err)
	}
}

func TestRecoverPanic_CapturesUnknown(t *testing.T) {
	var err error

	func() {
		defer util.RecoverPanic(&err)
		panic(25)
	}()

	if err == nil || err.Error() != "unknown panic" {
		t.Fatalf("expected error 'unknown panic', got: %v", err)
	}
}
