package cache

import (
	"context"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type CacheMock struct {
	Mock mock.Mock
}

func (c *CacheMock) FindAirportByICAOID(ctx context.Context, icaoID string) (*model.Airport, error) {
	return &model.Airport{}, nil
}

func (c *CacheMock) CacheAirport(ctx context.Context, icaoID string, data *model.Airport, expiration time.Duration) error {
	return nil
}

func (c *CacheMock) FindWeatherByLocation(ctx context.Context, location string) (*weather_dto.WeatherDto, error) {
	return nil, nil
}

func (c *CacheMock) CacheWeather(ctx context.Context, location string, data *weather_dto.WeatherDto, expiration time.Duration) error {
	return nil
}
