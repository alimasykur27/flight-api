package service

import (
	"context"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
)

type IAirportService interface {
	Create(ctx context.Context, r dto.AirportRequestDto) dto.AirportDto
	FindAll(ctx context.Context, p queryparams.QueryParams) pagination_dto.PaginationDto
	FindByID(ctx context.Context, id string) (dto.AirportDto, error)
	Update(ctx context.Context, id string, u dto.AirportUpdateDto) (dto.AirportDto, error)
	Delete(ctx context.Context, id string) error
}
