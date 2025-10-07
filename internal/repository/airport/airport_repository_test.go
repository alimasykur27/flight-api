package repository_airport

import (
	"context"
	"database/sql"
	"database/sql/driver"
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"flight-api/pkg/logger"
	"flight-api/util"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var log = logger.NewLogger(logger.DEBUG_LEVEL)

// ---------- HELPER FUNCTIONS ----------
var selectAirportQuery string = `SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status,
	country, state, state_full, county, city, ownership, "use",
	manager, manager_phone, latitude, latitude_sec, longitude, longitude_sec, elevation,
	control_tower, unicom, ctaf, effective_date, created_at, updated_at
FROM airports 
WHERE id = $1 
LIMIT 1
`

func newCols() []string {
	return []string{
		"id", "site_number", "icao_id", "faa_id", "iata_id", "name",
		"type", "status", "country", "state", "state_full", "county", "city",
		"ownership", "use", "manager", "manager_phone",
		"latitude", "latitude_sec", "longitude", "longitude_sec",
		"elevation", "control_tower", "unicom", "ctaf",
		"effective_date", "created_at", "updated_at",
	}
}

func newMinCols() []string {
	return []string{
		"id",
		"site_number", "icao_id", "faa_id", "iata_id", "name",
		"type", "status",
		"created_at", "updated_at",
	}
}

func successRow(
	id *uuid.UUID,
	site, icao, faa, iata, name *string,
	typ enum.FasilityTypeEnum,
	status *bool,
	country, state, stateFull, county, city *string,
	ownership enum.OwnershipEnum,
	use enum.UseTypeEnum,
	manager, managerPhone, latitude, latitudeSec, longitude, longitudeSec *string,
	elevation *int64,
	controlTower *bool,
	unicom, ctaf *string,
	effectiveDate *time.Time,
	createdAt, updatedAt time.Time,
) *sqlmock.Rows {
	return sqlmock.NewRows(newCols()).AddRow(
		id.String(),
		site, icao, faa, iata, name,
		typ,
		status,
		country, state, stateFull, county, city,
		ownership,
		use,
		manager, managerPhone, latitude, latitudeSec, longitude, longitudeSec,
		elevation,
		controlTower,
		unicom, ctaf,
		effectiveDate,
		createdAt, updatedAt,
	)
}

func buildRowsFindAll(limit, offset int) *sqlmock.Rows {
	end := offset + limit
	if end > len(dataDummy) {
		end = len(dataDummy)
	}
	rs := sqlmock.NewRows(newMinCols())

	for _, data := range dataDummy[offset:end] {
		a := data.row
		rs.AddRow(
			data.id.String(),
			a.SiteNumber, a.ICAOID, a.FAAID, a.IATAID, a.Name,
			a.Type,
			a.Status,
			a.CreatedAt,
			a.UpdatedAt,
		)
	}

	return rs
}

func buildRowsByName(nameLike string, limit, offset int) (*sqlmock.Rows, int) {
	data := make([]model.Airport, 0)

	for i := range dataDummy {
		data = append(data, dataDummy[i].row)
	}

	pat := strings.ToLower(strings.Trim(nameLike, "%"))
	filtered := make([]model.Airport, 0)

	for _, a := range data {
		if strings.Contains(strings.ToLower(*a.Name), pat) {
			filtered = append(filtered, a)
		}
	}
	// ORDER BY icao_id
	sort.Slice(filtered, func(i, j int) bool { return *filtered[i].ICAOID < *filtered[j].ICAOID })
	total := len(filtered)

	end := offset + limit
	if end > total {
		end = total
	}
	var page []model.Airport
	if offset < total {
		page = filtered[offset:end]
	} else {
		page = []model.Airport{}
	}

	rows := sqlmock.NewRows(newCols())
	for _, a := range page {
		rows.AddRow(
			a.ID, a.SiteNumber, a.ICAOID, a.FAAID, a.IATAID, a.Name, a.Type, a.Status,
			a.Country, a.State, a.StateFull, a.County, a.City, a.Ownership, a.Use,
			a.Manager, a.ManagerPhone, a.Latitude, a.LatitudeSec, a.Longitude, a.LongitudeSec, a.Elevation,
			a.ControlTower, a.Unicom, a.CTAF, a.EffectiveDate, a.CreatedAt, a.UpdatedAt,
		)
	}
	return rows, total
}

// ------Data Dummy--------
var timeNow = time.Now()
var sliceId = map[string]uuid.UUID{
	"SJC": uuid.New(),
	"LAX": uuid.New(),
	"SFO": uuid.New(),
}

var dataDummy = []struct {
	id    uuid.UUID
	row   model.Airport
	label string
}{
	{
		id:    sliceId["SJC"],
		label: "SJC",
		row: func() model.Airport {
			id := sliceId["SJC"]
			return model.Airport{
				ID:            &id,
				SiteNumber:    util.Ptr("12345"),
				ICAOID:        util.Ptr("ICAO123"),
				FAAID:         util.Ptr("FAA9"),
				IATAID:        util.Ptr("IT9"),
				Name:          util.Ptr("Norman Y. Mineta"),
				Type:          enum.AIRPORT,
				Status:        util.Ptr(true),
				Country:       util.Ptr("US"),
				State:         util.Ptr("CA"),
				StateFull:     util.Ptr("California"),
				County:        util.Ptr("Santa Clara"),
				City:          util.Ptr("San Jose"),
				Ownership:     enum.OWN_PUBLIC,
				Use:           enum.USE_PUBLIC,
				Manager:       util.Ptr("Jane Doe"),
				ManagerPhone:  util.Ptr("+1-555-0100"),
				Latitude:      util.Ptr("37.3639"),
				LatitudeSec:   util.Ptr("21.8"),
				Longitude:     util.Ptr("-121.9289"),
				LongitudeSec:  nil,
				Elevation:     util.Ptr(int64(17)),
				ControlTower:  util.Ptr(true),
				Unicom:        util.Ptr("122.95"),
				CTAF:          util.Ptr("118.00"),
				EffectiveDate: nil,
				CreatedAt:     &timeNow,
				UpdatedAt:     &timeNow,
			}
		}(),
	},
	{
		id:    sliceId["LAX"],
		label: "LAX",
		row: func() model.Airport {
			id := sliceId["LAX"]
			return model.Airport{
				ID:            &id,
				SiteNumber:    util.Ptr("67890"),
				ICAOID:        util.Ptr("KLAX"),
				FAAID:         util.Ptr("FAA1"),
				IATAID:        util.Ptr("LAX"),
				Name:          util.Ptr("Los Angeles Intl"),
				Type:          enum.AIRPORT,
				Status:        util.Ptr(true),
				Country:       util.Ptr("US"),
				State:         util.Ptr("CA"),
				StateFull:     util.Ptr("California"),
				County:        util.Ptr("Los Angeles"),
				City:          util.Ptr("Los Angeles"),
				Ownership:     enum.OWN_PUBLIC,
				Use:           enum.USE_PUBLIC,
				Manager:       util.Ptr("John Smith"),
				ManagerPhone:  util.Ptr("+1-555-0200"),
				Latitude:      util.Ptr("33.9416"),
				LatitudeSec:   util.Ptr("00.0"),
				Longitude:     util.Ptr("-118.4085"),
				LongitudeSec:  nil,
				Elevation:     util.Ptr(int64(125)),
				ControlTower:  util.Ptr(true),
				Unicom:        util.Ptr("122.80"),
				CTAF:          util.Ptr("119.80"),
				EffectiveDate: nil,
				CreatedAt:     &timeNow,
				UpdatedAt:     &timeNow,
			}
		}(),
	},
	{
		id:    sliceId["SFO"],
		label: "SFO",
		row: func() model.Airport {
			id := sliceId["SFO"]
			return model.Airport{
				ID:            &id,
				SiteNumber:    util.Ptr("54321"),
				ICAOID:        util.Ptr("KSFO"),
				FAAID:         util.Ptr("FAA4"),
				IATAID:        util.Ptr("SFO"),
				Name:          util.Ptr("San Francisco Intl"),
				Type:          enum.AIRPORT,
				Status:        util.Ptr(true),
				Country:       util.Ptr("US"),
				State:         util.Ptr("CA"),
				StateFull:     util.Ptr("California"),
				County:        util.Ptr("San Mateo"),
				City:          util.Ptr("San Francisco"),
				Ownership:     enum.OWN_PUBLIC,
				Use:           enum.USE_PUBLIC,
				Manager:       util.Ptr("Alice Johnson"),
				ManagerPhone:  util.Ptr("+1-555-0300"),
				Latitude:      util.Ptr("37.7749"),
				LatitudeSec:   util.Ptr("49.0"),
				Longitude:     util.Ptr("-122.4194"),
				LongitudeSec:  util.Ptr("25.0"),
				Elevation:     util.Ptr(int64(13)),
				ControlTower:  util.Ptr(true),
				Unicom:        util.Ptr("123.00"),
				CTAF:          util.Ptr("121.50"),
				EffectiveDate: nil,
				CreatedAt:     &timeNow,
				UpdatedAt:     &timeNow,
			}
		}(),
	},
}

// ---------- UNIT TESTS ----------
func TestNewAirportRepository(t *testing.T) {
	res := NewAirportRepository(log)
	assert.NotNil(t, res)
	assert.IsType(t, &AirportRepository{}, res)
}

func TestAirportRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := &AirportRepository{logger: log}

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	defer util.CommitOrRollback(tx)

	// ---------- arrange input ----------
	req := airport_dto.AirportRequestDto{
		SiteNumber:    util.Ptr("12345"),
		ICAOID:        util.Ptr("KJFK"),
		FAAID:         util.Ptr("JFK"),
		IATAID:        util.Ptr("JFK"),
		Name:          util.Ptr("John F. Kennedy International Airport"),
		Type:          enum.AIRPORT,
		Status:        util.Ptr(true),
		Country:       util.Ptr("USA"),
		State:         util.Ptr("NY"),
		StateFull:     util.Ptr("New York"),
		County:        util.Ptr("Queens"),
		City:          util.Ptr("New York"),
		Ownership:     enum.OWN_PUBLIC,
		Use:           enum.USE_PUBLIC,
		Manager:       util.Ptr("Jane Doe"),
		ManagerPhone:  util.Ptr("+1-555-1234"),
		Latitude:      util.Ptr("40.6413 N"),
		LatitudeSec:   util.Ptr("38.0"),
		Longitude:     util.Ptr("73.7781 W"),
		LongitudeSec:  nil, // keep it nil to test
		Elevation:     util.Ptr(int64(13)),
		ControlTower:  util.Ptr(true),
		Unicom:        util.Ptr("123.45"),
		CTAF:          util.Ptr("123.45"),
		EffectiveDate: nil,
	}
	modelInput := airport_dto.AirportRequestToAirport(req)

	// ---------- expect INSERT ----------
	insertRe := regexp.MustCompile(`(?s)INSERT\s+INTO\s+airports\s*\(.*?\)\s*VALUES\s*\(.*?\)\s*RETURNING\s+id`)
	newID := uuid.New().String()

	mock.ExpectQuery(insertRe.String()).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newID))

	// ---------- expect SELECT (FindByID) ----------
	now := time.Now()
	rows := sqlmock.NewRows(newCols()).AddRow(
		newID,
		modelInput.SiteNumber,
		modelInput.ICAOID,
		modelInput.FAAID,
		modelInput.IATAID,
		modelInput.Name,
		modelInput.Type,
		modelInput.Status,
		modelInput.Country,
		modelInput.State,
		modelInput.StateFull,
		modelInput.County,
		modelInput.City,
		modelInput.Ownership,
		modelInput.Use,
		modelInput.Manager,
		modelInput.ManagerPhone,
		modelInput.Latitude,
		modelInput.LatitudeSec,
		modelInput.Longitude,
		modelInput.LongitudeSec,
		modelInput.Elevation,
		modelInput.ControlTower,
		modelInput.Unicom,
		modelInput.CTAF,
		modelInput.EffectiveDate,
		now,
		now,
	)
	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(newID).
		WillReturnRows(rows)

	// Expect commit
	mock.ExpectCommit()

	// Act
	out, err := repo.Insert(context.Background(), tx, modelInput)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.ID)
	assert.Equal(t, newID, out.ID.String())
	assert.Equal(t, modelInput.SiteNumber, out.SiteNumber)
	assert.Equal(t, modelInput.ICAOID, out.ICAOID)
	assert.Equal(t, modelInput.FAAID, out.FAAID)
	assert.Equal(t, modelInput.IATAID, out.IATAID)
	assert.Equal(t, modelInput.Name, out.Name)
	assert.Equal(t, modelInput.Type, out.Type)
	assert.Equal(t, modelInput.Status, out.Status)
	assert.Equal(t, modelInput.Country, out.Country)
	assert.Equal(t, modelInput.State, out.State)
	assert.Equal(t, modelInput.StateFull, out.StateFull)
	assert.Equal(t, modelInput.County, out.County)
	assert.Equal(t, modelInput.City, out.City)
	assert.Equal(t, modelInput.Ownership, out.Ownership)
	assert.Equal(t, modelInput.Use, out.Use)
	assert.Equal(t, modelInput.Manager, out.Manager)
	assert.Equal(t, modelInput.ManagerPhone, out.ManagerPhone)
	assert.Equal(t, modelInput.Latitude, out.Latitude)
	assert.Equal(t, modelInput.LatitudeSec, out.LatitudeSec)
	assert.Equal(t, modelInput.Longitude, out.Longitude)
	assert.Equal(t, modelInput.LongitudeSec, out.LongitudeSec)
	assert.Equal(t, modelInput.Elevation, out.Elevation)
	assert.Equal(t, modelInput.ControlTower, out.ControlTower)
	assert.Equal(t, modelInput.Unicom, out.Unicom)
	assert.Equal(t, modelInput.CTAF, out.CTAF)
	assert.NotNil(t, out.CreatedAt)
	assert.NotNil(t, out.UpdatedAt)
	assert.WithinDuration(t, now, *out.CreatedAt, time.Second)
	assert.WithinDuration(t, now, *out.UpdatedAt, time.Second)
}

