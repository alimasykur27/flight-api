package service_airport

import (
	"context"
	"database/sql"
	airport_dto "flight-api/internal/dto/airport"
	dto "flight-api/internal/dto/airport"
	location_dto "flight-api/internal/dto/location"
	queryparams "flight-api/internal/dto/query_params"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/enum"
	"flight-api/internal/model"
	repository_airport "flight-api/internal/repository/airport"
	service_weather "flight-api/internal/service/weather"
	"flight-api/pkg/logger"
	"flight-api/util"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var timeNow = time.Now()
var sliceId = map[string]uuid.UUID{
	"KJFK": uuid.New(),
	"KLAX": uuid.New(),
	"KSFO": uuid.New(),
}

var dataDummy = []struct {
	id    uuid.UUID
	row   model.Airport
	label string
}{
	{
		id:    sliceId["KJFK"],
		label: "KJFK",
		row: func() model.Airport {
			id := sliceId["KJFK"]
			return model.Airport{
				ID:            &id,
				SiteNumber:    util.Ptr("12345"),
				ICAOID:        util.Ptr("KJFK"),
				FAAID:         util.Ptr("FAA2"),
				IATAID:        util.Ptr("JFK"),
				Name:          util.Ptr("John F. Kennedy Intl"),
				Type:          enum.AIRPORT,
				Status:        util.Ptr(true),
				Country:       util.Ptr("US"),
				State:         util.Ptr("NY"),
				StateFull:     util.Ptr("New York"),
				County:        util.Ptr("Queens"),
				City:          util.Ptr("New York"),
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
		id:    sliceId["KLAX"],
		label: "KLAX",
		row: func() model.Airport {
			id := sliceId["KLAX"]
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
		id:    sliceId["KSFO"],
		label: "KSFO",
		row: func() model.Airport {
			id := sliceId["KSFO"]
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

var Object = "weather"
var dataDummyWeather = map[string]weather_dto.WeatherDto{
	"new_york": weather_dto.WeatherDto{
		Location: &location_dto.LocationDto{
			Name:           util.Ptr("New York"),
			Region:         util.Ptr("New York"),
			Country:        util.Ptr("United States of America"),
			Lat:            util.Ptr(40.7142),
			Lon:            util.Ptr(-74.0064),
			TzId:           util.Ptr("America/New_York"),
			LocaltimeEpoch: util.Ptr(1759794491),
			Localtime:      util.Ptr("2025-10-06 19:48"),
		},
		Object: &Object,
		Current: &weather_dto.CurrentWeatherDto{
			LastUpdatedEpoch: util.Ptr(1759794300),
			LastUpdated:      util.Ptr("2025-10-06 19:45"),
			TempC:            util.Ptr(21.1),
			TempF:            util.Ptr(70.0),
			IsDay:            util.Ptr(uint8(0)),
		},
	},
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

// -------- SERVICE CONSTRUCTOR --------

func newDeps(t *testing.T) (*logger.Logger, *validator.Validate, *sql.DB, sqlmock.Sqlmock, *repository_airport.AirportRepositoryMock, *service_weather.WeatherServiceMock, IAirportService) {
	t.Helper()

	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)

	repoMock := &repository_airport.AirportRepositoryMock{}
	wMock := &service_weather.WeatherServiceMock{}

	svc := NewAirportService(log, val, db, repoMock, wMock)

	return log, val, db, dbmock, repoMock, wMock, svc
}

func TestAirportService_New(t *testing.T) {
	_, _, db, _, repoMock, wMock, svc := newDeps(t)
	defer db.Close()

	assert.NotNil(t, svc)
	// tipe konkret
	_, ok := svc.(*AirportService)
	assert.True(t, ok)
	assert.NotNil(t, repoMock)
	assert.NotNil(t, wMock)
}

func TestAirportService_Create_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}

	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	// Expect tx dari service
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// Arrange request
	req := airport_dto.AirportRequestDto{
		ICAOID:    util.Ptr("KJFK"),
		Name:      util.Ptr("John F. Kennedy International Airport"),
		City:      util.Ptr("New York"),
		Country:   util.Ptr("USA"),
		Latitude:  util.Ptr("40.6413"),
		Longitude: util.Ptr("-73.7781"),
	}

	// Arrange output model yg diharapkan dari repo.Insert
	newId := uuid.New()
	timeNow := time.Now()
	expectedModel := model.Airport{
		ID:            &newId,
		SiteNumber:    nil,
		ICAOID:        util.Ptr("KJFK"),
		FAAID:         nil,
		IATAID:        nil,
		Name:          util.Ptr("John F. Kennedy International Airport"),
		Type:          nil,
		Status:        nil,
		Country:       util.Ptr("USA"),
		State:         nil,
		StateFull:     nil,
		County:        nil,
		City:          util.Ptr("New York"),
		Ownership:     nil,
		Use:           nil,
		Manager:       nil,
		ManagerPhone:  nil,
		Latitude:      util.Ptr("40.6413"),
		LatitudeSec:   nil,
		Longitude:     util.Ptr("-73.7781"),
		LongitudeSec:  nil,
		Elevation:     nil,
		ControlTower:  nil,
		Unicom:        nil,
		CTAF:          nil,
		EffectiveDate: nil,
		CreatedAt:     &timeNow,
		UpdatedAt:     &timeNow,
	}

	// Expect: repo.Insert dipanggil dengan ctx apapun, tx valid, dan airport model yang terbentuk dari request
	repoMock.Mock.
		On(
			"Insert",
			mock.Anything, // ctx
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			mock.AnythingOfType("model.Airport"),
		).
		Return(expectedModel, nil).
		Once()

	// Act
	ctx := context.Background()
	out := svc.Create(ctx, req)

	assert.NotNil(t, out)
	assert.NotNil(t, out.ID)
	assert.Equal(t, *req.ICAOID, *out.ICAOID)
	assert.Equal(t, *req.Name, *out.Name)
	assert.Equal(t, *req.City, *out.City)
	assert.Equal(t, *req.Country, *out.Country)
	assert.Equal(t, *req.Latitude, *out.Latitude)
	assert.Equal(t, *req.Longitude, *out.Longitude)

	// Pastikan ekspektasi mock terpenuhi
	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_FindAll_Success_WithNextTrue(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}

	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	// Query params: Limit kecil, total besar → Next = true
	q := queryparams.QueryParams{
		Limit:  2,
		Offset: 0,
		Page:   1,
	}

	// Expect transaksi
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// Data dari repo
	a1 := dataDummy[0].row
	a2 := dataDummy[1].row
	list := []model.Airport{a1, a2}
	total := 5 // (offset + limit) < total ⇒ (0 + 2) < 5 ⇒ Next: true

	// Expect panggilan repo dengan tx != nil dan args limit/offset sesuai
	repoMock.Mock.
		On("FindAll",
			mock.Anything, // ctx
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			mock.MatchedBy(func(m map[string]interface{}) bool {
				lim, lok := m["limit"].(int)
				off, ook := m["offset"].(int)
				return lok && ook && lim == q.Limit && off == q.Offset
			}),
		).
		Return(list, total, nil).
		Once()

	// Act
	out := svc.FindAll(context.Background(), q)

	// Assert pagination dto
	require.Equal(t, "pagination", out.Object)
	require.Equal(t, total, out.Total)
	require.NotNil(t, out.Meta)
	require.Equal(t, q.Limit, out.Meta.Limit)
	require.Equal(t, q.Page, out.Meta.Page)
	require.True(t, out.Meta.Next)

	// Records length
	require.Len(t, out.Records, len(list))

	// Record pertama
	rec1 := out.Records[0].(airport_dto.AirportRecordDto)
	require.Equal(t, *a1.ID, *rec1.ID)
	require.Equal(t, *a1.ICAOID, *rec1.ICAOID)
	require.Equal(t, *a1.Name, *rec1.Name)

	// Record kedua
	rec2 := out.Records[1].(airport_dto.AirportRecordDto)
	require.Equal(t, *a2.ID, *rec2.ID)
	require.Equal(t, *a2.ICAOID, *rec2.ICAOID)
	require.Equal(t, *a2.Name, *rec2.Name)

	// Pastikan ekspektasi mock terpenuhi
	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_FindAll_Success_NoNext(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}

	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	// (offset + limit) == total ⇒ Next: false
	q := queryparams.QueryParams{
		Limit:  2,
		Offset: 4,
		Page:   3,
	}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	a1 := dataDummy[0].row
	a2 := dataDummy[1].row
	list := []model.Airport{a1, a2}
	total := 6 // 4+2 == 6 ⇒ Next false

	repoMock.Mock.
		On("FindAll",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			mock.MatchedBy(func(m map[string]interface{}) bool {
				return m["limit"] == q.Limit && m["offset"] == q.Offset
			}),
		).
		Return(list, total, nil).
		Once()

	out := svc.FindAll(context.Background(), q)

	require.Equal(t, "pagination", out.Object)
	require.Equal(t, total, out.Total)
	require.NotNil(t, out.Meta)
	require.Equal(t, q.Limit, out.Meta.Limit)
	require.Equal(t, q.Page, out.Meta.Page)
	require.False(t, out.Meta.Next)
	require.Len(t, out.Records, len(list))

	// Record pertama
	rec1 := out.Records[0].(airport_dto.AirportRecordDto)
	require.Equal(t, *a1.ID, *rec1.ID)
	require.Equal(t, *a1.ICAOID, *rec1.ICAOID)
	require.Equal(t, *a1.Name, *rec1.Name)

	// Record kedua
	rec2 := out.Records[1].(airport_dto.AirportRecordDto)
	require.Equal(t, *a2.ID, *rec2.ID)
	require.Equal(t, *a2.ICAOID, *rec2.ICAOID)
	require.Equal(t, *a2.Name, *rec2.Name)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_FindAll_RepoError_Rollback(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}

	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	q := queryparams.QueryParams{Limit: 10, Offset: 0, Page: 1}

	dbmock.ExpectBegin()
	dbmock.ExpectRollback() // karena error

	repoMock.Mock.
		On("FindAll",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			mock.AnythingOfType("map[string]interface {}"),
		).
		Return(nil, 0, assertErr("db failure")).
		Once()

	require.Panics(t, func() {
		_ = svc.FindAll(context.Background(), q)
	})

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_FindByID_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	targetID := sliceId["KJFK"]
	expectedModel := dataDummy[0].row
	expectedDto := dto.ToAirportDto(expectedModel)

	// transaksi: begin + commit (tidak panic)
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// expect repo: tx valid & id sesuai
	repoMock.Mock.
		On("FindByID",
			mock.Anything, // ctx
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			targetID.String(),
		).
		Return(expectedModel, nil).
		Once()

	// act
	got, err := svc.FindByID(context.Background(), targetID.String())

	// assert
	require.NoError(t, err)
	require.Equal(t, expectedDto, got)

	assert.NotNil(t, got.ID)
	assert.Equal(t, *expectedModel.ICAOID, *got.ICAOID)
	assert.Equal(t, *expectedModel.Name, *got.Name)
	assert.Equal(t, *expectedModel.City, *got.City)
	assert.Equal(t, *expectedModel.Country, *got.Country)
	assert.Equal(t, *expectedModel.Latitude, *got.Latitude)
	assert.Equal(t, *expectedModel.Longitude, *got.Longitude)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_FindByID_NotFound(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	unknownID := uuid.New().String()

	// transaksi tetap commit karena tidak ada panic
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	repoMock.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			unknownID,
		).
		Return(model.Airport{}, util.ErrNotFound).
		Once()

	// act
	got, err := svc.FindByID(context.Background(), unknownID)

	// assert
	require.Error(t, err)
	require.Equal(t, util.ErrNotFound, err)
	require.Equal(t, dto.AirportDto{}, got) // kosong sesuai code

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

// -------- UPDATE --------
func TestAirportService_Update_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	// Arrange
	id := sliceId["KJFK"]
	existing := dataDummy[0].row // KJFK

	// update fields
	newName := "JFK Intl (Renamed)"
	newCity := "NYC"
	updatedTime := time.Now().Add(1 * time.Hour) // pastikan UpdatedAt berubah
	u := airport_dto.AirportUpdateDto{
		Name: util.Ptr(newName),
		City: util.Ptr(newCity),
	}

	// transaksi: begin + commit
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// repo.FindByID -> return existing
	repoMock.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id.String(),
		).
		Return(existing, nil).
		Once()

	// expect Update menerima airport yg sudah terisi perubahan
	updatedAirport := existing
	*updatedAirport.Name = newName
	*updatedAirport.City = newCity
	updatedAirport.UpdatedAt = &updatedTime

	// repo.Update -> return updatedAirport
	repoMock.Mock.
		On("Update",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id.String(),
			mock.AnythingOfType("model.Airport"),
		).
		Return(updatedAirport, nil).
		Once()

	// Act
	got, err := svc.Update(context.Background(), id.String(), u)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got.ID)
	require.Equal(t, newName, *got.Name)
	require.Equal(t, newCity, *got.City)
	require.Equal(t, *existing.ICAOID, *got.ICAOID)
	require.Equal(t, *existing.Country, *got.Country)
	require.Equal(t, *existing.Latitude, *got.Latitude)
	require.Equal(t, *existing.Longitude, *got.Longitude)
	require.True(t, got.UpdatedAt.After(*existing.UpdatedAt))

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_Update_FindByID_NotFound(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	id := uuid.New().String()

	dbmock.ExpectBegin()
	dbmock.ExpectCommit() // tidak panic → commit

	repoMock.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id,
		).
		Return(model.Airport{}, util.ErrNotFound).
		Once()

	got, err := svc.Update(context.Background(), id, airport_dto.AirportUpdateDto{City: util.Ptr("X")})

	require.Error(t, err)
	require.Equal(t, util.ErrNotFound, err)
	require.Equal(t, dto.AirportDto{}, got)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

func TestAirportService_Update_UpdateRepoError(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repoMock := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weatherMock := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repoMock, weatherMock)

	id := sliceId["KLAX"]
	existing := dataDummy[1].row // KLAX
	airport_dto := airport_dto.AirportUpdateDto{
		Name: util.Ptr("SEA Intl"),
		City: util.Ptr("Seattle"),
	}

	dbmock.ExpectBegin()
	dbmock.ExpectRollback() // karena PanicIfError pada Update

	// FindByID OK
	repoMock.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id.String(),
		).
		Return(existing, nil).
		Once()

	// Update error -> panic
	repoMock.Mock.
		On("Update",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id.String(),
			mock.AnythingOfType("model.Airport"),
		).
		Return(model.Airport{}, assertErr("update failed")).
		Once()

	require.Panics(t, func() {
		_, _ = svc.Update(context.Background(), id.String(), airport_dto)
	})

	require.NoError(t, dbmock.ExpectationsWereMet())
	repoMock.Mock.AssertExpectations(t)
}

// -------- DELETE --------
func TestAirportService_Delete_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	existingID := sliceId["KSFO"]
	existing := dataDummy[2].row // KSFO
	id := existingID.String()

	// transaksi: Begin + Commit
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// FindByID OK → return data
	repo.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id,
		).
		Return(existing, nil).
		Once()

	// Delete OK
	repo.Mock.
		On("Delete",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id,
		).
		Return(nil).
		Once()

	// act
	err = svc.Delete(context.Background(), id)

	// assert
	require.NoError(t, err)
	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
}

