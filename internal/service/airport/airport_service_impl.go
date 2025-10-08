package service_airport

import (
	"context"
	"database/sql"
	airport_dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	weather_dto "flight-api/internal/dto/weather"
	repository_airport "flight-api/internal/repository/airport"
	service_weather "flight-api/internal/service/weather"

	"flight-api/pkg/logger"
	"flight-api/util"

	"github.com/go-playground/validator"
)

type AirportService struct {
	logger            *logger.Logger
	validate          *validator.Validate
	db                *sql.DB
	airportRepository repository_airport.IAirportRepository
	weatherService    service_weather.IWeatherService
}

func NewAirportService(
	logger *logger.Logger,
	validate *validator.Validate,
	db *sql.DB,
	airportRepository repository_airport.IAirportRepository,
	weatherService service_weather.IWeatherService,
) IAirportService {
	return &AirportService{
		logger:            logger,
		validate:          validate,
		db:                db,
		airportRepository: airportRepository,
		weatherService:    weatherService,
	}
}

func (s *AirportService) Seeding(ctx context.Context, reqs []string) ([]airport_dto.AirportDto, error) {
	s.logger.Debug("[Seeding] Seeding airport data...")

	return nil, nil
}

func (s *AirportService) Create(ctx context.Context, r airport_dto.AirportRequestDto) (airport_dto.AirportDto, error) {
	s.logger.Debug("[Create] Creating new airport...")

	err := s.validate.Struct(r)
	if err != nil {
		return airport_dto.AirportDto{}, nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Errorf("[Create] Failed to begin transaction: %v", err)
		return airport_dto.AirportDto{}, util.ErrInternalServer
	}
	defer util.CommitOrRollback(tx)

	isExists, err := s.airportRepository.FindExistsByICAOID(ctx, tx, *r.ICAOID)
	if err != nil {
		s.logger.Errorf("[Create] Failed to check existing airport: %v", err)
		return airport_dto.AirportDto{}, util.ErrInternalServer
	}

	if isExists {
		s.logger.Warnf("[Create] Airport with ICAO ID %s already exists", *r.ICAOID)
		return airport_dto.AirportDto{}, util.ErrConflict
	}

	airport := airport_dto.AirportRequestToAirport(r)
	airport, err = s.airportRepository.Insert(ctx, tx, airport)
	if err != nil {
		s.logger.Errorf("[Create] Failed to insert airport: %v", err)
		return airport_dto.AirportDto{}, err
	}

	data := airport_dto.ToAirportDto(airport)
	return data, nil
}

func (s *AirportService) FindAll(ctx context.Context, query queryparams.QueryParams) (pagination_dto.PaginationDto, error) {
	s.logger.Debug("[FindAll] Fetching all airports...")

	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Errorf("[FindAll] Failed to begin transaction: %v", err)
		return pagination_dto.PaginationDto{}, util.ErrInternalServer
	}
	defer util.CommitOrRollback(tx)

	args := map[string]interface{}{
		"limit":  query.Limit,
		"offset": query.Offset,
	}
	airports, total, err := s.airportRepository.FindAll(ctx, tx, args)
	if err != nil {
		s.logger.Errorf("[FindAll] Failed to fetch airports: %v", err)
		return pagination_dto.PaginationDto{}, util.ErrInternalServer
	}

	airportRecords := airport_dto.ToAirportRecordDtos(airports)
	records := util.ToInterfaces(airportRecords)
	hasNext := (query.Offset + query.Limit) < total

	response := pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: records,
		Total:   total,
		Meta: &pagination_dto.PaginationMetaDto{
			Limit: query.Limit,
			Page:  query.Page,
			Next:  hasNext,
		},
	}

	return response, nil
}