func TestAirportRepository_SyncAirport(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := &AirportRepository{logger: log}

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	defer util.CommitOrRollback(tx)

	// ---------- arrange input ----------
	req := airport_dto.AirportRequestDto{
		SiteNumber:    util.Ptr("12345"),
		ICAOID:        util.Ptr("KJFK"),
		FAAID:         util.Ptr("JFK"),
		IATAID:        util.Ptr("JFK"),
		Name:          util.Ptr("John F. Kennedy International Airport"),
		Type:          enum.AIRPORT,
		Status:        util.Ptr(true),
		Country:       util.Ptr("USA"),
		State:         util.Ptr("NY"),
		StateFull:     util.Ptr("New York"),
		County:        util.Ptr("Queens"),
		City:          util.Ptr("New York"),
		Ownership:     enum.OWN_PUBLIC,
		Use:           enum.USE_PUBLIC,
		Manager:       util.Ptr("Jane Doe"),
		ManagerPhone:  util.Ptr("+1-555-1234"),
		Latitude:      util.Ptr("40.6413 N"),
		LatitudeSec:   util.Ptr("38.0"),
		Longitude:     util.Ptr("73.7781 W"),
		LongitudeSec:  nil, // keep it nil to test
		Elevation:     util.Ptr(int64(13)),
		ControlTower:  util.Ptr(true),
		Unicom:        util.Ptr("123.45"),
		CTAF:          util.Ptr("123.45"),
		EffectiveDate: nil,
	}
	modelInput := airport_dto.AirportRequestToAirport(req)

	// ---------- expect INSERT ----------
	insertRe := regexp.MustCompile(`(?s)INSERT\s+INTO\s+airports\s*\(.*?\)\s*VALUES\s*\(.*?\)\s*RETURNING\s+id`)
	newID := uuid.New().String()

	mock.ExpectQuery(insertRe.String()).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newID))

	// ---------- expect SELECT (FindByID) ----------
	cols := newCols()
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(newID).
		WillReturnRows(
			sqlmock.NewRows(cols).AddRow(
				newID,
				modelInput.SiteNumber,
				modelInput.ICAOID,
				modelInput.FAAID,
				modelInput.IATAID,
				modelInput.Name,
				modelInput.Type,
				modelInput.Status,
				modelInput.Country,
				modelInput.State,
				modelInput.StateFull,
				modelInput.County,
				modelInput.City,
				modelInput.Ownership,
				modelInput.Use,
				modelInput.Manager,
				modelInput.ManagerPhone,
				modelInput.Latitude,
				modelInput.LatitudeSec,
				modelInput.Longitude,
				modelInput.LongitudeSec,
				modelInput.Elevation,
				modelInput.ControlTower,
				modelInput.Unicom,
				modelInput.CTAF,
				modelInput.EffectiveDate,
				now,
				now,
			),
		)

	// Expect commit
	mock.ExpectCommit()

	// Act
	out, err := repo.SyncAirport(context.Background(), tx, modelInput)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.ID)
	assert.Equal(t, newID, out.ID.String())
	assert.Equal(t, modelInput.SiteNumber, out.SiteNumber)
	assert.Equal(t, modelInput.ICAOID, out.ICAOID)
	assert.Equal(t, modelInput.FAAID, out.FAAID)
	assert.Equal(t, modelInput.IATAID, out.IATAID)
	assert.Equal(t, modelInput.Name, out.Name)
	assert.Equal(t, modelInput.Type, out.Type)
	assert.Equal(t, modelInput.Status, out.Status)
	assert.Equal(t, modelInput.Country, out.Country)
	assert.Equal(t, modelInput.State, out.State)
	assert.Equal(t, modelInput.StateFull, out.StateFull)
	assert.Equal(t, modelInput.County, out.County)
	assert.Equal(t, modelInput.City, out.City)
	assert.Equal(t, modelInput.Ownership, out.Ownership)
	assert.Equal(t, modelInput.Use, out.Use)
	assert.Equal(t, modelInput.Manager, out.Manager)
	assert.Equal(t, modelInput.ManagerPhone, out.ManagerPhone)
	assert.Equal(t, modelInput.Latitude, out.Latitude)
	assert.Equal(t, modelInput.LatitudeSec, out.LatitudeSec)
	assert.Equal(t, modelInput.Longitude, out.Longitude)
	assert.Equal(t, modelInput.LongitudeSec, out.LongitudeSec)
	assert.Equal(t, modelInput.Elevation, out.Elevation)
	assert.Equal(t, modelInput.ControlTower, out.ControlTower)
	assert.Equal(t, modelInput.Unicom, out.Unicom)
	assert.Equal(t, modelInput.CTAF, out.CTAF)
	assert.NotNil(t, out.CreatedAt)
	assert.NotNil(t, out.UpdatedAt)
	assert.WithinDuration(t, now, *out.CreatedAt, time.Second)
	assert.WithinDuration(t, now, *out.UpdatedAt, time.Second)
}

