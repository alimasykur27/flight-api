package repository_airport

import (
	"context"
	"database/sql"
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"flight-api/pkg/logger"
	"flight-api/util"
	"strings"

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

func (r *AirportRepository) Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (*model.Airport, error) {
	SQL := `
		INSERT INTO airports (
			site_number, icao_id, faa_id, iata_id, name, 
			type, status, country, state, state_full, 
			county, city, ownership, "use", manager, 
			manager_phone, latitude, latitude_sec, longitude, longitude_sec,
			elevation, control_tower, unicom, ctaf, effective_date,
			sync_status, sync_message
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, 
			$11, $12, $13, $14, $15, 
			$16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25,
			$26, $27
		) 
		RETURNING 
			id,
			site_number, icao_id, faa_id, iata_id, name, 
			type, status, country, state, state_full, 
			county, city, ownership, "use", manager, 
			manager_phone, latitude, latitude_sec, longitude, longitude_sec,
			elevation, control_tower, unicom, ctaf, effective_date,
			sync_status, sync_message, created_at, updated_at
	`

	row := tx.QueryRowContext(
		ctx,
		strings.TrimSpace(SQL),
		airport.SiteNumber, airport.ICAOID, airport.FAAID, airport.IATAID, airport.Name,
		airport.Type, airport.Status, airport.Country, airport.State, airport.StateFull,
		airport.County, airport.City, airport.Ownership, airport.Use, airport.Manager,
		airport.ManagerPhone, airport.Latitude, airport.LatitudeSec, airport.Longitude, airport.LongitudeSec,
		airport.Elevation, airport.ControlTower, airport.Unicom, airport.CTAF, airport.EffectiveDate,
		enum.SYNC_SYNCED.Int(), enum.SYNC_SYNCED.String(),
	)

	var result model.Airport
	err := row.Scan(
		&result.ID,
		&result.SiteNumber, &result.ICAOID, &result.FAAID, &result.IATAID, &result.Name,
		&result.Type, &result.Status, &result.Country, &result.State, &result.StateFull,
		&result.County, &result.City, &result.Ownership, &result.Use, &result.Manager,
		&result.ManagerPhone, &result.Latitude, &result.LatitudeSec, &result.Longitude, &result.LongitudeSec,
		&result.Elevation, &result.ControlTower, &result.Unicom, &result.CTAF, &result.EffectiveDate,
		&result.SyncStatus, &result.SyncMessage, &result.UpdatedAt, &result.CreatedAt,
	)

	if err != nil {
		r.logger.Errorf("[Insert] Failed to insert airport: %v", err)
		return nil, util.NewErrorException(util.ErrDatabase, "failed to insert airport: "+err.Error())
	}

	r.logger.Debugf("[Insert] Inserted airport with ID: %s", result.ID.String())
	return &result, nil
}

func (r *AirportRepository) FindByID(ctx context.Context, db *sql.DB, id string) (model.Airport, error) {
	SQL := `
		SELECT 
			id, 
			site_number, icao_id, faa_id, iata_id, name, 
			type, status, country, state, state_full, 
			county, city, ownership, "use", manager, 
			manager_phone, latitude, latitude_sec, longitude, longitude_sec, 
			elevation, control_tower, unicom, ctaf, effective_date, 
			sync_status, sync_message, updated_at, created_at
		FROM airports 
		WHERE id = $1 
		LIMIT 1
	`

	airportId, err := uuid.Parse(id)
	if err != nil {
		return model.Airport{}, util.NewErrorException(util.ErrInvalidInput, "invalid airport ID format: "+err.Error())
	}

	rows, err := db.QueryContext(ctx, strings.TrimSpace(SQL), airportId)
	if err != nil {
		r.logger.Errorf("failed to execute query select airport by ID: %v", err)
		return model.Airport{}, util.NewErrorException(util.ErrDatabase, "failed to execute query select airport by ID: "+err.Error())
	}
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
			&airport.County,
			&airport.City,
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
			&airport.EffectiveDate,
			&airport.SyncStatus,
			&airport.SyncMessage,
			&airport.CreatedAt,
			&airport.UpdatedAt,
		)

		if err != nil {
			r.logger.Errorf("failed to scan airport row: %v", err)
			return model.Airport{}, err
		}

		return airport, nil
	}

	r.logger.Warn("[FindByID] No airport found with the given ID")
	return model.Airport{}, util.ErrNotFound
}

