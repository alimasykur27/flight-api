package repository

import (
	"context"
	"database/sql"
	"flight-api/internal/model"
	"flight-api/pkg/logger"
	"flight-api/util"

	"github.com/google/uuid"
)

type AirportRepository struct {
	logger *logger.Logger
}

func NewAirportRepository(l *logger.Logger) IAirportRepository {
	return &AirportRepository{
		logger: l,
	}
}

func (r *AirportRepository) Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error) {
	SQL := `
		INSERT INTO airports (
			site_number, icao_id, faa_id, iata_id, name, type, status,
			country, state, state_full, county, city,
			ownership, "use", manager, manager_phone,
			latitude, latitude_sec, longitude, longitude_sec,
			elevation, control_tower, unicom, ctaf
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24
		)
		RETURNING id`

	row := tx.QueryRowContext(
		ctx,
		SQL,
		airport.SiteNumber,
		airport.ICAOID,
		airport.FAAID,
		airport.IATAID,
		airport.Name,
		airport.Type,
		airport.Status,
		airport.Country,
		airport.State,
		airport.StateFull,
		airport.County,
		airport.City,
		airport.Ownership,
		airport.Use,
		airport.Manager,
		airport.ManagerPhone,
		airport.Latitude,
		airport.LatitudeSec,
		airport.Longitude,
		airport.LongitudeSec,
		airport.Elevation,
		airport.ControlTower,
		airport.Unicom,
		airport.CTAF,
	)

	var id string
	row.Scan(&id)

	result, err := r.FindByID(ctx, tx, id)
	util.PanicIfError(err)

	return result, nil
}

func (r *AirportRepository) FindByID(ctx context.Context, tx *sql.Tx, id string) (model.Airport, error) {
	SQL := `
SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status,
	country, state, state_full, county, city, ownership, "use",
	manager, manager_phone, latitude, latitude_sec, longitude, longitude_sec, elevation,
	control_tower, unicom, ctaf, created_at, updated_at
FROM airports 
WHERE id = $1 
LIMIT 1
`
	airportId, err := uuid.Parse(id)
	if err != nil {
		return model.Airport{}, util.ErrNotFound
	}

	rows, err := tx.QueryContext(ctx, SQL, airportId)
	util.PanicIfError(err)
	defer rows.Close()

	airport := model.Airport{}
	if rows.Next() {
		err := rows.Scan(
			&airport.ID,
			&airport.SiteNumber,
			&airport.ICAOID,
			&airport.FAAID,
			&airport.IATAID,
			&airport.Name,
			&airport.Type,
			&airport.Status,
			&airport.Country,
			&airport.State,
			&airport.StateFull,
			&airport.City,
			&airport.County,
			&airport.Ownership,
			&airport.Use,
			&airport.Manager,
			&airport.ManagerPhone,
			&airport.Latitude,
			&airport.LatitudeSec,
			&airport.Longitude,
			&airport.LongitudeSec,
			&airport.Elevation,
			&airport.ControlTower,
			&airport.Unicom,
			&airport.CTAF,
			&airport.CreatedAt,
			&airport.UpdatedAt,
		)
		util.PanicIfError(err)
		return airport, nil
	} else {
		return model.Airport{}, util.ErrNotFound
	}
}

func (r *AirportRepository) FindAll(ctx context.Context, tx *sql.Tx, args ...interface{}) ([]model.Airport, int, error) {
	SQL := `SELECT id, site_number, icao_id, faa_id, iata_id, type, status, created_at, updated_at
			FROM airports 
			ORDER BY icao_id
			LIMIT $1
			OFFSET $2`

	limit, offset, err := util.ParsePagination(args...)
	util.PanicIfError(err)

	var total int
	row := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM airports`)
	err = row.Scan(&total)
	util.PanicIfError(err)

	rows, err := tx.QueryContext(ctx, SQL, limit, offset)
	util.PanicIfError(err)
	defer rows.Close()

	var airports []model.Airport
	for rows.Next() {
		airport := model.Airport{}
		err := rows.Scan(
			&airport.ID,
			&airport.SiteNumber,
			&airport.ICAOID,
			&airport.FAAID,
			&airport.Name,
			&airport.Type,
			&airport.Status,
			&airport.CreatedAt,
			&airport.UpdatedAt,
		)
		util.PanicIfError(err)
		airports = append(airports, airport)
	}

	return airports, total, nil
}

func (r *AirportRepository) Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error) {
	SQL := `
		UPDATE airports SET
			site_number = $1,
			faa_id = $2,
			iata_id = $3,
			name = $4,
			type = $5,
			status = $6,
			country = $7,
			state = $8,
			state_full = $9,
			county = $10,
			city = $11,
			ownership = $12,
			"use" = $13,
			manager = $14,
			manager_phone = $15,
			latitude = $16,
			latitude_sec = $17,
			longitude = $18,
			longitude_sec = $19,
			elevation = $20,
			control_tower = $21,
			unicom = $22,
			ctaf = $23,
			effective_date = $24,
			updated_at = NOW()
		WHERE id = $25
		RETURNING id
	`

	airportId, err := uuid.Parse(id)
	if err != nil {
		return model.Airport{}, util.ErrNotFound
	}

	row := tx.QueryRowContext(
		ctx,
		SQL,
		airport.SiteNumber,
		airport.FAAID,
		airport.IATAID,
		airport.Name,
		airport.Type,
		airport.Status,
		airport.Country,
		airport.State,
		airport.StateFull,
		airport.County,
		airport.City,
		airport.Ownership,
		airport.Use,
		airport.Manager,
		airport.ManagerPhone,
		airport.Latitude,
		airport.LatitudeSec,
		airport.Longitude,
		airport.LongitudeSec,
		airport.Elevation,
		airport.ControlTower,
		airport.Unicom,
		airport.CTAF,
		airport.EffectiveDate,
		airportId,
	)

	var updatedID string
	err = row.Scan(&updatedID)
	if err == sql.ErrNoRows {
		return model.Airport{}, util.ErrNotFound
	}
	util.PanicIfError(err)

	updatedAirport, err := r.FindByID(ctx, tx, updatedID)
	util.PanicIfError(err)

	return updatedAirport, nil
}

func (r *AirportRepository) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := `DELETE FROM airports WHERE id = $1`

	airportId, err := uuid.Parse(id)
	if err != nil {
		return util.ErrNotFound
	}

	result, err := tx.ExecContext(ctx, SQL, airportId)
	util.PanicIfError(err)

	rowsAffected, err := result.RowsAffected()
	util.PanicIfError(err)

	if rowsAffected == 0 {
		return util.ErrNotFound
	}

	return nil
}