// ---------- UNIT TESTS FOR FindByID ----------
func TestAirportRepository_FindByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := NewAirportRepository(log)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	defer util.CommitOrRollback(tx)

	// ---------- data dummy ----------
	for i := range dataDummy {
		data := dataDummy[i]
		id := data.id

		// Expect query for each data
		row := data.row
		mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
			WithArgs(id).
			WillReturnRows(
				successRow(
					&id,
					row.SiteNumber, row.ICAOID, row.FAAID, row.IATAID, row.Name,
					row.Type,
					row.Status,
					row.Country, row.State, row.StateFull, row.County, row.City,
					row.Ownership,
					row.Use,
					row.Manager, row.ManagerPhone, row.Latitude, row.LatitudeSec, row.Longitude, row.LongitudeSec,
					row.Elevation,
					row.ControlTower,
					row.Unicom, row.CTAF,
					row.EffectiveDate,
					timeNow, timeNow,
				),
			)
	}

	// Expect commit
	mock.ExpectCommit()

	// Act - Assert
	ctx := context.Background()
	for _, test := range dataDummy {
		out, err := repo.FindByID(ctx, tx, test.id.String())
		assert.NoError(t, err, "unexpected error for %s", test.label)
		assert.NotNil(t, out, "nil output for %s", test.label)
		assert.Equal(t, test.id.String(), out.ID.String(), "mismatched ID for %s", test.label)
		assert.NotNil(t, out.CreatedAt, "nil CreatedAt for %s", test.label)
		assert.NotNil(t, out.UpdatedAt, "nil UpdatedAt for %s", test.label)
		assert.WithinDuration(t, timeNow, *out.CreatedAt, time.Second, "CreatedAt not within duration for %s", test.label)
		assert.WithinDuration(t, timeNow, *out.UpdatedAt, time.Second, "UpdatedAt not within duration for %s", test.label)
	}
}

func TestAirportRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := NewAirportRepository(log)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	defer util.CommitOrRollback(tx)

	// id valid tapi tidak ada row
	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(newCols())) // 0 row

	mock.ExpectCommit()

	_, err = repo.FindByID(context.Background(), tx, id.String())
	assert.ErrorIs(t, err, util.ErrNotFound)
}

func TestAirportRepository_FindByID_InvalidUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := NewAirportRepository(log)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	defer util.CommitOrRollback(tx)

	// INVALID UUID → fungsi return ErrNotFound sebelum query.
	// Jadi JANGAN set ExpectQuery apapun di test ini.

	mock.ExpectCommit()

	_, err = repo.FindByID(context.Background(), tx, "12345")
	assert.ErrorIs(t, err, util.ErrNotFound)
}

func TestAirportRepository_FindByID_ErrNoRows_FromQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	repo := NewAirportRepository(log)

	// begin tx
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	// expect rollback
	mock.ExpectRollback()

	// act
	_, gotErr := repo.FindByID(context.Background(), tx, id.String())
	assert.ErrorIs(t, gotErr, util.ErrNotFound)
	_ = tx.Rollback()
}

// ---------- UNIT TESTS FOR FindAll ----------
func TestAirportRepository_FindAll(t *testing.T) {
	// table-driven: tiap case punya limit/offset & ekspektasi panjang output
	cases := []struct {
		name        string
		limit       int
		offset      int
		expectedLen int
	}{
		{"explicit 10/0 (clamped by data)", 10, 0, len(dataDummy)},
		{"defaults 10/0 (nil di prod, kita treat sama)", 10, 0, len(dataDummy)},
		{"paged 2/2", 2, 2, len(dataDummy) - 2},
	}

	// query string yang dipakai repo
	selectAll := `
		SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status, created_at, updated_at
		FROM airports 
		ORDER BY icao_id
		LIMIT $1
		OFFSET $2`
	selectAllQ := regexp.QuoteMeta(strings.TrimSpace(selectAll))
	countQ := regexp.MustCompile(`SELECT\s+COUNT\(\*\)\s+FROM\s+airports`)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			// begin tx
			mock.ExpectBegin()
			var tx *sql.Tx
			tx, err = db.Begin()
			assert.NoError(t, err)

			// SELECT expectation (sesuai limit/offset)
			mock.ExpectQuery(selectAllQ).
				WithArgs(tc.limit, tc.offset).
				WillReturnRows(buildRowsFindAll(tc.limit, tc.offset))

			// COUNT expectation
			mock.ExpectQuery(countQ.String()).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(dataDummy)))

			// commit
			mock.ExpectCommit()

			// -------------------
			// ACT
			ctx := context.Background()
			repo := NewAirportRepository(log)
			args := map[string]interface{}{"limit": tc.limit, "offset": tc.offset}

			// -------------------
			airports, _, err := repo.FindAll(ctx, tx, args) // panggil method punyamu
			assert.NoError(t, err)
			assert.NotNil(t, airports)
			assert.Equal(t, tc.expectedLen, len(airports), "mismatched length")

			for i := range airports {
				expected := dataDummy[tc.offset+i]
				got := airports[i]
				assert.Equal(t, expected.id, *got.ID)
				assert.Equal(t, *expected.row.SiteNumber, *got.SiteNumber)
				assert.Equal(t, *expected.row.ICAOID, *got.ICAOID)
				assert.Equal(t, *expected.row.FAAID, *got.FAAID)
				assert.Equal(t, *expected.row.IATAID, *got.IATAID)
				assert.Equal(t, *expected.row.Name, *got.Name)
				assert.Equal(t, expected.row.Type, got.Type)
				assert.Equal(t, *expected.row.Status, *got.Status)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
				assert.WithinDuration(t, timeNow, *got.CreatedAt, time.Second)
				assert.WithinDuration(t, timeNow, *got.UpdatedAt, time.Second)
			}

			// selesai -> commit
			assert.NoError(t, tx.Commit())
		})
	}
}

