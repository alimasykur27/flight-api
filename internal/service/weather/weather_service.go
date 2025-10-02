package service_weather

import (
	"context"
	weather_dto "flight-api/internal/dto/weather"
)

type IWeatherService interface {
	GetWeatherCondition(ctx context.Context, loc *string) (*weather_dto.WeatherDto, error)
}
