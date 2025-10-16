package repository_airport

import (
	"context"
	"database/sql"
	"flight-api/internal/model"
)

type IAirportRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (*model.Airport, error)
	FindAll(ctx context.Context, db *sql.DB, args map[string]interface{}) ([]model.Airport, int, error)
	FindBySearchName(ctx context.Context, db *sql.DB, name string, args map[string]interface{}) ([]model.Airport, int, error)
	FindByID(ctx context.Context, db *sql.DB, id string) (model.Airport, error)
	FindExistsByICAOID(ctx context.Context, db *sql.DB, icaoId string) (*bool, error)
	FindByICAOID(ctx context.Context, db *sql.DB, icaoId string) (model.Airport, error)
	Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}
