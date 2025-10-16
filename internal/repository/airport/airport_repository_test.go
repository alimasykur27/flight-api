package repository_airport

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
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

// ------------ QUERY ---------------
var insertAirportQuery string = `
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

var syncAirportQuery string = `
	INSERT INTO airports (
		site_number, icao_id, faa_id, iata_id, name, 
		type, status, country, state, state_full, 
		county, city, ownership, "use", manager, 
		manager_phone, latitude, latitude_sec, longitude, longitude_sec,
		elevation, control_tower, unicom, ctaf, effective_date,
		sync_status, sync_message
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12,
		$13, $14, $15, $16,
		$17, $18, $19, $20,
		$21, $22, $23, $24, $25,
		20, 'synced'
	)
	RETURNING 
		id,
		site_number, icao_id, faa_id, iata_id, name, 
		type, status, country, state, state_full, 
		county, city, ownership, "use", manager, 
		manager_phone, latitude, latitude_sec, longitude, longitude_sec,
		elevation, control_tower, unicom, ctaf, effective_date,
		sync_status, sync_message, updated_at, created_at
`

var selectAirportQuery string = `
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

var findBySearchNameQuery string = `
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

// ---------- HELPER FUNCTIONS ----------
func airportTableCols() []string {
	return []string{
		"id", "site_number", "icao_id", "faa_id", "iata_id", "name",
		"type", "status", "country", "state", "state_full", "county", "city",
		"ownership", "use", "manager", "manager_phone",
		"latitude", "latitude_sec", "longitude", "longitude_sec",
		"elevation", "control_tower", "unicom", "ctaf", "effective_date",
		"sync_status", "sync_message", "updated_at", "created_at",
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
	sync_status enum.SyncStatusEnum,
	sync_message *string,
	updatedAt, createdAt time.Time,
) *sqlmock.Rows {
	return sqlmock.NewRows(airportTableCols()).
		AddRow(
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
			sync_status,
			sync_message,
			updatedAt, createdAt,
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

	rows := sqlmock.NewRows(airportTableCols())
	for _, a := range page {
		rows.AddRow(
			a.ID,
			a.SiteNumber, a.ICAOID, a.FAAID, a.IATAID, a.Name,
			a.Type, a.Status, a.Country, a.State, a.StateFull,
			a.County, a.City, a.Ownership, a.Use, a.Manager,
			a.ManagerPhone, a.Latitude, a.LatitudeSec, a.Longitude, a.LongitudeSec,
			a.Elevation, a.ControlTower, a.Unicom, a.CTAF, a.EffectiveDate,
			a.SyncStatus, a.SyncMessage, a.UpdatedAt, a.CreatedAt,
		)
	}
	return rows, total
}

func SetupTesting(t *testing.T) (*logger.Logger, context.Context, *sql.DB, sqlmock.Sqlmock, func(), *AirportRepository) {
	log := logger.NewLogger(logger.DEBUG_LEVEL)

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	cleanup := func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}

	repoMock := &AirportRepository{logger: log}
	return log, ctx, db, mock, cleanup, repoMock
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
	_, _, _, _, _, repoMock := SetupTesting(t)
	assert.NotNil(t, repoMock)
	assert.IsType(t, &AirportRepository{}, repoMock)
	assert.NotNil(t, repoMock.logger)
	assert.IsType(t, &logger.Logger{}, repoMock.logger)
}

func TestAirportRepository_Insert(t *testing.T) {
	_, ctx, db, mock, cleanup, repoMock := SetupTesting(t)
	defer cleanup()

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
	newID := uuid.New().String()

	// Expectation
	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(strings.TrimSpace(insertAirportQuery)),
	).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows(airportTableCols()).
				AddRow(
					newID,
					modelInput.SiteNumber, modelInput.ICAOID, modelInput.FAAID, modelInput.IATAID, modelInput.Name,
					modelInput.Type, modelInput.Status, modelInput.Country, modelInput.State, modelInput.StateFull,
					modelInput.County, modelInput.City, modelInput.Ownership, modelInput.Use, modelInput.Manager,
					modelInput.ManagerPhone, modelInput.Latitude, modelInput.LatitudeSec, modelInput.Longitude, modelInput.LongitudeSec,
					modelInput.Elevation, modelInput.ControlTower, modelInput.Unicom, modelInput.CTAF, modelInput.EffectiveDate,
					enum.SYNC_NEW, enum.SYNC_NEW.String(), timeNow, timeNow,
				),
		)
	mock.ExpectCommit()

	// Act
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	assert.NoError(t, err)
	defer util.FinishTx(tx, &err)

	out, errRepo := repoMock.Insert(ctx, tx, modelInput)

	// Assert
	assert.Nil(t, errRepo)
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
	assert.WithinDuration(t, timeNow, *out.CreatedAt, time.Second)
	assert.WithinDuration(t, timeNow, *out.UpdatedAt, time.Second)
}

// ---------- UNIT TESTS FOR FindByID ----------
func TestAirportRepository_FindByID_Success(t *testing.T) {
	_, ctx, db, mock, cleanup, repoMock := SetupTesting(t)
	defer cleanup()

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
					enum.SYNC_NEW,
					util.Ptr(enum.SYNC_NEW.String()),
					timeNow,
					timeNow,
				),
			)
	}

	// Act - Assert
	for _, test := range dataDummy {
		out, err := repoMock.FindByID(ctx, db, test.id.String())
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
	_, _, db, mock, cleanup, repoMock := SetupTesting(t)
	defer cleanup()

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(airportTableCols()))

	_, err := repoMock.FindByID(context.Background(), db, id.String())
	assert.ErrorIs(t, err, util.ErrNotFound)
}

func TestAirportRepository_FindByID_InvalidUUID(t *testing.T) {
	_, _, db, _, cleanup, repoMock := SetupTesting(t)
	defer cleanup()

	_, err := repoMock.FindByID(context.Background(), db, "12345")
	assert.ErrorIs(t, err, util.ErrNotFound)
}

func TestAirportRepository_FindByID_ErrNoRows(t *testing.T) {
	_, _, db, mock, cleanup, repoMock := SetupTesting(t)
	defer cleanup()

	// arrange
	id := uuid.New()

	// expect
	mock.ExpectQuery(regexp.QuoteMeta(selectAirportQuery)).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	// act
	_, gotErr := repoMock.FindByID(context.Background(), db, id.String())
	assert.ErrorIs(t, gotErr, util.ErrNotFound)
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
	selectAllQuery := `
		SELECT id, site_number, icao_id, faa_id, iata_id, name, type, status, created_at, updated_at
		FROM airports 
		ORDER BY icao_id
		LIMIT $1
		OFFSET $2
	`
	countQ := regexp.MustCompile(`SELECT\s+COUNT\(\*\)\s+FROM\s+airports`)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// -------------------
			// ARRANGE
			_, _, db, mock, cleanup, repoMock := SetupTesting(t)
			defer cleanup()

			// COUNT expectation
			mock.ExpectQuery(countQ.String()).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(dataDummy)))

			// SELECT expectation (sesuai limit/offset)
			mock.ExpectQuery(regexp.QuoteMeta(strings.TrimSpace(selectAllQuery))).
				WithArgs(tc.limit, tc.offset).
				WillReturnRows(buildRowsFindAll(tc.limit, tc.offset))

			// ACT
			ctx := context.Background()
			args := map[string]interface{}{"limit": tc.limit, "offset": tc.offset}
			airports, _, err := repoMock.FindAll(ctx, db, args)

			// -------------------
			// ASSERT
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
		})
	}
}

// ---------- UNIT TEST For FindBySearchName ---------
func TestAirportRepository_FindBySearchName(t *testing.T) {
	// query dari repo
	selectQ := regexp.QuoteMeta(strings.TrimSpace(findBySearchNameQuery))
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
			_, _, db, mock, cleanup, repoMock := SetupTesting(t)
			defer cleanup()

			// search pattern
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

			args := map[string]interface{}{"limit": c.limit, "offset": c.offset}

			out, gotTotal, err := repoMock.FindBySearchName(context.Background(), db, c.search, args)
			assert.NoError(t, err)
			assert.Equal(t, c.expectTotal, gotTotal)
			assert.Equal(t, c.expectLen, len(out))
		})
	}
}

// ---------- UNIT TESTS FOR FindByICAO ----------
func TestAirportRepository_FindByICAOID(t *testing.T) {
	selectByICAOID := `
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
	query := regexp.QuoteMeta(strings.TrimSpace(selectByICAOID))

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
				rows := sqlmock.NewRows(airportTableCols()).
					AddRow(
						data.ID,
						data.SiteNumber, data.ICAOID, data.FAAID, data.IATAID, data.Name,
						data.Type, data.Status, data.Country, data.State, data.StateFull,
						data.County, data.City, data.Ownership, data.Use, data.Manager,
						data.ManagerPhone, data.Latitude, data.LatitudeSec, data.Longitude, data.LongitudeSec,
						data.Elevation, data.ControlTower, data.Unicom, data.CTAF, data.EffectiveDate,
						enum.SYNC_NEW, enum.SYNC_NEW.String(), timeNow, timeNow,
					)
				m.ExpectQuery(query).
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
				empty := sqlmock.NewRows(airportTableCols())
				m.ExpectQuery(query).
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
				m.ExpectQuery(query).
					WithArgs("KERR").
					WillReturnError(fmt.Errorf("boom"))
			},
			expectErr: nil,
			expectOK:  false,
		},
		{
			name: "query returns sql.ErrNoRows",
			icao: "KNONE",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(query).
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
			_, _, db, mock, cleaup, repoMock := SetupTesting(t)
			defer cleaup()

			// Setup mock by case
			tc.setupMock(mock)

			if tc.name == "db error -> panic" {
				assert.Panics(t, func() {
					_, _ = repoMock.FindByICAOID(context.Background(), db, tc.icao)
				})
				return
			}

			if tc.expectPanic {
				assert.Panics(t, func() {
					_, _ = repoMock.FindByICAOID(context.Background(), db, tc.icao)
				})
			}

			got, err := repoMock.FindByICAOID(context.Background(), db, tc.icao)
			assert.ErrorIs(t, err, tc.expectErr)
			if tc.expectOK {
				assert.Equal(t, "KSFO", *got.ICAOID)
				assert.Equal(t, "San Francisco Intl", *got.Name)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
			}
		})
	}
}