// ---------- UNIT TEST For FindBySearchName ---------
func TestAirportRepository_FindBySearchName(t *testing.T) {
	// query dari repo
	selectSQL := `
SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status,
        country, state, state_full, county, city, ownership, "use",
        manager, manager_phone, latitude, latitude_sec, longitude, longitude_sec, elevation,
        control_tower, unicom, ctaf, effective_date, created_at, updated_at
FROM airports 
WHERE LOWER(name) LIKE LOWER($3)
ORDER BY icao_id
LIMIT $1
OFFSET $2`
	selectQ := regexp.QuoteMeta(strings.TrimSpace(selectSQL))
	countRe := regexp.MustCompile(`(?is)SELECT\s+COUNT\(\*\)\s+FROM\s+airports\s+WHERE\s+name\s+ILIKE\s+\$1`)

	cases := []struct {
		name        string
		search      string
		limit       int
		offset      int
		expectLen   int
		expectTotal int
	}{
		{"match 'Intl' first page", "Intl", 2, 0, 2, 2},
		{"match 'a' second page", "a", 2, 2, 1, 3},
		{"case-insensitive", "INTL", 10, 0, 2, 2},
		{"no match", "zzz", 10, 0, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			mock.ExpectBegin()
			tx, err := db.Begin()
			assert.NoError(t, err)

			searchPattern := "%" + c.search + "%"

			// SELECT expectation
			rows, total := buildRowsByName(searchPattern, c.limit, c.offset)

			// args: LIMIT $1, OFFSET $2, LIKE $3
			margs := []driver.Value{c.limit, c.offset, searchPattern}
			mock.ExpectQuery(selectQ).
				WithArgs(margs...).
				WillReturnRows(rows)

			// COUNT expectation
			mock.ExpectQuery(countRe.String()).
				WithArgs(searchPattern).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))

			mock.ExpectCommit()

			repo := NewAirportRepository(log) // sesuaikan
			args := map[string]interface{}{"limit": c.limit, "offset": c.offset}

			out, gotTotal, err := repo.FindBySearchName(context.Background(), tx, c.search, args)
			assert.NoError(t, err)
			assert.Equal(t, c.expectTotal, gotTotal)
			assert.Equal(t, c.expectLen, len(out))

			assert.NoError(t, tx.Commit())
		})
	}
}

