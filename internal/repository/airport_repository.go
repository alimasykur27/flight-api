package repository

import (
	"context"
	"database/sql"
	"flight-api/internal/model"
)

type IAirportRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error)
	FindAll(ctx context.Context, tx *sql.Tx, args ...interface{}) ([]model.Airport, int, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string) (model.Airport, error)
	Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}
