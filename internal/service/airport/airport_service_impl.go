package service_airport

import (
	"context"
	"database/sql"
	airport_dto "flight-api/internal/dto/airport"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	"flight-api/internal/model"
	repository_airport "flight-api/internal/repository/airport"

	"flight-api/pkg/logger"
	"flight-api/util"

	"github.com/go-playground/validator"
)

type AirportService struct {
	logger            *logger.Logger
	validate          *validator.Validate
	db                *sql.DB
	airportRepository repository_airport.IAirportRepository
}

func NewAirportService(
	logger *logger.Logger,
	validate *validator.Validate,
	db *sql.DB,
	airportRepository repository_airport.IAirportRepository,
) IAirportService {
	return &AirportService{
		logger:            logger,
		validate:          validate,
		db:                db,
		airportRepository: airportRepository,
	}
}

func (s *AirportService) Create(ctx context.Context, r dto.AirportRequestDto) dto.AirportDto {
	err := s.validate.Struct(r)
	util.PanicIfError(err)

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport := airport_dto.AirportRequestToAirport(r)
	airport, err = s.airportRepository.Insert(ctx, tx, airport)
	util.PanicIfError(err)

	data := dto.ToAirportDto(airport)
	return data
}

func (s *AirportService) FindAll(ctx context.Context, query queryparams.QueryParams) pagination_dto.PaginationDto {
	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airports, total, err := s.airportRepository.FindAll(ctx, tx, query.Limit, query.Offset)
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
		Meta: pagination_dto.PaginationMetaDto{
			Limit: query.Limit,
			Page:  query.Page,
			Next:  hasNext,
		},
	}
}

func (s *AirportService) FindByID(ctx context.Context, id string) (dto.AirportDto, error) {
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

	s.fillUpdatableFields(airport, u)
	updatedAirport, err := s.airportRepository.Update(ctx, tx, id, airport)
	util.PanicIfError(err)

	return dto.ToAirportDto(updatedAirport), nil
}

func (s *AirportService) Delete(ctx context.Context, id string) error {
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

func (s *AirportService) fillUpdatableFields(airport model.Airport, u dto.AirportUpdateDto) {
	util.UpdateString(airport.SiteNumber, u.SiteNumber)
	util.UpdateString(airport.FAAID, u.FAAID)
	util.UpdateString(airport.IATAID, u.IATAID)
	util.UpdateString(airport.Name, u.Name)
	util.UpdateString(airport.Type, (*string)(u.Type))
	util.UpdateBool(airport.Status, u.Status)
	util.UpdateString(airport.Country, u.Country)
	util.UpdateString(airport.State, u.State)
	util.UpdateString(airport.StateFull, u.StateFull)
	util.UpdateString(airport.County, u.County)
	util.UpdateString(airport.City, u.City)
	util.UpdateString(airport.Ownership, (*string)(u.Ownership))
	util.UpdateString(airport.Use, (*string)(u.Use))
	util.UpdateString(airport.Manager, u.Manager)
	util.UpdateString(airport.ManagerPhone, u.ManagerPhone)
	util.UpdateString(airport.Latitude, u.Latitude)
	util.UpdateString(airport.LatitudeSec, u.LatitudeSec)
	util.UpdateString(airport.Longitude, u.Longitude)
	util.UpdateString(airport.LongitudeSec, u.LongitudeSec)
	util.UpdateInt(airport.Elevation, u.Elevation)
	util.UpdateBool(airport.ControlTower, u.ControlTower)
	util.UpdateString(airport.Unicom, u.Unicom)
	util.UpdateString(airport.CTAF, u.CTAF)
	util.UpdateTime(airport.EffectiveDate, u.EffectiveDate)
}
