package service_weather

import (
	"context"
	"errors"
	"flight-api/config"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
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
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *WeatherService) GetWeatherCondition(ctx context.Context, loc *string) (*weather_dto.WeatherDto, error) {
	// Call WeatherAPIs
	s.logger.Debug("[GetWeatherCondition] Fetching weather data from Weather APIs...")

	if loc == nil {
		s.logger.Errorf("[GetWeatherCondition] 'loc' query parameter is required")
		return nil, util.NewErrorException(util.ErrBadRequest, "'loc' query parameter is required")
	}

	location := url.QueryEscape(strings.ToUpper(*loc))
	currentWeatherUrl := s.cfg.WeatherURL + "/v1/current.json"
	URL := currentWeatherUrl + "?key=" + s.cfg.WeatherAPIKey + "&q=" + location

	s.logger.Debugf("[GetWeatherCondition] Location: %s", *loc)
	resp, err := s.client.Get(URL)
	if err != nil {
		s.logger.Errorf("[GetWeatherCondition] Failed to fetch weather data: %v", err)

		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return nil, util.NewErrorException(util.ErrTimeout, "request to weather API timed out")
		}

		return nil, util.NewErrorException(util.ErrInternalServer, "failed to fetch weather data: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("[GetWeatherCondition] Failed to fetch weather data: %s", resp.Status)

		errBody, _ := io.ReadAll(resp.Body)
		var errData map[string]map[string]interface{}
		err = util.ParseJSON(errBody, &errData)
		if err != nil {
			s.logger.Errorf("[GetWeatherCondition] Failed to unmarshal error response body: %v", err)
			return nil, util.NewErrorException(util.ErrInternalServer, "failed to unmarshal error response body: "+err.Error())
		}

		errorCode, ok := errData["error"]["code"].(float64)
		if !ok {
			s.logger.Errorf("[GetWeatherCondition] Failed to assert error code to float64: %v", errData["error"]["code"])
			return nil, util.NewErrorException(util.ErrInternalServer, "failed to assert error code to float64: "+err.Error())
		}

		switch errorCode {
		case 1006:
			s.logger.Error("[GetWeatherCondition] Error code: ", errorCode)
			return nil, util.NewErrorException(util.ErrNotFound, "no matching location found")
		case 1003:
			s.logger.Error("[GetWeatherCondition] Error code: ", errorCode)
			return nil, util.NewErrorException(util.ErrBadRequest, "bad request")
		case 1002:
			s.logger.Error("[GetWeatherCondition] Error code: ", errorCode)
			return nil, util.NewErrorException(util.ErrUnauthorized, "unauthorized")
		default:
			return nil, util.NewErrorException(util.ErrBadRequest, "bad request")
		}
	}

	s.logger.Debug("[GetWeatherCondition] Successfully fetched weather data.")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("[GetWeatherCondition] Failed to read response body: %v", err)
		return nil, util.NewErrorException(util.ErrInternalServer, "failed to read response body: "+err.Error())
	}

	data := weather_dto.WeatherDto{
		Object:   util.Ptr("weather"),
		Location: nil,
		Current:  nil,
	}

	err = util.ParseJSON(body, &data)
	if err != nil {
		s.logger.Debugf("[GetWeatherCondition] Failed to unmarshal response body: %v", err)
		return nil, util.NewErrorException(util.ErrInternalServer, "failed to unmarshal response body: "+err.Error())
	}

	s.logger.Debug("[GetWeatherCondition] Finished fetching weather data.")
	return &data, nil
}
