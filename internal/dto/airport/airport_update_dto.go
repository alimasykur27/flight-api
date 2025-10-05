package airport_dto

import (
	"flight-api/internal/enum"
	"time"
)

type AirportUpdateDto struct {
	SiteNumber    *string            `json:"site_number" validate:"omitempty"`
	ICAOID        *string            `json:"icao_id" validate:"omitempty"`
	FAAID         *string            `json:"faa_id" validate:"omitempty"`
	IATAID        *string            `json:"iata_id" validate:"omitempty"`
	Name          *string            `json:"name" validate:"omitempty"`
	Type          *string            `json:"type" validate:"omitempty,oneof=small_airport medium_airport large_airport seaplane_hydrant heliport balloonport closed"`
	Status        *bool              `json:"status" validate:"omitempty"`
	Country       *string            `json:"country" validate:"omitempty"`
	State         *string            `json:"state" validate:"omitempty"`
	StateFull     *string            `json:"state_full" validate:"omitempty"`
	County        *string            `json:"county" validate:"omitempty"`
	City          *string            `json:"city" validate:"omitempty"`
	Ownership     enum.OwnershipEnum `json:"owership" validate:"omitempty,ownership"`
	Use           enum.UseTypeEnum   `json:"use" validate:"omitempty,use"`
	Manager       *string            `json:"manager" validate:"omitempty"`
	ManagerPhone  *string            `json:"manager_phone" validate:"omitempty"`
	Latitude      *string            `json:"latitude" validate:"omitempty"`
	LatitudeSec   *string            `json:"latitude_sec" validate:"omitempty"`
	Longitude     *string            `json:"longitude" validate:"omitempty"`
	LongitudeSec  *string            `json:"longitude_sec" validate:"omitempty"`
	Elevation     *int64             `json:"elevation" validate:"omitempty"`
	ControlTower  *bool              `json:"control_tower" validate:"omitempty"`
	Unicom        *string            `json:"unicom" validate:"omitempty"`
	CTAF          *string            `json:"ctaf" validate:"omitempty"`
	EffectiveDate *time.Time         `json:"effective_date" validate:"omitempty"`
}
