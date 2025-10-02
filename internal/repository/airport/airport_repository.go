package repository_airport

import (
	"context"
	"database/sql"
	"flight-api/internal/model"
)

type IAirportRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error)
	SyncAirport(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error)
	FindAll(ctx context.Context, tx *sql.Tx, args map[string]interface{}) ([]model.Airport, int, error)
	FindBySearchName(ctx context.Context, tx *sql.Tx, name string, args map[string]interface{}) ([]model.Airport, int, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string) (model.Airport, error)
	FindExistsByICAOID(ctx context.Context, tx *sql.Tx, icaoId string) (bool, error)
	FindByICAOID(ctx context.Context, tx *sql.Tx, icaoId string) (model.Airport, error)
	Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}