func TestAirportService_Delete_NotFound(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	id := uuid.New().String()

	// Begin + Commit (tidak panic)
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	repo.Mock.
		On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			id,
		).
		Return(model.Airport{}, util.ErrNotFound).
		Once()

	// Delete tidak boleh dipanggil
	// act
	err = svc.Delete(context.Background(), id)

	// assert
	require.Error(t, err)
	require.Equal(t, util.ErrNotFound, err)
	repo.Mock.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
}

// -------_ getWeatherConditionByCode --------
func TestGetWeatherCondition_EmptyCodeAndName(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	queryParam := queryparams.QueryParams{
		Limit:  10,
		Offset: 0,
		Page:   1,
	}
	out, err := svc.GetWeatherCondition(context.Background(), "", "", queryParam)
	require.Error(t, err)
	require.Nil(t, out)
	require.Equal(t, util.ErrBadRequest, err)

	// Pastikan tidak ada panggilan repo/weather
	repo.Mock.AssertNotCalled(t, "FindByICAOID", mock.Anything, mock.Anything, mock.Anything)
	weather.Mock.AssertNotCalled(t, "GetWeatherCondition", mock.Anything, mock.Anything)
}

func TestGetWeatherCondition_ByCode_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	code := "KJFK"
	airport := dataDummy[0].row // KJFK

	// transaksi: Begin + Commit
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// repo FindByICAOID OK
	repo.Mock.
		On("FindByICAOID",
			mock.Anything, // ctx
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			code,
		).
		Return(airport, nil).
		Once()

	// weather.GetWeatherCondition OK
	var weatherData = dataDummyWeather["new_york"]
	weatherResp := &weatherData
	weather.Mock.
		On("GetWeatherCondition",
			mock.Anything,
			airport.City, // *string
		).
		Return(weatherResp, nil).
		Once()

	// act
	queryParam := queryparams.QueryParams{
		Limit:  10,
		Offset: 0,
		Page:   1,
	}
	out, err := svc.GetWeatherCondition(context.Background(), code, "", queryParam)
	require.NoError(t, err)
	require.NotNil(t, out)

	require.Equal(t, "pagination", out.Object)
	require.Equal(t, 1, out.Total)
	require.Len(t, out.Records, 1)

	// Rekor pertama harus AirportWeatherDto
	rec0, ok := out.Records[0].(dto.AirportWeatherDto)
	require.True(t, ok, "record[0] bukan dto.AirportWeatherDto")
	require.Equal(t, "airport_weather", rec0.Object)
	require.NotNil(t, rec0.Airport)
	require.Equal(t, "KJFK", *rec0.Airport.ICAOID)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	weather.Mock.AssertExpectations(t)
}