func (r *AirportRepository) FindAll(ctx context.Context, db *sql.DB, args map[string]interface{}) ([]model.Airport, int, error) {
	limit, offset := util.ParsePagination(args)

	// Get total count
	var total int
	TotalSQL := `SELECT COUNT(*) FROM airports`
	row := db.QueryRowContext(ctx, TotalSQL)
	err := row.Scan(&total)
	if err != nil {
		r.logger.Errorf("[FindAll] failed to get total count of airports: %v", err)
		return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to execute query select count: "+err.Error())
	}

	// Find All Airports
	SQL := `
		SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status, created_at, updated_at
		FROM airports 
		ORDER BY icao_id
		LIMIT $1
		OFFSET $2
	`

	rows, err := db.QueryContext(ctx, strings.TrimSpace(SQL), limit, offset)
	if err != nil {
		r.logger.Errorf("[FindAll] failed to fetch airports: %v", err)
		return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to execute query select airports: "+err.Error())
	}

	defer rows.Close()

	var airports []model.Airport
	for rows.Next() {
		airport := model.Airport{}
		err := rows.Scan(
			&airport.ID,
			&airport.SiteNumber,
			&airport.ICAOID,
			&airport.FAAID,
			&airport.IATAID,
			&airport.Name,
			&airport.Type,
			&airport.Status,
			&airport.CreatedAt,
			&airport.UpdatedAt,
		)

		if err != nil {
			r.logger.Errorf("[FindAll] Failed to scan airport row: %v", err)
			return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to scan airport row: "+err.Error())
		}

		airports = append(airports, airport)
	}

	return airports, total, nil
}

func (r *AirportRepository) FindBySearchName(ctx context.Context, db *sql.DB, name string, args map[string]interface{}) ([]model.Airport, int, error) {
	r.logger.Debugf("[FindBySearchName] Find airports by name: %s", name)

	limit, offset := util.ParsePagination(args)
	searchName := "%" + name + "%"

	SQL := `
		SELECT 
			id, 
			site_number, icao_id, faa_id, iata_id, name, 
			type, status, country, state, state_full, 
			county, city, ownership, "use", manager, 
			manager_phone, latitude, latitude_sec, longitude, longitude_sec, 
			elevation, control_tower, unicom, ctaf, effective_date, 
			sync_status, sync_message, updated_at, created_at
		FROM airports 
		WHERE 
			LOWER(name) LIKE LOWER($3)
		ORDER BY icao_id
		LIMIT $1
		OFFSET $2
	`

	rows, err := db.QueryContext(ctx, strings.TrimSpace(SQL), limit, offset, searchName)
	if err != nil {
		r.logger.Error("[FindBySearchName] Failed to execute query select airport")
		return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to execute query select airports: "+err.Error())
	}
	defer rows.Close()

	var airports []model.Airport
	for rows.Next() {
		airport := model.Airport{}
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
			&airport.County,
			&airport.City,
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
			&airport.EffectiveDate,
			&airport.SyncStatus,
			&airport.SyncMessage,
			&airport.UpdatedAt,
			&airport.CreatedAt,
		)

		if err != nil {
			r.logger.Errorf("[FindBySearchName] Failed to scan airport row: %v", err)
			return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to scan airport row: "+err.Error())
		}

		airports = append(airports, airport)
	}

	var total int
	TotalSQL := `SELECT COUNT(*) FROM airports WHERE name ILIKE $1`

	row := db.QueryRowContext(ctx, TotalSQL, searchName)
	err = row.Scan(&total)
	if err != nil {
		r.logger.Error("[FindBySearchName] Failed to execute query select count airport")
		return nil, 0, util.NewErrorException(util.ErrDatabase, "failed to execute query select count airport: "+err.Error())
	}

	return airports, total, nil
}

