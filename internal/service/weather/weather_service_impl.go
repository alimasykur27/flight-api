package service_weather

import (
	"context"
	"flight-api/config"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type WeatherService struct {
	logger *logger.Logger
	cfg    *config.Config
	client *http.Client
}

func NewWeatherService(logger *logger.Logger, cfg *config.Config) IWeatherService {
	return &WeatherService{
		logger: logger,
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (s *WeatherService) GetWeatherCondition(ctx context.Context, loc *string) (*weather_dto.WeatherDto, error) {
	// Call WeatherAPIs
	s.logger.Debug("[GetWeatherCondition] Fetching weather data from Weather APIs...")

	if loc == nil {
		s.logger.Errorf("[GetWeatherCondition] Failed to fetch weather data: %v", util.ErrBadRequest)
		return nil, util.ErrBadRequest
	}

	currentWeatherUrl := s.cfg.WeatherURL + "/v1/current.json"

	location := url.QueryEscape(strings.ToUpper(*loc))
	URL := currentWeatherUrl + "?key=" + s.cfg.WeatherAPIKey + "&q=" + location

	s.logger.Debugf("Request URL: %s", URL)
	resp, err := s.client.Get(URL)

	if err == http.ErrHandlerTimeout {
		s.logger.Errorf("[GetWeatherCondition] Failed to fetch weather data: %v", err)
		return nil, util.ErrInternalServer
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("[GetWeatherCondition] Failed to fetch weather data: %s", resp.Status)

		errBody, _ := io.ReadAll(resp.Body)

		var errData map[string]map[string]interface{}
		err = util.ParseJSON(errBody, &errData)
		if err != nil {
			return nil, util.ErrInternalServer
		}

		errorCode, ok := errData["error"]["code"].(float64)
		if !ok {
			s.logger.Errorf("[GetWeatherCondition] Failed to assert error code to float64: %v", errData["error"]["code"])
			return nil, util.ErrInternalServer
		}

		if errorCode == 1006 {
			s.logger.Debug("Error code: ", errorCode)
			return nil, util.ErrNotFound
		}

		return nil, util.ErrBadRequest
	}
	s.logger.Debug("[GetWeatherCondition] Successfully fetched weather data.")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("[GetWeatherCondition] Failed to read response body: %v", err)
		return nil, util.ErrInternalServer
	}

	data := weather_dto.WeatherDto{
		Object:   util.Ptr("weather"),
		Location: nil,
		Current:  nil,
	}

	err = util.ParseJSON(body, &data)
	if err != nil {
		s.logger.Debugf("Failed to unmarshal response body: %v", err)
		return nil, util.ErrInternalServer
	}

	s.logger.Debug("Finished fetching weather data.")
	return &data, nil
}
