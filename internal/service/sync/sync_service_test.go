package service_sync

import (
	"context"
	"database/sql"
	airport_dto "flight-api/internal/dto/airport"
	sync_dto "flight-api/internal/dto/sync"
	"flight-api/internal/model"
	repository_airport "flight-api/internal/repository/airport"
	service_aviation "flight-api/internal/service/aviation"
	"flight-api/pkg/logger"
	"flight-api/util"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type assertErr string

func (e assertErr) Error() string { return string(e) }

func TestSyncAirports_Mixed_OneInserted_OneSkipped(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	validate := util.NewValidator()
	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	avi := &service_aviation.AviationServiceMock{Mock: mock.Mock{}}
	svc := NewSyncService(logger, validate, db, repo, avi)

	req := sync_dto.SyncAirportRequest{ICAOCodes: []string{"KJFK", "KSEA"}}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// JFK sudah ada, SEA belum
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KJFK").Return(true, nil).Once()
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KSEA").Return(false, nil).Once()

	// aviation mengembalikan data untuk KSEA
	seaReq := airport_dto.AirportRequestDto{
		ICAOID:    util.Ptr("KSEA"),
		Name:      util.Ptr("Seattle-Tacoma International Airport"),
		City:      util.Ptr("Seattle"),
		Country:   util.Ptr("USA"),
		Latitude:  util.Ptr("47.4502"),
		Longitude: util.Ptr("-122.3088"),
	}
	avi.Mock.On("FetchAirportData", mock.Anything, []string{"KSEA"}).
		Return(map[string]airport_dto.AirportRequestDto{"KSEA": seaReq}, nil).
		Once()

	// insert untuk KSEA
	newIdSEA := uuid.New()
	timeNow := time.Now()
	seaModel := model.Airport{
		ID:        &newIdSEA,
		ICAOID:    seaReq.ICAOID,
		Name:      seaReq.Name,
		City:      seaReq.City,
		Country:   seaReq.Country,
		Latitude:  seaReq.Latitude,
		Longitude: seaReq.Longitude,
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
	}

	repo.Mock.On("Insert",
		mock.Anything,
		mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
		mock.MatchedBy(func(m model.Airport) bool {
			return *m.ICAOID == "KSEA"
		}),
	).Return(seaModel, nil).Once()

	out, err := svc.SyncAirports(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, out, 2)

	// Urutan hasil: sesuai loop ICAO di req
	require.Equal(t, "KJFK", out[0].ICAOCode)
	require.Equal(t, "Skipped", out[0].Status)
	require.Nil(t, out[0].Airport)

	require.Equal(t, "KSEA", out[1].ICAOCode)
	require.Equal(t, "Inserted", out[1].Status)
	require.NotNil(t, out[1].Airport)
	require.Equal(t, "KSEA", *out[1].Airport.ICAOID)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	avi.Mock.AssertExpectations(t)
}

func TestSyncAirports_NotFoundFromAviation(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	validate := util.NewValidator()
	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	avi := &service_aviation.AviationServiceMock{Mock: mock.Mock{}}
	svc := NewSyncService(logger, validate, db, repo, avi)

	req := sync_dto.SyncAirportRequest{ICAOCodes: []string{"KXXX"}}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// Tidak ada di DB
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KXXX").Return(false, nil).Once()

	// Aviation tidak mengembalikan data (key tidak ada)
	avi.Mock.On("FetchAirportData", mock.Anything, []string{"KXXX"}).
		Return(map[string]airport_dto.AirportRequestDto{}, nil).
		Once()

	out, err := svc.SyncAirports(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "KXXX", out[0].ICAOCode)
	require.Equal(t, "Not Found", out[0].Status)
	require.Nil(t, out[0].Airport)

	// Tidak ada upsert
	repo.Mock.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything, mock.Anything)
	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	avi.Mock.AssertExpectations(t)
}

func TestSyncAirports_ErrorOnExistsCheck_ReturnsError(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	validate := util.NewValidator()
	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	avi := &service_aviation.AviationServiceMock{Mock: mock.Mock{}}
	svc := NewSyncService(logger, validate, db, repo, avi)

	req := sync_dto.SyncAirportRequest{ICAOCodes: []string{"KJFK"}}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit() // tidak panic → commit

	repo.Mock.On("FindExistsByICAOID",
		mock.Anything,
		mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
		"KJFK",
	).Return(false, assertErr("db failure")).Once()

	out, err := svc.SyncAirports(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, out)

	// Aviation & Insert tidak terpanggil
	avi.Mock.AssertNotCalled(t, "FetchAirportData", mock.Anything, mock.Anything)
	repo.Mock.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything, mock.Anything)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
}

func TestSyncAirports_ErrorOnAviationFetch_ReturnsError(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	validate := util.NewValidator()
	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	avi := &service_aviation.AviationServiceMock{Mock: mock.Mock{}}
	svc := NewSyncService(logger, validate, db, repo, avi)
	req := sync_dto.SyncAirportRequest{ICAOCodes: []string{"KSEA", "KPDX"}}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// KSEA belum ada, KPDX sudah ada → only fetch KSEA
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KSEA").Return(false, nil).Once()
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KPDX").Return(true, nil).Once()

	avi.Mock.On("FetchAirportData", mock.Anything, []string{"KSEA"}).
		Return(nil, assertErr("aviation timeout")).
		Once()

	out, err := svc.SyncAirports(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, out)

	repo.Mock.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything, mock.Anything)
	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	avi.Mock.AssertExpectations(t)
}

func TestSyncAirports_ErrorOnInsert_ReturnsError(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	validate := util.NewValidator()
	repo := &repository_airport.AirportRepositoryMock{Mock: mock.Mock{}}
	avi := &service_aviation.AviationServiceMock{Mock: mock.Mock{}}
	svc := NewSyncService(logger, validate, db, repo, avi)

	req := sync_dto.SyncAirportRequest{ICAOCodes: []string{"KLAX"}}

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	// KLAX belum ada
	repo.Mock.On("FindExistsByICAOID", mock.Anything, mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }), "KLAX").Return(false, nil).Once()

	// Aviation balikin data KLAX
	laxReq := airport_dto.AirportRequestDto{
		ICAOID:    util.Ptr("KLAX"),
		Name:      util.Ptr("Los Angeles International Airport"),
		City:      util.Ptr("Los Angeles"),
		Country:   util.Ptr("USA"),
		Latitude:  util.Ptr("33.9416"),
		Longitude: util.Ptr("-118.4085"),
	}
	avi.Mock.On("FetchAirportData", mock.Anything, []string{"KLAX"}).
		Return(map[string]airport_dto.AirportRequestDto{"KLAX": laxReq}, nil).
		Once()

	// Insert error
	repo.Mock.On("Insert",
		mock.Anything,
		mock.MatchedBy(func(tx *sql.Tx) bool { return tx != nil }),
		mock.MatchedBy(func(m model.Airport) bool { return *m.ICAOID == "KLAX" }),
	).Return(model.Airport{}, assertErr("insert failed")).Once()

	out, err := svc.SyncAirports(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, out)

	require.NoError(t, dbmock.ExpectationsWereMet())
	repo.Mock.AssertExpectations(t)
	avi.Mock.AssertExpectations(t)
}
