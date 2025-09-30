package service

import (
	"context"
	"database/sql"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	"flight-api/internal/model"
	"flight-api/internal/repository"
	"flight-api/pkg/logger"
	"flight-api/util"
	"fmt"

	"github.com/go-playground/validator"
)

type AirportService struct {
	airportRepository repository.IAirportRepository
	db                *sql.DB
	validate          *validator.Validate
	logger            *logger.Logger
}

func NewAirportService(airportRepository repository.IAirportRepository, db *sql.DB, validate *validator.Validate, logger *logger.Logger) IAirportService {
	return &AirportService{
		airportRepository: airportRepository,
		db:                db,
		validate:          validate,
		logger:            logger,
	}
}

func (s *AirportService) Create(ctx context.Context, r dto.AirportRequestDto) dto.AirportDto {
	err := s.validate.Struct(r)
	util.PanicIfError(err)

	tx, err := s.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	airport := model.Airport{
		SiteNumber:   util.ToSqlNullString(r.SiteNumber),
		ICAOID:       *r.ICAOID,
		FAAID:        util.ToSqlNullString(r.FAAID),
		IATAID:       util.ToSqlNullString(r.IATAID),
		Name:         util.ToSqlNullString(r.Name),
		Type:         util.ToSqlNullString((*string)(r.Type)),
		Status:       util.ToSqlNullBool(r.Status),
		Country:      util.ToSqlNullString(r.Country),
		State:        util.ToSqlNullString(r.State),
		StateFull:    util.ToSqlNullString(r.StateFull),
		County:       util.ToSqlNullString(r.County),
		City:         util.ToSqlNullString(r.City),
		Ownership:    util.ToSqlNullString((*string)(r.Ownership)),
		Use:          util.ToSqlNullString((*string)(r.Use)),
		Manager:      util.ToSqlNullString(r.Manager),
		ManagerPhone: util.ToSqlNullString(r.ManagerPhone),
		Latitude:     util.ToSqlNullString(r.Latitude),
		LatitudeSec:  util.ToSqlNullString(r.LatitudeSec),
		Longitude:    util.ToSqlNullString(r.Longitude),
		LongitudeSec: util.ToSqlNullString(r.LongitudeSec),
		Elevation:    util.ToSqlNullInt64(r.Elevation),
		ControlTower: util.ToSqlNullBool(r.ControlTower),
		Unicom:       util.ToSqlNullString(r.Unicom),
		CTAF:         util.ToSqlNullString(r.CTAF),
	}

	airport, err = s.airportRepository.Insert(ctx, tx, airport)
	util.PanicIfError(err)

	return dto.ToAirportDto(airport)
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

	fmt.Println("Update Service")

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

func (s *AirportService) fillUpdatableFields(airport *model.Airport, u dto.AirportUpdateDto) {
	util.ApplyUpdates(&airport.SiteNumber, util.ToSqlNullString(u.SiteNumber))
	util.ApplyUpdates(&airport.FAAID, util.ToSqlNullString(u.FAAID))
	util.ApplyUpdates(&airport.IATAID, util.ToSqlNullString(u.IATAID))
	util.ApplyUpdates(&airport.Name, util.ToSqlNullString(u.Name))
	util.ApplyUpdates(&airport.Type, util.ToSqlNullString((*string)(u.Type)))
	util.ApplyUpdates(&airport.Status, util.ToSqlNullBool(u.Status))
	util.ApplyUpdates(&airport.Country, util.ToSqlNullString(u.Country))
	util.ApplyUpdates(&airport.State, util.ToSqlNullString(u.State))
	util.ApplyUpdates(&airport.StateFull, util.ToSqlNullString(u.StateFull))
	util.ApplyUpdates(&airport.County, util.ToSqlNullString(u.County))
	util.ApplyUpdates(&airport.City, util.ToSqlNullString(u.City))
	util.ApplyUpdates(&airport.Ownership, util.ToSqlNullString((*string)(u.Ownership)))
	util.ApplyUpdates(&airport.Use, util.ToSqlNullString((*string)(u.Use)))
	util.ApplyUpdates(&airport.Manager, util.ToSqlNullString(u.Manager))
	util.ApplyUpdates(&airport.ManagerPhone, util.ToSqlNullString(u.ManagerPhone))
	util.ApplyUpdates(&airport.Latitude, util.ToSqlNullString(u.Latitude))
	util.ApplyUpdates(&airport.LatitudeSec, util.ToSqlNullString(u.LatitudeSec))
	util.ApplyUpdates(&airport.Longitude, util.ToSqlNullString(u.Longitude))
	util.ApplyUpdates(&airport.LongitudeSec, util.ToSqlNullString(u.LongitudeSec))
	util.ApplyUpdates(&airport.Elevation, util.ToSqlNullInt64(u.Elevation))
	util.ApplyUpdates(&airport.ControlTower, util.ToSqlNullBool(u.ControlTower))
	util.ApplyUpdates(&airport.Unicom, util.ToSqlNullString(u.Unicom))
	util.ApplyUpdates(&airport.CTAF, util.ToSqlNullString(u.CTAF))
	util.ApplyUpdates(&airport.EffectiveDate, util.ToSqlNullTime(u.EffectiveDate))
}