// ---------- UNIT TESTS FOR FindByICAO ----------
func TestAirportRepository_FindByICAOID(t *testing.T) {
	query := `SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status,
				country, state, state_full, county, city, ownership, "use",
				manager, manager_phone, latitude, latitude_sec, longitude, longitude_sec, elevation,
				control_tower, unicom, ctaf, effective_date, created_at, updated_at
		FROM airports 
		WHERE icao_id = $1 
		LIMIT 1`
	q := regexp.QuoteMeta(strings.TrimSpace(query))

	cases := []struct {
		name        string
		icao        string
		setupMock   func(sqlmock.Sqlmock)
		expectErr   error
		expectOK    bool
		expectPanic bool
	}{
		{
			name: "success",
			icao: "KSFO",
			setupMock: func(m sqlmock.Sqlmock) {
				data := dataDummy[2].row
				rows := sqlmock.NewRows(newCols()).AddRow(
					data.ID, data.SiteNumber, data.ICAOID, data.FAAID, data.IATAID, data.Name, data.Type, data.Status,
					data.Country, data.State, data.StateFull, data.County, data.City, data.Ownership, data.Use,
					data.Manager, data.ManagerPhone, data.Latitude, data.LatitudeSec, data.Longitude, data.LongitudeSec, data.Elevation,
					data.ControlTower, data.Unicom, data.CTAF, data.EffectiveDate, data.CreatedAt, data.UpdatedAt,
				)
				m.ExpectQuery(q).
					WithArgs("KSFO").
					WillReturnRows(rows)
			},
			expectErr: nil,
			expectOK:  true,
		},
		{
			name: "not found (empty result)",
			icao: "KZZZ",
			setupMock: func(m sqlmock.Sqlmock) {
				empty := sqlmock.NewRows(newCols()) // 0 row
				m.ExpectQuery(q).
					WithArgs("KZZZ").
					WillReturnRows(empty)
			},
			expectErr: util.ErrNotFound,
			expectOK:  false,
		},
		{
			name: "db error -> panic",
			icao: "KERR",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(q).
					WithArgs("KERR").
					WillReturnError(fmt.Errorf("boom"))
			},
			expectErr: nil,   // kita assert panic, bukan error return
			expectOK:  false, // tidak dipakai
		},
		{
			name: "query returns sql.ErrNoRows",
			icao: "KNONE",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(q).
					WithArgs("KNONE").
					WillReturnError(sql.ErrNoRows) // harus map ke util.ErrNotFound
			},
			expectErr:   util.ErrNotFound,
			expectOK:    false,
			expectPanic: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			mock.ExpectBegin()
			tx, err := db.Begin()
			assert.NoError(t, err)

			tc.setupMock(mock)
			mock.ExpectCommit()

			repo := NewAirportRepository(log)

			if tc.name == "db error -> panic" {
				assert.Panics(t, func() {
					_, _ = repo.FindByICAOID(context.Background(), tx, tc.icao)
				})
				assert.NoError(t, tx.Commit())
				return
			}

			if tc.expectPanic {
				assert.Panics(t, func() {
					_, _ = repo.FindByICAOID(context.Background(), tx, tc.icao)
				})
			}

			got, err := repo.FindByICAOID(context.Background(), tx, tc.icao)
			assert.ErrorIs(t, err, tc.expectErr)
			if tc.expectOK {
				assert.Equal(t, "KSFO", *got.ICAOID)
				assert.Equal(t, "San Francisco Intl", *got.Name)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
			}

			assert.NoError(t, tx.Commit())
		})
	}
}