// ---------- UNIT TESTS FOR FindExistsByICAOID ----------
func TestAirportRepository_FindExistsByICAOID(t *testing.T) {
	checkExistsQuery := regexp.QuoteMeta(`SELECT 1 FROM airports WHERE icao_id = $1 LIMIT 1`)

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
				m.ExpectQuery(checkExistsQuery).
					WithArgs("KAAA").
					WillReturnRows(rows)
			},
			expectOK:  true,
			expectErr: false,
		},
		{
			name: "exists is not 1",
			icao: "KBBB",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(99) // nilai selain 1
				m.ExpectQuery(checkExistsQuery).
					WithArgs("KBBB").
					WillReturnRows(rows)
			},
			expectOK:  false,
			expectErr: false,
		},
		{
			name: "not found",
			icao: "KZZZ",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(checkExistsQuery).
					WithArgs("KZZZ").
					WillReturnError(sql.ErrNoRows)
			},
			expectOK:  false,
			expectErr: false,
		},
		{
			name: "db error -> panic",
			icao: "KERR",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(checkExistsQuery).
					WithArgs("KERR").
					WillReturnError(util.ErrInternalServer)
			},
			expectOK:  false,
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, ctx, db, mock, cleanup, repoMock := SetupTesting(t)
			defer cleanup()

			// Setup mock by case
			tc.setupMock(mock)

			if tc.expectErr {
				_, err := repoMock.FindExistsByICAOID(ctx, db, tc.icao)
				assert.Contains(t, err.Error(), util.ErrInternalServer.Error())
			} else {
				ok, err := repoMock.FindExistsByICAOID(ctx, db, tc.icao)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectOK, *ok)
				assert.NotNil(t, ok)
			}
		})
	}
}

