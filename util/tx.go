package util

import (
	"database/sql"
	"fmt"
)

func FinishTx(tx *sql.Tx, err *error) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		if err != nil && *err == nil {
			*err = fmt.Errorf("panic recovered in tx: %v", r)
		}
		panic(r)
	}

	if err == nil {
		_ = tx.Rollback()
		return
	}

	if *err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			*err = fmt.Errorf("op failed: %w; rollback failed: %v", *err, rbErr)
		}
		return
	}

	if cmErr := tx.Commit(); cmErr != nil {
		*err = cmErr
	}
}
