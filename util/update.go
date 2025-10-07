package util

import (
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/model"
	"time"
)

// UpdateString applies updates for a pointer to string.
func UpdateString(dst **string, src *string) {
	if src != nil {
		*dst = src
	}
}

// UpdateBool applies updates for a pointer to bool.
func UpdateBool(dst **bool, src *bool) {
	if src != nil {
		*dst = src
	}
}

// UpdateInt applies updates for a pointer to int.
func UpdateInt[T int64 | int32 | int16 | int](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}

// UpdateTime applies updates for a pointer to time.Time.
func UpdateTime(dst **time.Time, src *time.Time) {
	if src != nil {
		*dst = src
	}
}

// FillUpdatedFields
func FillUpdatableFields(airport *model.Airport, u airport_dto.AirportUpdateDto) {
	UpdateString(&airport.SiteNumber, u.SiteNumber)
	UpdateString(&airport.FAAID, u.FAAID)
	UpdateString(&airport.IATAID, u.IATAID)
	UpdateString(&airport.Name, u.Name)
	UpdateString(&airport.Type, (*string)(u.Type))
	UpdateBool(&airport.Status, u.Status)
	UpdateString(&airport.Country, u.Country)
	UpdateString(&airport.State, u.State)
	UpdateString(&airport.StateFull, u.StateFull)
	UpdateString(&airport.County, u.County)
	UpdateString(&airport.City, u.City)
	UpdateString(&airport.Ownership, (*string)(u.Ownership))
	UpdateString(&airport.Use, (*string)(u.Use))
	UpdateString(&airport.Manager, u.Manager)
	UpdateString(&airport.ManagerPhone, u.ManagerPhone)
	UpdateString(&airport.Latitude, u.Latitude)
	UpdateString(&airport.LatitudeSec, u.LatitudeSec)
	UpdateString(&airport.Longitude, u.Longitude)
	UpdateString(&airport.LongitudeSec, u.LongitudeSec)
	UpdateInt(&airport.Elevation, u.Elevation)
	UpdateBool(&airport.ControlTower, u.ControlTower)
	UpdateString(&airport.Unicom, u.Unicom)
	UpdateString(&airport.CTAF, u.CTAF)
	UpdateTime(&airport.EffectiveDate, u.EffectiveDate)
}