func TestGetWeatherCondition_ByCode_NotFound(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	code := "XXXX"

	// Begin + Commit (tidak ada panic)
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	repo.Mock.
		On("FindByICAOID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			code,
		).
		Return(model.Airport{}, util.ErrNotFound).
		Once()

	queryParam := queryparams.QueryParams{
		Limit:  10,
		Offset: 0,
		Page:   1,
	}
	out, err := svc.GetWeatherCondition(context.Background(), code, "", queryParam)
	require.Error(t, err)
	require.Nil(t, out)
	require.Equal(t, util.ErrNotFound, err)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	weather.Mock.AssertNotCalled(t, "GetWeatherCondition", mock.Anything, mock.Anything)
}

func TestGetWeatherCondition_ByCode_UnexpectedRepoError(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	code := "KERR"

	// Begin + Commit (servicemu return error biasa, bukan panic)
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	repo.Mock.
		On("FindByICAOID",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			code,
		).
		Return(model.Airport{}, assertErr("db down")).
		Once()

	queryParam := queryparams.QueryParams{
		Limit:  10,
		Offset: 0,
		Page:   1,
	}
	out, err := svc.GetWeatherCondition(context.Background(), code, "", queryParam)
	require.Error(t, err)
	require.Nil(t, out)
	require.Equal(t, util.ErrInternalServer, err)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	weather.Mock.AssertNotCalled(t, "GetWeatherCondition", mock.Anything, mock.Anything)
}