func (s *AirportService) FindByID(ctx context.Context, id string) (airport_dto.AirportDto, error) {
	s.logger.Debug("[FindByID] Fetching airport by ID...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport, err := s.airportRepository.FindByID(ctx, tx, id)

	if err != nil {
		return airport_dto.AirportDto{}, util.ErrNotFound
	}

	return airport_dto.ToAirportDto(airport), nil
}

func (s *AirportService) Update(ctx context.Context, id string, u airport_dto.AirportUpdateDto) (airport_dto.AirportDto, error) {
	s.logger.Debug("[Update] Updating airport...")

	err := s.validate.Struct(u)
	if err != nil {
		util.LogPanicError(err)
		return airport_dto.AirportDto{}, util.ErrBadRequest
	}

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport, err := s.airportRepository.FindByID(ctx, tx, id)

	if err == util.ErrNotFound {
		return airport_dto.AirportDto{}, util.ErrNotFound
	} else if err != nil {
		util.PanicIfError(err)
	}

	util.FillUpdatableFields(&airport, u)
	updatedAirport, err := s.airportRepository.Update(ctx, tx, id, airport)
	util.PanicIfError(err)

	s.logger.Debugf("[Update] Airport updated: %+v", updatedAirport)
	return airport_dto.ToAirportDto(updatedAirport), nil
}

func (s *AirportService) Delete(ctx context.Context, id string) error {
	s.logger.Debug("[Delete] Deleting airport...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	_, err = s.airportRepository.FindByID(ctx, tx, id)

	if err == util.ErrNotFound {
		return util.ErrNotFound
	}

	err = s.airportRepository.Delete(ctx, tx, id)
	util.PanicIfError(err)

	return nil
}

func (s *AirportService) GetWeatherCondition(ctx context.Context, code string, name string, query queryparams.QueryParams) (*pagination_dto.PaginationDto, error) {
	s.logger.Debugf("[GetWeatherCondition] Fetching weather data from Weather APIs...")

	var response *pagination_dto.PaginationDto
	var err error

	if code != "" {
		response, err = s.getWeatherConditionByCode(ctx, code)
	} else if name != "" {
		response, err = s.getWeatherConditionBySearchName(ctx, name, query)
	} else {
		return nil, util.ErrBadRequest
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *AirportService) getWeatherConditionByCode(ctx context.Context, code string) (*pagination_dto.PaginationDto, error) {
	s.logger.Debugf("[getWeatherConditionByCode] Fetching weather data from Weather APIs...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	// Find Airport By ICAO ID
	airport, err := s.airportRepository.FindByICAOID(ctx, tx, code)

	if err == util.ErrNotFound {
		return nil, util.ErrNotFound
	} else if err != nil {
		return nil, util.ErrInternalServer
	}

	// Get Airport Weather Condition
	weather, _ := s.weatherService.GetWeatherCondition(ctx, airport.City)

	data := []airport_dto.AirportWeatherDto{}
	airportWeather := airport_dto.AirportWeatherDto{
		Object:  "airport_weather",
		Code:    airport.ICAOID,
		Airport: util.Ptr(airport_dto.ToAirportDto(airport)),
		Weather: weather.Current,
	}
	data = append(data, airportWeather)

	response := pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: util.ToInterfaces(data),
		Total:   len(data),
		Meta:    nil,
	}

	return &response, nil
}

func (s *AirportService) getWeatherConditionBySearchName(ctx context.Context, name string, query queryparams.QueryParams) (*pagination_dto.PaginationDto, error) {
	s.logger.Debugf("[getWeatherConditionBySearchName] Fetching weather data from Weather APIs...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	args := map[string]interface{}{
		"limit":  query.Limit,
		"offset": query.Offset,
	}

	// Get Airport By Search Name
	airports, total, err := s.airportRepository.FindBySearchName(ctx, tx, name, args)
	if err != nil {
		return nil, err
	}

	records := []airport_dto.AirportWeatherDto{}
	for _, airport := range airports {
		// Get Airport Weather Condition
		weather, _ := s.weatherService.GetWeatherCondition(ctx, airport.City)

		var current *weather_dto.CurrentWeatherDto
		if weather == nil {
			current = nil
		} else {
			current = weather.Current
		}

		res := airport_dto.AirportWeatherDto{
			Object:  "airport_weather",
			Code:    airport.ICAOID,
			Airport: util.Ptr(airport_dto.ToAirportDto(airport)),
			Weather: current,
		}

		records = append(records, res)
	}

	result := pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: util.ToInterfaces(records),
		Total:   total,
		Meta: &pagination_dto.PaginationMetaDto{
			Limit: query.Limit,
			Page:  query.Page,
			Next:  false,
		},
	}

	return &result, nil
}
