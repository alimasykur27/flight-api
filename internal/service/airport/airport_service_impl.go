package service_airport

import (
	"context"
	"database/sql"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/model"
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

	s.fillUpdatableFields(&airport, u)
	updatedAirport, err := s.airportRepository.Update(ctx, tx, id, airport)
	util.PanicIfError(err)

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

	data := []dto.AirportWeatherDto{}
	airportWeather := dto.AirportWeatherDto{
		Object:  "airport_weather",
		Code:    airport.ICAOID,
		Airport: util.Ptr(dto.ToAirportDto(airport)),
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

func (s *AirportService) fillUpdatableFields(airport *model.Airport, u dto.AirportUpdateDto) {
	s.logger.Debug("[fillUpdatableFields] Filling updatable fields...")

	util.UpdateString(&airport.SiteNumber, u.SiteNumber)
	util.UpdateString(&airport.FAAID, u.FAAID)
	util.UpdateString(&airport.IATAID, u.IATAID)
	util.UpdateString(&airport.Name, u.Name)
	util.UpdateString(&airport.Type, (*string)(u.Type))
	util.UpdateBool(&airport.Status, u.Status)
	util.UpdateString(&airport.Country, u.Country)
	util.UpdateString(&airport.State, u.State)
	util.UpdateString(&airport.StateFull, u.StateFull)
	util.UpdateString(&airport.County, u.County)
	util.UpdateString(&airport.City, u.City)
	util.UpdateString(&airport.Ownership, (*string)(u.Ownership))
	util.UpdateString(&airport.Use, (*string)(u.Use))
	util.UpdateString(&airport.Manager, u.Manager)
	util.UpdateString(&airport.ManagerPhone, u.ManagerPhone)
	util.UpdateString(&airport.Latitude, u.Latitude)
	util.UpdateString(&airport.LatitudeSec, u.LatitudeSec)
	util.UpdateString(&airport.Longitude, u.Longitude)
	util.UpdateString(&airport.LongitudeSec, u.LongitudeSec)
	util.UpdateInt(&airport.Elevation, u.Elevation)
	util.UpdateBool(&airport.ControlTower, u.ControlTower)
	util.UpdateString(&airport.Unicom, u.Unicom)
	util.UpdateString(&airport.CTAF, u.CTAF)
	util.UpdateTime(&airport.EffectiveDate, u.EffectiveDate)
}