func (r *AirportRepository) FindExistsByICAOID(ctx context.Context, db *sql.DB, icaoId string) (*bool, error) {
	r.logger.Debugf("[FindExistsByICAOID] check existing airport by ICAO ID")
	var exists int

	SQL := `SELECT 1 FROM airports WHERE icao_id = $1 LIMIT 1`
	row := db.QueryRowContext(ctx, strings.TrimSpace(SQL), icaoId)
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("[FindExistsByICAOID] Failed to check existing airport by ICAO ID: %v", err)
		return nil, util.NewErrorException(util.ErrDatabase, "failed to check existing airport by ICAO ID: "+err.Error())
	}

	r.logger.Debugf("[FindExistsByICAOID] Exists value: %d", exists)
	if err == sql.ErrNoRows || exists != 1 {
		r.logger.Debugf("[FindExistsByICAOID] Airport with ICAO ID %s does not exist", icaoId)
		return util.Ptr(false), util.NewErrorException(nil, "airport does not exist: "+icaoId)
	} else if err != nil {
		return nil, util.NewErrorException(util.ErrDatabase, "failed to check existing airport by ICAO ID: "+err.Error())
	}

	r.logger.Debugf("[FindExistsByICAOID] Airport with ICAO ID %s exists", icaoId)
	return util.Ptr(true), nil
}

func (r *AirportRepository) FindByICAOID(ctx context.Context, db *sql.DB, icaoId string) (model.Airport, error) {
	SQL := `
		SELECT 
			id, 
			site_number, icao_id, faa_id, iata_id, name, 
			type, status, country, state, state_full, 
			county, city, ownership, "use", manager, 
			manager_phone, latitude, latitude_sec, longitude, longitude_sec, 
			elevation, control_tower, unicom, ctaf, effective_date,
			sync_status, sync_message, updated_at, created_at
		FROM airports 
		WHERE icao_id = $1 
		LIMIT 1
	`

	rows, err := db.QueryContext(ctx, strings.TrimSpace(SQL), icaoId)
	if err == sql.ErrNoRows {
		r.logger.Error("[FindByICAOID] No airport found")
		return model.Airport{}, util.NewErrorException(err, "failed to execute select query: "+err.Error())
	}

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
			&airport.County,
			&airport.City,
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
			&airport.EffectiveDate,
			&airport.SyncStatus,
			&airport.SyncMessage,
			&airport.UpdatedAt,
			&airport.CreatedAt,
		)

		if err != nil {
			r.logger.Errorf("[FindByICAOID] Failed to scan airport row: %v", err)
			return model.Airport{}, util.NewErrorException(util.ErrDatabase, "failed to scan airport row: "+err.Error())
		}

		return airport, nil
	} else {
		return model.Airport{}, util.NewErrorException(util.ErrNotFound, "airport not found")
	}
}

func (r *AirportRepository) Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error) {
	SQL := `
		UPDATE airports 
		SET
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
		RETURNING 
			id,
			site_number, icao_id, faa_id, iata_id, name,
			type, status, country, state, state_full,
			county, city, ownership, "use", manager,
			manager_phone, latitude, latitude_sec, longitude, longitude_sec,
			elevation, control_tower, unicom, ctaf, effective_date,
			sync_status, sync_message, updated_at, created_at
	`

	airportId, err := uuid.Parse(id)
	if err != nil {
		return model.Airport{}, util.ErrBadRequest
	}

	row := tx.QueryRowContext(
		ctx,
		strings.TrimSpace(SQL),
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

	err = row.Scan(
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
		&airport.County,
		&airport.City,
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
		&airport.EffectiveDate,
		&airport.SyncStatus,
		&airport.SyncMessage,
		&airport.CreatedAt,
		&airport.UpdatedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("Failed to update airport: %v", err)
		return model.Airport{}, err
	} else if err == sql.ErrNoRows {
		return model.Airport{}, util.ErrNotFound
	}

	return airport, nil
}

func (r *AirportRepository) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := `DELETE FROM airports WHERE id = $1`

	airportId, err := uuid.Parse(id)
	if err != nil {
		return util.NewErrorException(util.ErrInvalidInput, "invalid airport ID format: "+err.Error())
	}

	result, err := tx.ExecContext(ctx, SQL, airportId)
	if err != nil {
		return util.NewErrorException(util.ErrDatabase, "failed to execute delete query: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return util.NewErrorException(util.ErrDatabase, "failed to get rows affected: "+err.Error())
	}

	if rowsAffected == 0 {
		return util.NewErrorException(util.ErrNotFound, "airport not found")
	}

	return nil
}
