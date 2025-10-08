package cache

import (
	"context"
	"flight-api/config"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/internal/model"
	"flight-api/pkg/logger"
	"flight-api/util"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	logger      *logger.Logger
	cfg         *config.Config
	redisClient *redis.Client
	defaultKey  string
}

func NewCache(logger *logger.Logger, cfg *config.Config, redisClient *redis.Client) ICache {
	return &Cache{
		logger:      logger,
		cfg:         cfg,
		redisClient: redisClient,
		defaultKey:  cfg.ServiceName,
	}
}

func (c *Cache) FindAirportByICAOID(ctx context.Context, icaoID string) (*model.Airport, error) {
	var data model.Airport

	key := c.defaultKey + ":airport:" + icaoID
	cacheData, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	err = util.ParseJSON([]byte(cacheData), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Cache) CacheAirport(ctx context.Context, icaoID string, data *model.Airport, expiration time.Duration) error {
	key := c.defaultKey + ":airport:" + strings.ReplaceAll(strings.ToLower(icaoID), " ", "_")

	jsonData, err := util.ToJSON(data)
	if err != nil {
		c.logger.Errorf("[CacheWeather] Failed to marshal weather data to JSON: %v", err)
		return err
	}

	return c.redisClient.Set(ctx, key, jsonData, expiration).Err()
}

func (c *Cache) FindWeatherByLocation(ctx context.Context, location string) (*weather_dto.WeatherDto, error) {
	var data weather_dto.WeatherDto

	key := c.defaultKey + ":weather:" + location
	cacheData, err := c.redisClient.Get(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	c.logger.Debug("[GetWeatherCondition] Found cached weather data in Redis.")
	err = util.ParseJSON([]byte(cacheData), &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Cache) CacheWeather(ctx context.Context, location string, data *weather_dto.WeatherDto, expiration time.Duration) error {
	key := c.defaultKey + ":weather:" + location

	jsonData, err := util.ToJSON(data)
	if err != nil {
		c.logger.Errorf("[CacheWeather] Failed to marshal weather data to JSON: %v", err)
		return err
	}

	return c.redisClient.Set(ctx, key, jsonData, expiration).Err()
}
