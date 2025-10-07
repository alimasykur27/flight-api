package service_weather

import (
	"context"
	weather_dto "flight-api/internal/dto/weather"

	"github.com/stretchr/testify/mock"
)

type WeatherServiceMock struct {
	Mock mock.Mock
}

func (w *WeatherServiceMock) GetWeatherCondition(ctx context.Context, loc *string) (*weather_dto.WeatherDto, error) {
	args := w.Mock.Called(ctx, loc)
	var out *weather_dto.WeatherDto
	if v, ok := args.Get(0).(*weather_dto.WeatherDto); ok {
		out = v
	}
	return out, args.Error(1)
}
