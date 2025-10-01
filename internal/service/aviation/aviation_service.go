package service_aviation

import (
	"context"
	airport_dto "flight-api/internal/dto/airport"
)

type IAviationService interface {
	FetchAirportData(ctx context.Context, icaoCodes []string) (map[string]airport_dto.AirportRequestDto, error)
}
