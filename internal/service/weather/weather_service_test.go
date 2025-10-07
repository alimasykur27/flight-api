package service_weather

import (
	"context"
	"encoding/json"
	"errors"
	"flight-api/config"
	location_dto "flight-api/internal/dto/location"
	weather_dto "flight-api/internal/dto/weather"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var log = logger.NewLogger(logger.DEBUG_LEVEL)

// Dummy Data Weather API
var DummyData = map[string]map[string]map[string]interface{}{
	"jakarta": {
		"location": {
			"name":            "Jakarta",
			"region":          "Jakarta Raya",
			"country":         "Indonesia",
			"lat":             -6.2146,
			"lon":             106.8451,
			"tz_id":           "Asia/Jakarta",
			"localtime_epoch": 1759705426,
			"localtime":       "2025-10-06 06:03",
		},
		"current": {
			"last_updated_epoch": 1759705200,
			"last_updated":       "2025-10-06 06:00",
			"temp_c":             25.2,
			"temp_f":             77.4,
			"is_day":             1,
			"condition": map[string]interface{}{
				"text": "Mist",
				"icon": "//cdn.weatherapi.com/weather/64x64/day/143.png",
				"code": 1030,
			},
			"wind_mph":    3.8,
			"wind_kph":    6.1,
			"wind_degree": 133,
			"wind_dir":    "SE",
			"pressure_mb": 1011.0,
			"pressure_in": 29.85,
			"precip_mm":   0.0,
			"precip_in":   0.0,
			"humidity":    83,
			"cloud":       25,
			"feelslike_c": 27.1,
			"feelslike_f": 80.7,
			"windchill_c": 27.0,
			"windchill_f": 80.5,
			"heatindex_c": 29.7,
			"heatindex_f": 85.4,
			"dewpoint_c":  22.2,
			"dewpoint_f":  72.0,
			"vis_km":      4.0,
			"vis_miles":   2.0,
			"uv":          0.0,
			"gust_mph":    5.2,
			"gust_kph":    8.3,
			"short_rad":   1.34,
			"diff_rad":    0.67,
			"dni":         0.0,
			"gti":         0.67,
		},
	},
	"tokyo": {
		"location": {
			"name":            "Tokyo",
			"region":          "Tokyo",
			"country":         "Japan",
			"lat":             35.6895,
			"lon":             139.6917,
			"tz_id":           "Asia/Tokyo",
			"localtime_epoch": 1759705568,
			"localtime":       "2025-10-06 08:06",
		},
		"current": {
			"last_updated_epoch": 1759705200,
			"last_updated":       "2025-10-06 08:00",
			"temp_c":             24.3,
			"temp_f":             75.7,
			"is_day":             1,
			"condition": map[string]interface{}{
				"text": "Partly cloudy",
				"icon": "//cdn.weatherapi.com/weather/64x64/day/116.png",
				"code": 1003,
			},
			"wind_mph":    8.3,
			"wind_kph":    13.3,
			"wind_degree": 353,
			"wind_dir":    "N",
			"pressure_mb": 1010.0,
			"pressure_in": 29.83,
			"precip_mm":   0.0,
			"precip_in":   0.0,
			"humidity":    89,
			"cloud":       75,
			"feelslike_c": 26.0,
			"feelslike_f": 78.8,
			"windchill_c": 25.6,
			"windchill_f": 78.0,
			"heatindex_c": 27.4,
			"heatindex_f": 81.3,
			"dewpoint_c":  20.3,
			"dewpoint_f":  68.5,
			"vis_km":      10.0,
			"vis_miles":   6.0,
			"uv":          1.3,
			"gust_mph":    10.6,
			"gust_kph":    17.1,
			"short_rad":   58.85,
			"diff_rad":    30.72,
			"dni":         0.0,
			"gti":         28.94,
		},
	},
}

// Mock Weather API
func newWeatherMockServer(log *logger.Logger, routes map[string]map[string]map[string]interface{}, mode string, delay time.Duration) *httptest.Server {
	log.Info("Mock Weather API server started")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Mock Weather API: Received request: ", r.URL.String())

		// Handle error simulation
		switch mode {
		case "timeout":
			time.Sleep(delay * time.Second)
			w.Header().Set("Content-Type", "application/json")
			return
		case "success-invalid-json":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid json`))
			return
		case "invalid-json":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{invalid json`))
			return
		case "invalid-error-code-type":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				map[string]map[string]interface{}{
					"error": {
						"code":    "not-a-number",
						"message": "An error occurred.",
					},
				},
			)
			return
		case "unknown-error-code":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				map[string]map[string]interface{}{
					"error": {
						"code":    9999,
						"message": "Unknown error occurred.",
					},
				},
			)
			return
		}

		// Handle success
		key := r.URL.Query().Get("key")
		q := r.URL.Query().Get("q")

		if key == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				map[string]map[string]interface{}{
					"error": {
						"code":    1002,
						"message": "API key is invalid or not provided.",
					},
				},
			)
			return
		}

		if q == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				map[string]map[string]interface{}{
					"error": {
						"code":    1003,
						"message": "Parameter q is missing.",
					},
				},
			)
			return
		}

		data, exists := routes[strings.ToLower(q)]
		if !exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				map[string]map[string]interface{}{
					"error": {
						"code":    1006,
						"message": "No matching location found.",
					},
				},
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data)
		log.Infof("Mock Weather API: Served data for location %s", q)
	})

	return httptest.NewServer(h)
}