// ---------- UNIT TESTS FOR FindExistsByICAOID ----------
func TestAirportRepository_FindExistsByICAOID(t *testing.T) {
	updateQ := regexp.QuoteMeta(`SELECT 1 FROM airports WHERE icao_id = $1 LIMIT 1`)

	cases := []struct {
		name      string
		icao      string
		setupMock func(sqlmock.Sqlmock)
		expectOK  bool
		expectErr bool // dipakai untuk skenario panic
	}{
		{
			name: "exists",
			icao: "KAAA",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(1)
				m.ExpectQuery(updateQ).
					WithArgs("KAAA").
					WillReturnRows(rows)
			},
			expectOK:  true,
			expectErr: false,
		},
		{
			name: "not found",
			icao: "KZZZ",
			setupMock: func(m sqlmock.Sqlmock) {
				// QueryRowContext -> Scan akan terima sql.ErrNoRows
				m.ExpectQuery(updateQ).
					WithArgs("KZZZ").
					WillReturnError(sql.ErrNoRows)
				// Alternatif: WillReturnRows(sqlmock.NewRows([]string{"exists"})) juga oke
			},
			expectOK:  false,
			expectErr: false,
		},
		{
			name: "db error -> panic",
			icao: "KERR",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(updateQ).
					WithArgs("KERR").
					WillReturnError(fmt.Errorf("boom"))
			},
			expectOK:  false,
			expectErr: true, // kita expect panic
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			mock.ExpectBegin()
			tx, err := db.Begin()
			assert.NoError(t, err)

			tc.setupMock(mock)
			mock.ExpectCommit()

			repo := NewAirportRepository(log)

			if tc.expectErr {
				assert.Panics(t, func() {
					_, _ = repo.FindExistsByICAOID(context.Background(), tx, tc.icao)
				})
			} else {
				ok, err := repo.FindExistsByICAOID(context.Background(), tx, tc.icao)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectOK, ok)
			}

			assert.NoError(t, tx.Commit())
		})
	}
}

