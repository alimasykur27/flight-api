package service_airport

import (
	"context"
	airport_dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
)

type IAirportService interface {
	Seeding(ctx context.Context, reqs []string) ([]airport_dto.AirportDto, error)
	Create(ctx context.Context, r airport_dto.AirportRequestDto) (airport_dto.AirportDto, error)
	FindAll(ctx context.Context, p queryparams.QueryParams) (pagination_dto.PaginationDto, error)
	FindByID(ctx context.Context, id string) (airport_dto.AirportDto, error)
	Update(ctx context.Context, id string, u airport_dto.AirportUpdateDto) (airport_dto.AirportDto, error)
	Delete(ctx context.Context, id string) error
	GetWeatherCondition(ctx context.Context, code string, name string, query queryparams.QueryParams) (*pagination_dto.PaginationDto, error)
}
