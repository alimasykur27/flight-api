package service_aviation

import (
	"context"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	aviation_dto "flight-api/internal/dto/aviation"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net/http"
	"time"
)

type AviationService struct {
	logger *logger.Logger
	cfg    *config.Config
}

func NewAviationService(logger *logger.Logger, cfg *config.Config) IAviationService {
	return &AviationService{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *AviationService) FetchAirportData(ctx context.Context, icaoCodes []string) (map[string]airport_dto.AirportRequestDto, error) {
	var airportsData = make(map[string]airport_dto.AirportRequestDto)

	// Call Aviation API.
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	s.logger.Debug("[FecthAirportData] Fetching airport data from Aviation API...")
	for _, code := range icaoCodes {
		s.logger.Debugf("Fetching data for ICAO code: %s", code)
		URL := s.cfg.AviationURL + "/airports?apt=" + code

		s.logger.Debugf("Request URL: %s", URL)
		resp, err := client.Get(URL)

		if err != nil {
			s.logger.Errorf("[FetchiAirportData] fetching data for ICAO code %s: %v", code, err)
			switch err {
			case http.ErrHandlerTimeout:
				return nil, util.ErrInternalServer
			case http.ErrMissingFile:
				return nil, util.ErrNotFound
			default:
				return nil, util.ErrBadRequest
			}
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s.logger.Errorf("[FetchiAirportData] Failed to fetch data for ICAO code %s: %s", code, resp.Status)
			return nil, util.ErrBadRequest
		}
		s.logger.Debugf("[FetchiAirportData] Successfully fetched data for ICAO code %s", code)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Errorf("[FetchiAirportData] Failed to read response body for ICAO code %s: %v", code, err)
			return nil, util.ErrInternalServer
		}

		var data map[string][]aviation_dto.AviationAirportDto

		err = util.ParseJSON(body, &data)
		if err != nil {
			s.logger.Debugf("Failed to unmarshal response body for ICAO code %s: %v", code, err)
			return nil, util.ErrInternalServer
		}

		for _, code := range icaoCodes {
			s.logger.Debugf("Processing data for ICAO code: %s", code)
			airports, exists := data[code]

			if !exists || len(airports) == 0 {
				s.logger.Warnf("No data found for ICAO code: %s", code)
				airportsData[code] = airport_dto.AirportRequestDto{}
				continue
			}

			airport := airports[0]
			airportData := aviation_dto.ToAirportRequestDto(airport)
			airportsData[code] = airportData
		}
	}

	s.logger.Debug("Finished fetching airport data.")
	return airportsData, nil
}
