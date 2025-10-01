package sync_dto

type SyncAirportRequest struct {
	ICAOCodes []string `json:"icao_codes" validate:"required,dive,required"`
}
