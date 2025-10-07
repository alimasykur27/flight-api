package service_aviation

import (
	"context"
	airport_dto "flight-api/internal/dto/airport"

	"github.com/stretchr/testify/mock"
)

type AviationServiceMock struct {
	Mock mock.Mock
}

// FetchAirportData(ctx, icaoCodes) (map[string]AirportRequestDto, error)
func (m *AviationServiceMock) FetchAirportData(ctx context.Context, icaoCodes []string) (map[string]airport_dto.AirportRequestDto, error) {
	args := m.Mock.Called(ctx, icaoCodes)

	var out map[string]airport_dto.AirportRequestDto
	if v, ok := args.Get(0).(map[string]airport_dto.AirportRequestDto); ok {
		out = v
	}
	return out, args.Error(1)
}