// Mock Tripper
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type errReadCloser struct{ err error }

func (e errReadCloser) Read(p []byte) (int, error) { return 0, e.err }
func (e errReadCloser) Close() error               { return nil }

// -------------------- Unit Test -------------------

func TestNewWeatherService(t *testing.T) {
	cfg := &config.Config{
		WeatherURL:    "http://example.com",
		WeatherAPIKey: "test_api_key",
	}
	svc := NewWeatherService(log, cfg)
	assert.NotNil(t, svc, "WeatherService should not be nil")
	assert.Equal(t, log, svc.(*WeatherService).logger, "Logger should match")
	assert.Equal(t, cfg, svc.(*WeatherService).cfg, "Config should match")
	assert.NotNil(t, svc.(*WeatherService).client, "HTTP client should not be nil")
	assert.Equal(t, 60*time.Second, svc.(*WeatherService).client.Timeout, "HTTP client timeout should be 60 seconds")
}

func TestGetWeatherCondition_Success(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	tests := []struct {
		name     string
		input    string
		expected *weather_dto.WeatherDto
	}{
		{
			name:  "Valid Location - Jakarta",
			input: "jakarta",
			expected: &weather_dto.WeatherDto{
				Object: util.Ptr("weather"),
				Location: &location_dto.LocationDto{
					Name:           util.Ptr("Jakarta"),
					Country:        util.Ptr("Indonesia"),
					Lat:            util.Ptr(-6.2146),
					Lon:            util.Ptr(106.8451),
					TzId:           util.Ptr("Asia/Jakarta"),
					LocaltimeEpoch: util.Ptr(1759705426),
					Localtime:      util.Ptr("2025-10-06 06:03"),
				},
				Current: &weather_dto.CurrentWeatherDto{
					LastUpdatedEpoch: util.Ptr(1759705200),
					LastUpdated:      util.Ptr("2025-10-06 06:00"),
					TempC:            util.Ptr(25.2),
					TempF:            util.Ptr(77.4),
					IsDay:            util.Ptr(uint8(1)),
					Condition: &weather_dto.ConditionDto{
						Text: util.Ptr("Mist"),
						Icon: util.Ptr("//cdn.weatherapi.com/weather/64x64/day/143.png"),
						Code: util.Ptr(1030),
					},
					WindMph:    util.Ptr(3.8),
					WindKph:    util.Ptr(6.1),
					WindDegree: util.Ptr(133),
					WindDir:    util.Ptr("SE"),
					PressureMb: util.Ptr(1011.0),
					PressureIn: util.Ptr(29.85),
					PrecipMm:   util.Ptr(0.0),
					PrecipIn:   util.Ptr(0.0),
					Humidity:   util.Ptr(83),
					Cloud:      util.Ptr(25),
					FeelslikeC: util.Ptr(27.1),
					FeelslikeF: util.Ptr(80.7),
					WindchillC: util.Ptr(27.0),
					WindchillF: util.Ptr(80.5),
					HeatindexC: util.Ptr(29.7),
					HeatindexF: util.Ptr(85.4),
					DewpointC:  util.Ptr(22.2),
					DewpointF:  util.Ptr(72.0),
					VisKm:      util.Ptr(4.0),
					VisMiles:   util.Ptr(2.0),
					Uv:         util.Ptr(0.0),
					GustMph:    util.Ptr(5.2),
					GustKph:    util.Ptr(8.3),
					ShortRad:   util.Ptr(1.34),
					DiffRad:    util.Ptr(0.67),
					DNI:        util.Ptr(0.0),
					GTI:        util.Ptr(0.67),
				},
			},
		},
		{
			name:  "Valid Location - Tokyo",
			input: "tokyo",
			expected: &weather_dto.WeatherDto{
				Object: util.Ptr("weather"),
				Location: &location_dto.LocationDto{
					Name:           util.Ptr("Tokyo"),
					Country:        util.Ptr("Japan"),
					Lat:            util.Ptr(35.6895),
					Lon:            util.Ptr(139.6917),
					TzId:           util.Ptr("Asia/Tokyo"),
					LocaltimeEpoch: util.Ptr(1759705568),
					Localtime:      util.Ptr("2025-10-06 08:06"),
				},
				Current: &weather_dto.CurrentWeatherDto{
					LastUpdatedEpoch: util.Ptr(1759705200),
					LastUpdated:      util.Ptr("2025-10-06 08:00"),
					TempC:            util.Ptr(24.3),
					TempF:            util.Ptr(75.7),
					IsDay:            util.Ptr(uint8(1)),
					Condition: &weather_dto.ConditionDto{
						Text: util.Ptr("Partly cloudy"),
						Icon: util.Ptr("//cdn.weatherapi.com/weather/64x64/day/116.png"),
						Code: util.Ptr(1003),
					},
					WindMph:    util.Ptr(8.3),
					WindKph:    util.Ptr(13.3),
					WindDegree: util.Ptr(353),
					WindDir:    util.Ptr("N"),
					PressureMb: util.Ptr(1010.0),
					PressureIn: util.Ptr(29.83),
					PrecipMm:   util.Ptr(0.0),
					PrecipIn:   util.Ptr(0.0),
					Humidity:   util.Ptr(89),
					Cloud:      util.Ptr(75),
					FeelslikeC: util.Ptr(26.0),
					FeelslikeF: util.Ptr(78.8),
					WindchillC: util.Ptr(25.6),
					WindchillF: util.Ptr(78.0),
					HeatindexC: util.Ptr(27.4),
					HeatindexF: util.Ptr(81.3),
					DewpointC:  util.Ptr(20.3),
					DewpointF:  util.Ptr(68.5),
					VisKm:      util.Ptr(10.0),
					VisMiles:   util.Ptr(6.0),
					Uv:         util.Ptr(1.3),
					GustMph:    util.Ptr(10.6),
					GustKph:    util.Ptr(17.1),
					ShortRad:   util.Ptr(58.85),
					DiffRad:    util.Ptr(30.72),
					DNI:        util.Ptr(0.0),
					GTI:        util.Ptr(28.94),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetWeatherCondition(context.Background(), &tt.input)
			assert.NoError(t, err)

			// Check object type
			assert.Equal(t, "weather", *result.Object)

			// Check location
			assert.Equal(t, tt.expected.Location.Name, result.Location.Name, "Location field Name does not match")
			assert.Equal(t, tt.expected.Location.Country, result.Location.Country, "Location field Country does not match")
			assert.Equal(t, tt.expected.Location.Lat, result.Location.Lat, "Location field Lat does not match")
			assert.Equal(t, tt.expected.Location.Lon, result.Location.Lon, "Location field Lon does not match")
			assert.Equal(t, tt.expected.Location.TzId, result.Location.TzId, "Location field TzId does not match")
			assert.Equal(t, tt.expected.Location.LocaltimeEpoch, result.Location.LocaltimeEpoch, "Location field LocaltimeEpoch does not match")
			assert.Equal(t, tt.expected.Location.Localtime, result.Location.Localtime, "Location field Localtime does not match")

			// Check current weather
			assert.Equal(t, tt.expected.Current.LastUpdatedEpoch, result.Current.LastUpdatedEpoch, "Current field LastUpdatedEpoch does not match")
			assert.Equal(t, tt.expected.Current.LastUpdated, result.Current.LastUpdated, "Current field LastUpdated does not match")
			assert.Equal(t, tt.expected.Current.TempC, result.Current.TempC, "Current field TempC does not match")
			assert.Equal(t, tt.expected.Current.TempF, result.Current.TempF, "Current field TempF does not match")
			assert.Equal(t, tt.expected.Current.IsDay, result.Current.IsDay, "Current field IsDay does not match")

			// Check condition
			assert.Equal(t, tt.expected.Current.Condition.Text, result.Current.Condition.Text, "Condition field Text does not match")
			assert.Equal(t, tt.expected.Current.Condition.Icon, result.Current.Condition.Icon, "Condition field Icon does not match")
			assert.Equal(t, tt.expected.Current.Condition.Code, result.Current.Condition.Code, "Condition field Code does not match")
		})
	}
}

func TestGetWeatherCondition_Success_InvalidJSON(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "success-invalid-json", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	loc := "xxx"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrInternalServer)
}

