package service_airport

import (
	"context"
	"database/sql"
	"flight-api/config"
	"flight-api/internal/cache"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/model"
	repository_airport "flight-api/internal/repository/airport"
	service_weather "flight-api/internal/service/weather"
	"time"

	"flight-api/pkg/logger"
	"flight-api/util"

	"github.com/go-playground/validator"
)

type AirportService struct {
	logger            *logger.Logger
	cfg               *config.Config
	validate          *validator.Validate
	db                *sql.DB
	cache             cache.ICache
	airportRepository repository_airport.IAirportRepository
	weatherService    service_weather.IWeatherService
}

func NewAirportService(
	logger *logger.Logger,
	cfg *config.Config,
	validate *validator.Validate,
	db *sql.DB,
	cache cache.ICache,
	airportRepository repository_airport.IAirportRepository,
	weatherService service_weather.IWeatherService,
) IAirportService {
	return &AirportService{
		logger:            logger,
		cfg:               cfg,
		validate:          validate,
		db:                db,
		cache:             cache,
		airportRepository: airportRepository,
		weatherService:    weatherService,
	}
}

func (s *AirportService) Create(ctx context.Context, r dto.AirportRequestDto) dto.AirportDto {
	s.logger.Debug("[Create] Creating new airport...")

	err := s.validate.Struct(r)
	util.PanicIfError(err)

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport := dto.AirportRequestToAirport(r)
	airport, err = s.airportRepository.Insert(ctx, tx, airport)
	util.PanicIfError(err)

	data := dto.ToAirportDto(airport)
	return data
}

func (s *AirportService) FindAll(ctx context.Context, query queryparams.QueryParams) pagination_dto.PaginationDto {
	s.logger.Debug("[FindAll] Fetching all airports...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	args := map[string]interface{}{
		"limit":  query.Limit,
		"offset": query.Offset,
	}
	airports, total, err := s.airportRepository.FindAll(ctx, tx, args)
	util.PanicIfError(err)

	airportRecords := dto.ToAirportRecordDtos(airports)
	records := make([]interface{}, len(airportRecords))
	for i, v := range airportRecords {
		records[i] = v
	}
	hasNext := (query.Offset + query.Limit) < total

	return pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: records,
		Total:   total,
		Meta: &pagination_dto.PaginationMetaDto{
			Limit: query.Limit,
			Page:  query.Page,
			Next:  hasNext,
		},
	}
}

func (s *AirportService) FindByID(ctx context.Context, id string) (dto.AirportDto, error) {
	s.logger.Debug("[FindByID] Fetching airport by ID...")

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport, err := s.airportRepository.FindByID(ctx, tx, id)

	if err != nil {
		return dto.AirportDto{}, util.ErrNotFound
	}

	return dto.ToAirportDto(airport), nil
}

func (s *AirportService) Update(ctx context.Context, id string, u dto.AirportUpdateDto) (dto.AirportDto, error) {
	s.logger.Debug("[Update] Updating airport...")

	err := s.validate.Struct(u)
	if err != nil {
		util.LogPanicError(err)
		return dto.AirportDto{}, util.ErrBadRequest
	}

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport, err := s.airportRepository.FindByID(ctx, tx, id)

	if err == util.ErrNotFound {
		return dto.AirportDto{}, util.ErrNotFound
	} else if err != nil {
		util.PanicIfError(err)
	}

	util.FillUpdatableFields(&airport, u)
	updatedAirport, err := s.airportRepository.Update(ctx, tx, id, airport)
	util.PanicIfError(err)

	s.logger.Debugf("[Update] Airport updated: %+v", updatedAirport)
	return dto.ToAirportDto(updatedAirport), nil
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

	airportWeather := []dto.AirportWeatherDto{}
	var airportData model.Airport
	var isExists bool = false

	// // Find Airport On redis first
	if s.cfg.RedisEnable {
		airportCache, err := s.cache.FindAirportByICAOID(ctx, code)

		if err == nil {
			s.logger.Debug("[getWeatherConditionByCode] Airport data cache is found on redis")
			isExists = true
			airportData = *airportCache
		} else {
			s.logger.Warn("[getWeatherConditionByCode] Airport data cache not found.")
			isExists = false
		}
	}

	// Find Airport By ICAO ID
	if !isExists {
		airportData, err := s.airportRepository.FindByICAOID(ctx, s.db, code)

		if err == util.ErrNotFound {
			return nil, util.ErrNotFound
		} else if err != nil {
			return nil, util.ErrInternalServer
		}

		if s.cfg.RedisEnable {
			c := code
			p := airportData

			go func(ctx context.Context, key string, payload model.Airport) {
				if err := s.cache.CacheAirport(ctx, key, &payload, 30*time.Minute); err != nil {
					s.logger.Errorf("[getWeatherConditionByCode] Failed to cache airport data to Redis: %v", err)
				} else {
					s.logger.Debug("[getWeatherConditionByCode] Successfully cached airport data to Redis.")
				}
			}(ctx, c, p)
		}
	}

	// Get Airport Weather Condition
	weather, _ := s.weatherService.GetWeatherCondition(ctx, airportData.City)

	currentWeatther := &weather_dto.CurrentWeatherDto{}
	if weather == nil {
		currentWeatther = nil
	} else {
		currentWeatther = weather.Current
	}

	data := dto.AirportWeatherDto{
		Object:  "airport_weather",
		Code:    airportData.ICAOID,
		Airport: util.Ptr(dto.ToAirportDto(airportData)),
		Weather: currentWeatther,
	}
	airportWeather = append(airportWeather, data)

	response := pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: util.ToInterfaces(airportWeather),
		Total:   len(airportWeather),
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

	records := []dto.AirportWeatherDto{}
	for _, airport := range airports {
		// Get Airport Weather Condition
		weather, _ := s.weatherService.GetWeatherCondition(ctx, airport.City)

		var current *weather_dto.CurrentWeatherDto
		if weather == nil {
			current = nil
		} else {
			current = weather.Current
		}

		res := dto.AirportWeatherDto{
			Object:  "airport_weather",
			Code:    airport.ICAOID,
			Airport: util.Ptr(dto.ToAirportDto(airport)),
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
