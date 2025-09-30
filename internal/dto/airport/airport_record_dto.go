package airport_dto

import (
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"flight-api/util"
	"time"

	"github.com/google/uuid"
)

type AirportRecordDto struct {
	ID         uuid.UUID             `json:"id"`
	Object     string                `json:"object"`
	SiteNumber string                `json:"site_number"`
	ICAOID     string                `json:"icao_id"`
	FAAID      string                `json:"faa_id"`
	IATAID     string                `json:"iata_id"`
	Name       string                `json:"name"`
	Type       enum.FasilityTypeEnum `json:"type"`
	Status     bool                  `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

func ToAirportRecordDto(m model.Airport) AirportRecordDto {
	return AirportRecordDto{
		ID:         m.ID,
		Object:     "airport",
		SiteNumber: util.FromSqlNull[string](m.SiteNumber),
		ICAOID:     m.ICAOID,
		FAAID:      util.FromSqlNull[string](m.FAAID),
		IATAID:     util.FromSqlNull[string](m.IATAID),
		Name:       util.FromSqlNull[string](m.Name),
		Type:       enum.ToFacilityType(util.FromSqlNull[string](m.Type)),
		Status:     util.FromSqlNull[bool](m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToAirportRecordDtos(models []model.Airport) []AirportRecordDto {
	dtos := make([]AirportRecordDto, len(models))
	for i, m := range models {
		dtos[i] = ToAirportRecordDto(m)
	}

	return dtos
}
