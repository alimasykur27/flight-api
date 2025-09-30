package airport_dto

import (
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"flight-api/util"
	"time"

	"github.com/google/uuid"
)

type AirportDto struct {
	ID            uuid.UUID             `json:"id"`
	Object        string                `json:"object"`
	SiteNumber    string                `json:"site_number"`
	ICAOID        string                `json:"icao_id"`
	FAAID         string                `json:"faa_id"`
	IATAID        string                `json:"iata_id"`
	Name          string                `json:"name"`
	Type          enum.FasilityTypeEnum `json:"type"`
	Status        bool                  `json:"status"`
	Country       string                `json:"country"`
	State         string                `json:"state"`
	StateFull     string                `json:"state_full"`
	County        string                `json:"county"`
	City          string                `json:"city" `
	Ownership     enum.OwnershipEnum    `json:"owership"`
	Use           enum.UseTypeEnum      `json:"use"`
	Manager       string                `json:"manager"`
	ManagerPhone  string                `json:"manager_phone"`
	Latitude      string                `json:"latitude"`
	LatitudeSec   string                `json:"latitude_sec"`
	Longitude     string                `json:"longitude"`
	LongitudeSec  string                `json:"longitude_sec"`
	Elevation     int64                 `json:"elevation"`
	ControlTower  bool                  `json:"control_tower"`
	Unicom        string                `json:"unicom"`
	CTAF          string                `json:"ctaf"`
	EffectiveDate time.Time             `json:"effective_date"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

func ToAirportDto(m model.Airport) AirportDto {
	return AirportDto{
		ID:            m.ID,
		Object:        "airport",
		SiteNumber:    util.FromSqlNull[string](m.SiteNumber),
		ICAOID:        m.ICAOID,
		FAAID:         util.FromSqlNull[string](m.FAAID),
		IATAID:        util.FromSqlNull[string](m.IATAID),
		Name:          util.FromSqlNull[string](m.Name),
		Type:          enum.ToFacilityType(util.FromSqlNull[string](m.Type)),
		Status:        util.FromSqlNull[bool](m.Status),
		Country:       util.FromSqlNull[string](m.Country),
		State:         util.FromSqlNull[string](m.State),
		StateFull:     util.FromSqlNull[string](m.StateFull),
		County:        util.FromSqlNull[string](m.County),
		City:          util.FromSqlNull[string](m.City),
		Ownership:     enum.ToOwnership(util.FromSqlNull[string](m.Ownership)),
		Use:           enum.ToUseType(util.FromSqlNull[string](m.Use)),
		Manager:       util.FromSqlNull[string](m.Manager),
		ManagerPhone:  util.FromSqlNull[string](m.ManagerPhone),
		Latitude:      util.FromSqlNull[string](m.Latitude),
		LatitudeSec:   util.FromSqlNull[string](m.LatitudeSec),
		Longitude:     util.FromSqlNull[string](m.Longitude),
		LongitudeSec:  util.FromSqlNull[string](m.LongitudeSec),
		Elevation:     util.FromSqlNull[int64](m.Elevation),
		ControlTower:  util.FromSqlNull[bool](m.ControlTower),
		Unicom:        util.FromSqlNull[string](m.Unicom),
		CTAF:          util.FromSqlNull[string](m.CTAF),
		EffectiveDate: util.FromSqlNull[time.Time](m.EffectiveDate),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToAirportDtos(m []model.Airport) []AirportDto {
	var airportSlice []AirportDto
	for _, airport := range m {
		airportSlice = append(airportSlice, ToAirportDto(airport))
	}

	return airportSlice
}