func TestGetWeatherCondition_BySearchName_Success(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	name := "International"
	q := queryparams.QueryParams{Limit: 2, Offset: 0, Page: 1}

	// transaksi: Begin + Commit
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// dummy airports
	a1 := dataDummy[0].row // KJFK
	a2 := dataDummy[1].row // KLAX
	list := []model.Airport{a1, a2}
	total := 2

	// repo.FindBySearchName dipanggil dengan name & args (limit/offset) sesuai
	repo.Mock.
		On("FindBySearchName",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			name,
			mock.MatchedBy(func(m map[string]interface{}) bool {
				lim, ok1 := m["limit"].(int)
				off, ok2 := m["offset"].(int)
				return ok1 && ok2 && lim == q.Limit && off == q.Offset
			}),
		).
		Return(list, total, nil).
		Once()

	// weather service dipanggil per-airport (city pointer)
	weatherData := dataDummyWeather["new_york"]
	var weather1 *weather_dto.WeatherDto = &weatherData
	weather.Mock.
		On("GetWeatherCondition", mock.Anything, a1.City).
		Return(weather1, nil).
		Once()

	var weather2 *weather_dto.WeatherDto = nil
	weather.Mock.
		On("GetWeatherCondition", mock.Anything, a2.City).
		Return(weather2, nil).
		Once()

	// act via public method yg rutenya ke getWeatherConditionBySearchName
	out, err := svc.GetWeatherCondition(context.Background(), "", name, q)

	// assert
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, "pagination", out.Object)
	require.Equal(t, total, out.Total)
	require.NotNil(t, out.Meta)
	require.Equal(t, q.Limit, out.Meta.Limit)
	require.Equal(t, q.Page, out.Meta.Page)
	require.False(t, out.Meta.Next)

	require.Len(t, out.Records, 2)

	// record pertama
	rec1, ok := out.Records[0].(dto.AirportWeatherDto)
	require.True(t, ok, "record[0] bukan dto.AirportWeatherDto")
	require.Equal(t, "airport_weather", rec1.Object)
	require.NotNil(t, rec1.Airport)
	require.Equal(t, "KJFK", *rec1.Airport.ICAOID)
	require.NotNil(t, rec1.Weather)
	require.Equal(t, *weather1.Current.TempC, *rec1.Weather.TempC)

	// record kedua
	rec2, ok := out.Records[1].(dto.AirportWeatherDto)
	require.True(t, ok, "record[1] bukan dto.AirportWeatherDto")
	require.Equal(t, "airport_weather", rec2.Object)
	require.NotNil(t, rec2.Airport)
	require.Equal(t, "KLAX", *rec2.Airport.ICAOID)
	require.Nil(t, rec2.Weather) // karena weather service return nil

	// Pastikan ekspektasi mock terpenuhi
	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	weather.Mock.AssertExpectations(t)
}

func TestGetWeatherCondition_BySearchName_RepoError(t *testing.T) {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	val := util.NewValidator()
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	weather := &service_weather.WeatherServiceMock{Mock: mock.Mock{}}
	svc := NewAirportService(log, val, db, repo, weather)

	name := "X"
	q := queryparams.QueryParams{Limit: 5, Offset: 10, Page: 4}

	// Begin + Commit (return error biasa, tidak panic)
	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	repo.Mock.
		On("FindBySearchName",
			mock.Anything,
			mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
			name,
			mock.MatchedBy(func(m map[string]interface{}) bool {
				return m["limit"] == q.Limit && m["offset"] == q.Offset
			}),
		).
		Return(nil, 0, assertErr("db error")).
		Once()

	out, err := svc.GetWeatherCondition(context.Background(), "", name, q)
	require.Error(t, err)
	require.Nil(t, out)

	// Weather tidak boleh terpanggil
	weather.Mock.AssertNotCalled(t, "GetWeatherCondition", mock.Anything, mock.Anything)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
}
