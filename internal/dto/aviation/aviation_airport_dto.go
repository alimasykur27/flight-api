package aviation_dto

import (
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/enum"
	"flight-api/util"
	"strings"
	"time"
)

type AviationAirportDto struct {
	SiteNumber              string `json:"site_number"`
	Type                    string `json:"type"`
	FacilityName            string `json:"facility_name"`
	FAAIdentifier           string `json:"faa_ident"`
	ICAOIdentifier          string `json:"icao_ident"`
	Region                  string `json:"region"`
	DistrictOffice          string `json:"district_office"`
	State                   string `json:"state"`
	StateFull               string `json:"state_full"`
	County                  string `json:"county"`
	City                    string `json:"city"`
	Ownership               string `json:"ownership"`
	Use                     string `json:"use"`
	Manager                 string `json:"manager"`
	ManagerPhone            string `json:"manager_phone"`
	Latitude                string `json:"latitude"`
	LatitudeSec             string `json:"latitude_sec"`
	Longitude               string `json:"longitude"`
	LongitudeSec            string `json:"longitude_sec"`
	Elevation               string `json:"elevation"`
	MagneticVariation       string `json:"magnetic_variation"`
	TPA                     string `json:"tpa"`
	VFRSectional            string `json:"vfr_sectional"`
	NotamFacilityIdentifier string `json:"notam_facility_ident"`
	Status                  string `json:"status"`
	ControlTower            string `json:"control_tower"`
	UNICOM                  string `json:"unicom"`
	CTAF                    string `json:"ctaf"`
	EffectiveDate           string `json:"effective_date"`
}

func ToAirportRequestDto(source AviationAirportDto) airport_dto.AirportRequestDto {
	return airport_dto.AirportRequestDto{
		SiteNumber:    &source.SiteNumber,
		ICAOID:        &source.ICAOIdentifier,
		FAAID:         &source.FAAIdentifier,
		IATAID:        &source.FAAIdentifier,
		Name:          &source.FacilityName,
		Type:          enum.ToFacilityType(strings.ToLower(source.Type)),
		Status:        util.Ptr((ToAirportStatus(source.Status))),
		Country:       nil,
		State:         &source.State,
		StateFull:     &source.StateFull,
		County:        &source.Region,
		City:          &source.City,
		Ownership:     ToAirportOwnership(source.Ownership),
		Use:           ToAirportUse(source.Use),
		Manager:       &source.Manager,
		ManagerPhone:  &source.ManagerPhone,
		Latitude:      &source.Latitude,
		LatitudeSec:   &source.LatitudeSec,
		Longitude:     &source.Longitude,
		LongitudeSec:  &source.LongitudeSec,
		Elevation:     util.ParseInt64Ptr(source.Elevation),
		ControlTower:  util.Ptr(ToControlTower(source.ControlTower)),
		Unicom:        &source.UNICOM,
		CTAF:          &source.CTAF,
		EffectiveDate: ToAirportEffectiveDate(source.EffectiveDate),
	}
}

func ToAirportStatus(status string) bool {
	return status == "O"
}

func ToControlTower(controlTower string) bool {
	return controlTower == "Y"
}

func ToAirportOwnership(ownership string) enum.OwnershipEnum {
	switch ownership {
	case "PU":
		return enum.OWN_PUBLIC
	case "PR":
		return enum.OWN_PRIVATE
	default:
		return enum.OWN_NIL
	}
}

func ToAirportUse(use string) enum.UseTypeEnum {
	switch use {
	case "PU":
		return enum.USE_PUBLIC
	case "PR":
		return enum.USE_PRIVATE
	default:
		return enum.USE_NIL
	}
}

func ToAirportEffectiveDate(effectiveDate string) *time.Time {
	if effectiveDate == "" {
		return nil
	}

	// convert to time.Time
	t, err := time.Parse("02/01/2006", effectiveDate)
	if err != nil {
		return nil
	}

	// return &t
	return &t
}
