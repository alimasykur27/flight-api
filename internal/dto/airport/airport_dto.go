package airport_dto

import (
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type AirportDto struct {
	ID            *uuid.UUID            `json:"id"`
	Object        *string               `json:"object"`
	SiteNumber    *string               `json:"site_number"`
	ICAOID        *string               `json:"icao_id"`
	FAAID         *string               `json:"faa_id"`
	IATAID        *string               `json:"iata_id"`
	Name          *string               `json:"name"`
	Type          enum.FasilityTypeEnum `json:"type"`
	Status        *bool                 `json:"status"`
	Country       *string               `json:"country"`
	State         *string               `json:"state"`
	StateFull     *string               `json:"state_full"`
	County        *string               `json:"county"`
	City          *string               `json:"city" `
	Ownership     enum.OwnershipEnum    `json:"owership"`
	Use           enum.UseTypeEnum      `json:"use"`
	Manager       *string               `json:"manager"`
	ManagerPhone  *string               `json:"manager_phone"`
	Latitude      *string               `json:"latitude"`
	LatitudeSec   *string               `json:"latitude_sec"`
	Longitude     *string               `json:"longitude"`
	LongitudeSec  *string               `json:"longitude_sec"`
	Elevation     *int64                `json:"elevation"`
	ControlTower  *bool                 `json:"control_tower"`
	Unicom        *string               `json:"unicom"`
	CTAF          *string               `json:"ctaf"`
	EffectiveDate *time.Time            `json:"effective_date"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

func NewAirportDto(
	id *uuid.UUID,
	siteNumber *string,
	icaoID *string,
	faaID *string,
	iataID *string,
	name *string,
	typee enum.FasilityTypeEnum,
	status *bool,
	country *string,
	state *string,
	stateFull *string,
	county *string,
	city *string,
	ownership enum.OwnershipEnum,
	use enum.UseTypeEnum,
	manager *string,
	managerPhone *string,
	latitude *string,
	latitudeSec *string,
	longitude *string,
	longitudeSec *string,
	elevation *int64,
	controlTower *bool,
	unicom *string,
	ctaf *string,
	effectiveDate *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) AirportDto {
	object := "airport"
	return AirportDto{
		ID:            id,
		Object:        &object,
		SiteNumber:    siteNumber,
		ICAOID:        icaoID,
		FAAID:         faaID,
		IATAID:        iataID,
		Name:          name,
		Type:          typee,
		Status:        status,
		Country:       country,
		State:         state,
		StateFull:     stateFull,
		County:        county,
		City:          city,
		Ownership:     ownership,
		Use:           use,
		Manager:       manager,
		ManagerPhone:  managerPhone,
		Latitude:      latitude,
		LatitudeSec:   latitudeSec,
		Longitude:     longitude,
		LongitudeSec:  longitudeSec,
		Elevation:     elevation,
		ControlTower:  controlTower,
		Unicom:        unicom,
		CTAF:          ctaf,
		EffectiveDate: effectiveDate,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func ToAirportDto(m model.Airport) AirportDto {
	object := "airport"
	return AirportDto{
		ID:            m.ID,
		Object:        &object,
		SiteNumber:    m.SiteNumber,
		ICAOID:        m.ICAOID,
		FAAID:         m.FAAID,
		IATAID:        m.IATAID,
		Name:          m.Name,
		Type:          m.Type,
		Status:        m.Status,
		Country:       m.Country,
		State:         m.State,
		StateFull:     m.StateFull,
		County:        m.County,
		City:          m.City,
		Ownership:     m.Ownership,
		Use:           m.Use,
		Manager:       m.Manager,
		ManagerPhone:  m.ManagerPhone,
		Latitude:      m.Latitude,
		LatitudeSec:   m.LatitudeSec,
		Longitude:     m.Longitude,
		LongitudeSec:  m.LongitudeSec,
		Elevation:     m.Elevation,
		ControlTower:  m.ControlTower,
		Unicom:        m.Unicom,
		CTAF:          m.CTAF,
		EffectiveDate: m.EffectiveDate,
		CreatedAt:     *m.CreatedAt,
		UpdatedAt:     *m.UpdatedAt,
	}
}

func ToAirportDtos(m []model.Airport) []AirportDto {
	var airportSlice []AirportDto
	for _, airport := range m {
		airportSlice = append(airportSlice, ToAirportDto(airport))
	}

	return airportSlice
}
