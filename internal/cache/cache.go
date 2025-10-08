package cache

import (
	"context"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/model"
	"time"
)

type ICache interface {
	FindAirportByICAOID(ctx context.Context, icaoID string) (*model.Airport, error)
	CacheAirport(ctx context.Context, icaoID string, data *model.Airport, expiration time.Duration) error
	FindWeatherByLocation(ctx context.Context, location string) (*weather_dto.WeatherDto, error)
	CacheWeather(ctx context.Context, location string, data *weather_dto.WeatherDto, expiration time.Duration) error
}