// ---------- UNIT TESTS FOR Update ----------
func TestAirportRepository_Update(t *testing.T) {
	// --- Arrange data ---
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

	// Update Query
	updateSQL := `
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

				ret := sqlmock.NewRows(airportTableCols()).
					AddRow(
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
						updatedAirport.SyncStatus,
						updatedAirport.SyncMessage,
						timeNow,
						timeUpdated,
					)
				m.ExpectQuery(updateQ).
					WithArgs(args...).
					WillReturnRows(ret)

				m.ExpectCommit()
			},
		},
		{
			name:      "not found (no rows returned)",
			id:        notFoundID,
			payload:   updatedAirport,
			expectErr: util.ErrNotFound,
			setupMock: func(m sqlmock.Sqlmock) {
				args := make([]driver.Value, 0, 25)
				for i := 0; i < 24; i++ {
					args = append(args, sqlmock.AnyArg())
				}
				args = append(args, notFoundID)

				m.ExpectQuery(updateQ).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)

				m.ExpectRollback()
			},
		},
		{
			name:      "invalid uuid",
			id:        invalidID,
			payload:   updatedAirport,
			expectErr: util.ErrBadRequest,
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectRollback()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, ctx, db, mock, cleanup, repoMock := SetupTesting(t)
			defer cleanup()

			// Setup Begin
			mock.ExpectBegin()

			// Setup mock by case
			c.setupMock(mock)

			// ACT
			tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
			assert.NoError(t, err)
			defer util.FinishTx(tx, &err)

			// ASSERT
			got, err := repoMock.Update(ctx, tx, c.id, c.payload)

			if c.expectErr != nil {
				assert.NotNil(t, err)
				assert.Error(t, err)
				assert.Equal(t, c.expectErr, err)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *c.payload.ID, *got.ID)
			}
		})
	}
}

// ---------- UNIT TESTS FOR Delete ----------
func TestAirportRepository_Delete(t *testing.T) {
	deleteQuery := `DELETE airports WHERE id = $1`

	// data dummy
	deletedIdSuccess := dataDummy[0].id // ID yang ada
	deletedIdNotFound := uuid.New()     // ID yang tidak ada
	deletedIdInvalid := "invalid-uuid"  // invalid format

	cases := []struct {
		name        string
		id          string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "existing ID",
			id:   deletedIdSuccess.String(),
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(strings.TrimSpace(deleteQuery))).
					WithArgs(deletedIdSuccess).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "non-existing ID",
			id:   deletedIdNotFound.String(),
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(strings.TrimSpace(deleteQuery))).
					WithArgs(deletedIdNotFound).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedErr: nil,
		},
		{
			name: "invalid UUID",
			id:   deletedIdInvalid,
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(strings.TrimSpace(deleteQuery))).
					WithArgs(deleteQuery).
					WillReturnError(errors.New("invalid uuid"))
			},
			expectedErr: util.ErrBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, ctx, db, mock, cleanup, repoMock := SetupTesting(t)
			defer cleanup()

			// Expectation
			mock.ExpectBegin()

			tc.setupMock(mock)

			// ACT
			tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
			assert.NoError(t, err)
			defer util.FinishTx(tx, &err)

			// Call Delete
			err = repoMock.Delete(ctx, tx, tc.id)

			// ASSERT
			if tc.expectedErr != nil {
				//
			} else {
				assert.Nil(t, err)
			}

			// finalize tx
			assert.NoError(t, tx.Commit())
		})
	}
}
