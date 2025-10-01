package service_sync

import (
	"context"
	"database/sql"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	sync_dto "flight-api/internal/dto/sync"
	repo_airport "flight-api/internal/repository/airport"
	service_aviation "flight-api/internal/service/aviation"
	"flight-api/pkg/logger"
	"flight-api/util"

	"github.com/go-playground/validator"
)

type SyncService struct {
	logger            *logger.Logger
	validate          *validator.Validate
	cfg               *config.Config
	db                *sql.DB
	airportRepository repo_airport.IAirportRepository
	aviationService   service_aviation.IAviationService
}

func NewSyncService(
	logger *logger.Logger,
	validate *validator.Validate,
	cfg *config.Config,
	db *sql.DB,
	airportRepository repo_airport.IAirportRepository,
	aviationService service_aviation.IAviationService,
) ISyncService {
	return &SyncService{
		logger:            logger,
		validate:          validate,
		cfg:               cfg,
		db:                db,
		airportRepository: airportRepository,
		aviationService:   aviationService,
	}
}

func (s *SyncService) SyncAirports(
	ctx context.Context,
	req sync_dto.SyncAirportRequest,
) ([]sync_dto.SyncAirportResponse, error) {
	err := s.validate.Struct(req)
	if err != nil {
		s.logger.Errorf("[SyncAirports] request validation failed %s", err)
		return nil, err
	}
	s.logger.Debug("[SyncAirports] request validated")

	// Check if ICAO codes are exists on the database.
	// If not, fetch data from Aviation API and store it in the database.
	// If exists, skip fetching data from Aviation API.
	var icaoCodesToFetch []string
	results := make([]sync_dto.SyncAirportResponse, 0)

	s.logger.Debug("[SyncAirports] Checking existing ICAO codes in the database...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	for _, code := range req.ICAOCodes {
		exists, err := s.airportRepository.FindExistsByICAOID(ctx, tx, code)
		if err != nil {
			s.logger.Errorf("[SyncAirports] failed to check if ICAO code %s exists: %v", code, err)
			return nil, err
		}
		if !exists {
			s.logger.Debugf("[SyncAirports] ICAO code %s does not exist in the database.", code)
			icaoCodesToFetch = append(icaoCodesToFetch, code)
		} else {
			result := sync_dto.SyncAirportResponse{
				ICAOCode: code,
				Airport:  nil,
				Status:   "Skipped",
				Message:  "ICAO code already exists in the database. Skipping fetch.",
			}
			results = append(results, result)
			s.logger.Debugf("[SyncAirports] ICAO code %s already exists in the database. Skipping fetch.", code)
		}
	}

	// Fetch data from Aviation API.
	s.logger.Debugf("[SyncAirports] Fetching data for ICAO codes: %v", icaoCodesToFetch)
	airportData, err := s.aviationService.FetchAirportData(ctx, icaoCodesToFetch)
	if err != nil {
		s.logger.Errorf("[SyncAirports] failed to fetch airport data from Aviation API: %v", err)
		return nil, err
	}
	s.logger.Debug("[SyncAirports] successfully fetched airport data from Aviation API")

	// Store fetched data in the database.
	s.logger.Debug("[SyncAirports] Storing fetched airport data in the database...")

	for _, code := range icaoCodesToFetch {
		data, exists := airportData[code]

		// If data is empty struct or nil, and not exists, skip inserting to the database.
		if !exists || (data == airport_dto.AirportRequestDto{}) {
			s.logger.Warnf("[SyncAirports] No data found for ICAO code %s from Aviation API. Skipping...", code)
			result := sync_dto.SyncAirportResponse{
				ICAOCode: code,
				Airport:  nil,
				Status:   "Not Found",
				Message:  "No data found from Aviation API.",
			}
			results = append(results, result)
			continue
		}

		airportPayload := airport_dto.AirportRequestToAirport(data)
		airportModel, err := s.airportRepository.Insert(ctx, tx, airportPayload)

		if err != nil {
			s.logger.Errorf("[SyncAirports] failed to insert airport data for ICAO code %s: %v", code, err)
			return nil, err
		}

		s.logger.Debugf("[SyncAirports] Successfully inserted airport data for ICAO code %s", code)
		airportDto := airport_dto.ToAirportDto(airportModel)
		result := sync_dto.SyncAirportResponse{
			ICAOCode: code,
			Airport:  &airportDto,
			Status:   "Inserted",
			Message:  "Airport data successfully inserted",
		}
		results = append(results, result)

		s.logger.Debugf("[SyncAirports] Airport data for ICAO code %s Inserted", code)
	}

	s.logger.Debugf("[SyncAirports] Successfully synced airport data")
	return results, err
}