// ---------- UNIT TESTS FOR Update ----------
func TestAirportRepository_Update(t *testing.T) {
	// --- arrange data
	existingID := dataDummy[1].id
	existingIDStr := dataDummy[1].id.String()
	notFoundID := uuid.New().String()
	invalidID := "invalid-uuid"

	// Existing airport data needed for update
	// existingAirport := dataDummy[1].row
	updatedAirport := dataDummy[1].row

	// Update input
	// Only manager and manager_phone property
	newManager := "Updated manager"
	newManagerPhone := "1-xxx-xxx-xx"
	updateInput := airport_dto.AirportUpdateDto{
		Manager:      &newManager,
		ManagerPhone: &newManagerPhone,
	}

	util.FillUpdatableFields(&updatedAirport, updateInput)

	// Time updated
	timeUpdated := time.Now()

	// SQL update (trimmed di repo)
	updateSQL := `
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
	updateQ := regexp.QuoteMeta(strings.TrimSpace(updateSQL))

	cases := []struct {
		name      string
		id        string
		payload   model.Airport
		expectErr error
		setupMock func(mock sqlmock.Sqlmock)
	}{
		{
			name:      "success",
			id:        existingIDStr,
			payload:   updatedAirport,
			expectErr: nil,
			setupMock: func(m sqlmock.Sqlmock) {
				args := make([]driver.Value, 0, 25)
				for i := 0; i < 24; i++ {
					args = append(args, sqlmock.AnyArg())
				}
				args = append(args, existingID)

				ret := sqlmock.NewRows([]string{"id"}).AddRow(existingID)
				m.ExpectQuery(updateQ).
					WithArgs(args...).
					WillReturnRows(ret)

					// FindByID
				rows := sqlmock.NewRows(newCols()).AddRow(
					existingID,
					updatedAirport.SiteNumber,
					updatedAirport.ICAOID,
					updatedAirport.FAAID,
					updatedAirport.IATAID,
					updatedAirport.Name,
					updatedAirport.Type,
					updatedAirport.Status,
					updatedAirport.Country,
					updatedAirport.State,
					updatedAirport.StateFull,
					updatedAirport.County,
					updatedAirport.City,
					updatedAirport.Ownership,
					updatedAirport.Use,
					updatedAirport.Manager,
					updatedAirport.ManagerPhone,
					updatedAirport.Latitude,
					updatedAirport.LatitudeSec,
					updatedAirport.Longitude,
					updatedAirport.LongitudeSec,
					updatedAirport.Elevation,
					updatedAirport.ControlTower,
					updatedAirport.Unicom,
					updatedAirport.CTAF,
					updatedAirport.EffectiveDate,
					timeNow,
					timeUpdated,
				)
				m.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
					WithArgs(existingID).
					WillReturnRows(rows)
			},
		},
		{
			name:      "not found (no rows returned)",
			id:        notFoundID,
			payload:   updatedAirport,
			expectErr: util.ErrNotFound,
			setupMock: func(m sqlmock.Sqlmock) {
				// susun args: 24 field + uuid di posisi $25
				args := make([]driver.Value, 0, 25)
				for i := 0; i < 24; i++ {
					args = append(args, sqlmock.AnyArg())
				}
				args = append(args, notFoundID)

				m.ExpectQuery(updateQ).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows) // bikin row.Scan() return ErrNoRows
				// tidak ada FindByID
			},
		},
		{
			name:      "invalid uuid",
			id:        invalidID,
			payload:   updatedAirport,
			expectErr: util.ErrNotFound,
			setupMock: func(m sqlmock.Sqlmock) {
				// Tidak ada query karena gagal parse sebelum QueryRowContext
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			mock.ExpectBegin()
			tx, err := db.Begin()
			defer util.CommitOrRollback(tx)
			assert.NoError(t, err)

			// set ekspektasi per skenario
			c.setupMock(mock)

			// Expect Commit
			mock.ExpectCommit()

			repo := NewAirportRepository(log)

			got, err := repo.Update(context.Background(), tx, c.id, c.payload)
			assert.ErrorIs(t, err, c.expectErr)
			if c.expectErr == nil {
				assert.Equal(t, existingID.String(), got.ID.String())
			}

			if c.expectErr == nil {
				assert.Equal(t, existingIDStr, got.ID.String())
				assert.Equal(t, newManager, *got.Manager)
				assert.Equal(t, newManagerPhone, *got.ManagerPhone)
				assert.Equal(t, timeUpdated, *got.UpdatedAt)
				assert.Equal(t, timeNow, *got.CreatedAt)
			}
		})
	}
}

// ---------- UNIT TESTS FOR Delete ----------
func TestAirportRepository_Delete(t *testing.T) {
	// data dummy
	deletedIdSuccess := dataDummy[0].id // ID yang ada
	deletedIdNotFound := uuid.New()     // ID yang tidak ada
	deletedIdInvalid := "invalid-uuid"  // invalid format

	cases := []struct {
		name        string
		id          string
		expectedErr error
	}{
		{"existing ID", deletedIdSuccess.String(), nil},
		{"non-existing ID", deletedIdNotFound.String(), util.ErrNotFound},
		{"invalid UUID", deletedIdInvalid, util.ErrNotFound},
	}

	// siapkan regex query DELETE
	deleteRe := regexp.MustCompile(`(?s)DELETE\s+FROM\s+airports\s+WHERE\s+id\s*=\s*\$1`)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, mock.ExpectationsWereMet())
				_ = db.Close()
			}()

			// begin tx
			mock.ExpectBegin()
			tx, err := db.Begin()
			assert.NoError(t, err)

			// Pasang expectation hanya jika UUID valid (karena kodemu parse dulu sebelum Exec)
			if uid, err := uuid.Parse(tc.id); err == nil {
				affected := int64(0)
				if tc.expectedErr == nil {
					affected = 1 // sukses: 1 row affected
				}
				mock.ExpectExec(deleteRe.String()).
					WithArgs(uid).
					WillReturnResult(sqlmock.NewResult(0, affected))
			}
			// commit selalu (repo.Delete tidak commit/rollback)
			mock.ExpectCommit()

			// ACT
			repo := NewAirportRepository(log)
			err = repo.Delete(context.Background(), tx, tc.id)

			// ASSERT
			assert.ErrorIs(t, err, tc.expectedErr)

			// finalize tx
			assert.NoError(t, tx.Commit())
		})
	}
}
