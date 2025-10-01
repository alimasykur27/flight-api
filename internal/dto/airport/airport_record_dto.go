package airport_dto

import (
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type AirportRecordDto struct {
	ID         *uuid.UUID            `json:"id"`
	Object     *string               `json:"object"`
	SiteNumber *string               `json:"site_number"`
	ICAOID     *string               `json:"icao_id"`
	FAAID      *string               `json:"faa_id"`
	IATAID     *string               `json:"iata_id"`
	Name       *string               `json:"name"`
	Type       enum.FasilityTypeEnum `json:"type"`
	Status     *bool                 `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

func ToAirportRecordDto(m model.Airport) AirportRecordDto {
	object := "airport"
	return AirportRecordDto{
		ID:         m.ID,
		Object:     &object,
		SiteNumber: m.SiteNumber,
		ICAOID:     m.ICAOID,
		FAAID:      m.FAAID,
		IATAID:     m.IATAID,
		Name:       m.Name,
		Type:       m.Type,
		Status:     m.Status,
		CreatedAt:  *m.CreatedAt,
		UpdatedAt:  *m.UpdatedAt,
	}
}

func ToAirportRecordDtos(models []model.Airport) []AirportRecordDto {
	dtos := make([]AirportRecordDto, len(models))
	for i, m := range models {
		dtos[i] = ToAirportRecordDto(m)
	}

	return dtos
}
