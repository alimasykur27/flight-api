package sync_dto

import airport_dto "flight-api/internal/dto/airport"

type SyncAirportResponse struct {
	ICAOCode string                  `json:"icao_code"`
	Airport  *airport_dto.AirportDto `json:"airport"`
	Status   string                  `json:"status"`
	Message  string                  `json:"message"`
}