func TestGetWeatherCondition_Error(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Invalid Location",
			input:    "invalid_location",
			expected: util.ErrNotFound,
		},
		{
			name:     "Empty Location",
			input:    "",
			expected: util.ErrBadRequest,
		},
		{
			name:     "Whitespace Location",
			input:    "   ",
			expected: util.ErrNotFound,
		},
		{
			name:     "Special Characters",
			input:    "@#$%^&*()!",
			expected: util.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetWeatherCondition(context.Background(), &tt.input)
			assert.ErrorIs(t, err, tt.expected)
		})
	}

}

func TestGetWeatherCondition_Error_NilLocation(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	_, err := svc.GetWeatherCondition(context.Background(), nil)
	assert.ErrorIs(t, err, util.ErrBadRequest)
}

func TestGetWeatherCondition_Error_Unauthorized(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL: srv.URL,
			// Missing API Key
		},
		logger: log,
		client: &http.Client{},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrUnauthorized)
}

func TestGetWeatherCondition_Error_Request(t *testing.T) {
	// Transport that returns error on any request
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("weird client error")
	})
	client := &http.Client{Transport: rt}

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    "http://invalid-url", // Invalid URL to trigger error
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: client,
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.NotNil(t, err, "expected error")
	assert.ErrorIs(t, err, util.ErrInternalServer, "expected ErrInternalServer")
}

