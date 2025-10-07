package service_aviation

import (
	"context"
	"errors"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	aviation_dto "flight-api/internal/dto/aviation"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type AviationService struct {
	logger *logger.Logger
	cfg    *config.Config
	client *http.Client
}

func NewAviationService(logger *logger.Logger, cfg *config.Config) IAviationService {
	return &AviationService{
		logger: logger,
		cfg:    cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *AviationService) FetchAirportData(ctx context.Context, icaoCodes []string) (map[string]airport_dto.AirportRequestDto, error) {
	var airportsData = make(map[string]airport_dto.AirportRequestDto)

	// Call Aviation API
	s.logger.Debug("[FecthAirportData] Fetching airport data from Aviation API...")

	// Combine slice into separated comma string
	apt := strings.Join(icaoCodes, ",")
	URL := s.cfg.AviationURL + "/airports?apt=" + apt
	s.logger.Debugf("Request URL: %s", URL)

	s.logger.Debugf("[FetchiAirportData] ICAO code: %s", apt)
	resp, err := s.client.Get(URL)
	if err != nil {
		s.logger.Errorf("[FetchiAirportData] Error fetching data for ICAO code %s: %v", apt, err)

		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return nil, util.ErrGatewayTimeout
		}

		return nil, util.ErrInternalServer

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("[FetchiAirportData] Failed to fetch data for ICAO code %s: %s", apt, resp.Status)
		return nil, util.ErrBadRequest
	}
	s.logger.Debugf("[FetchiAirportData] Successfully fetched data for ICAO code %s", apt)

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("[FetchiAirportData] Failed to read response body for ICAO code %s: %v", apt, err)
		return nil, util.ErrInternalServer
	}

	var bodyData map[string][]aviation_dto.AviationAirportDto

	err = util.ParseJSON(rawBody, &bodyData)
	if err != nil {
		s.logger.Errorf("[FetchiAirportData] Failed to unmarshal response body for ICAO code %s: %v", apt, err)
		return nil, util.ErrInternalServer
	}

	for _, code := range icaoCodes {
		s.logger.Debugf("[FetchiAirportData] Processing data for ICAO code: %s", code)
		airports, exists := bodyData[code]

		if !exists || len(airports) == 0 {
			s.logger.Warnf("[FetchiAirportData] No data found for ICAO code: %s", code)
			airportsData[code] = airport_dto.AirportRequestDto{}
			continue
		}

		airport := airports[0]
		airportData := aviation_dto.ToAirportRequestDto(airport)
		airportsData[code] = airportData
	}

	s.logger.Debug("Finished fetching airport data.")
	return airportsData, nil
}
