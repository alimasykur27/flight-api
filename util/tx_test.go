package util_test

import (
	"flight-api/util"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCommitOrRollback(t *testing.T) {
	t.Run("valid - commit when no panic", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error creating sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectCommit()

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unexpected error begin: %v", err)
		}

		func() {
			defer util.CommitOrRollback(tx)
			// no panic here → should commit
		}()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("invalid - rollback when panic", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error creating sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectRollback()

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unexpected error begin: %v", err)
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		}()

		func() {
			defer util.CommitOrRollback(tx)
			panic("boom!") // force rollback
		}()
	})
}