func TestGetWeatherCondition_Error_Timeout(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "timeout", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL: srv.URL,
		},
		logger: log,
		client: &http.Client{
			Timeout: 5 * time.Millisecond,
		},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrGatewayTimeout)
}

func TestGetWeatherCondition_Error_InvalidJSON(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "invalid-json", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrInternalServer)
}

func TestGetWeatherCondition_Error_InvalidErrorCodeType(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "invalid-error-code-type", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.NotNil(t, err, "expected error")
	assert.ErrorIs(t, err, util.ErrInternalServer, "expected ErrInternalServer")
}

func TestGetWeatherCondition_Error_UnknownErrorCode(t *testing.T) {
	srv := newWeatherMockServer(log, DummyData, "unknown-error-code", 2)
	defer srv.Close()

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    srv.URL,
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrBadRequest)
}

func TestGetWeatherCondition_Error_Request_InvalidURL(t *testing.T) {
	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    "http://invalid-url", // Invalid URL to trigger error
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: &http.Client{},
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.ErrorIs(t, err, util.ErrInternalServer, "expected ErrInternalServer")
}

func TestGetWeatherCondition_Error_ReadBody(t *testing.T) {
	// Transport that returns a response with a body that errors on Read
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       errReadCloser{err: io.ErrUnexpectedEOF},
			Header:     make(http.Header),
		}, nil
	})

	client := &http.Client{Transport: rt}

	svc := &WeatherService{
		cfg: &config.Config{
			WeatherURL:    "http://valid-url", // Valid URL but body will error
			WeatherAPIKey: "test_api_key",
		},
		logger: log,
		client: client,
	}

	loc := "jakarta"
	_, err := svc.GetWeatherCondition(context.Background(), &loc)
	assert.NotNil(t, err, "expected error")
	assert.ErrorIs(t, err, util.ErrInternalServer, "expected ErrInternalServer")
}
