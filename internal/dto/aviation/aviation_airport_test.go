package aviation_dto_test

import (
	airport_dto "flight-api/internal/dto/airport"
	aviation_dto "flight-api/internal/dto/aviation"
	"flight-api/internal/enum"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestToAirportRequestDto(t *testing.T) {
	source := aviation_dto.AviationAirportDto{
		SiteNumber:              "A",
		Type:                    "AIRPORT",
		FacilityName:            "A",
		FAAIdentifier:           "A",
		ICAOIdentifier:          "A",
		Region:                  "A",
		DistrictOffice:          "A",
		State:                   "A",
		StateFull:               "A",
		County:                  "A",
		City:                    "A",
		Ownership:               "PU",
		Use:                     "PU",
		Manager:                 "A",
		ManagerPhone:            "A",
		Latitude:                "",
		LatitudeSec:             "",
		Longitude:               "",
		LongitudeSec:            "",
		Elevation:               "13",
		MagneticVariation:       "",
		TPA:                     "",
		VFRSectional:            "",
		NotamFacilityIdentifier: "",
		Status:                  "O",
		ControlTower:            "Y",
		UNICOM:                  "A",
		CTAF:                    "A",
		EffectiveDate:           "01/02/2006",
	}

	result := aviation_dto.ToAirportRequestDto(source)

	assert.IsType(t, airport_dto.AirportRequestDto{}, result)
	assert.Equal(t, "A", *result.SiteNumber)
	assert.Equal(t, "A", *result.ICAOID)
	assert.Equal(t, "A", *result.FAAID)
	assert.Equal(t, "A", *result.IATAID)
	assert.Equal(t, "A", *result.Name)
	assert.Equal(t, "airport", *result.Type)
	assert.Equal(t, true, *result.Status)
	if result.Country == nil {
		assert.Nil(t, result.Country)
	} else {
		assert.Equal(t, "A", *result.Country)
	}
	assert.Equal(t, "A", *result.State)
	assert.Equal(t, "A", *result.StateFull)
	assert.Equal(t, "A", *result.County)
	assert.Equal(t, "A", *result.City)
	assert.Equal(t, enum.OWN_PUBLIC, result.Ownership)
	assert.Equal(t, enum.USE_PUBLIC, result.Use)
	assert.Equal(t, "A", *result.Manager)
	assert.Equal(t, "A", *result.ManagerPhone)
	assert.Equal(t, "", *result.Latitude)
	assert.Equal(t, "", *result.LatitudeSec)
	assert.Equal(t, "", *result.Longitude)
	assert.Equal(t, "", *result.LongitudeSec)
	assert.Equal(t, int64(13), *result.Elevation)
	assert.Equal(t, true, *result.ControlTower)
	assert.Equal(t, "A", *result.Unicom)
	assert.Equal(t, "A", *result.CTAF)
	assert.Equal(t, "2006-02-01T00:00:00Z", result.EffectiveDate.Format(time.RFC3339Nano))
}

func TestToAirportStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
		logger   *test.Hook
	}{
		{
			name:     "Open",
			status:   "O",
			expected: true,
		},
		{
			name:     "Closed",
			status:   "C",
			expected: false,
		},
		{
			name:     "Invalid",
			status:   "Invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aviation_dto.ToAirportStatus(tt.status)
			assert.IsType(t, false, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToControlTower(t *testing.T) {
	tests := []struct {
		name         string
		controlTower string
		expected     bool
	}{
		{
			name:         "Y",
			controlTower: "Y",
			expected:     true,
		},
		{
			name:         "N",
			controlTower: "N",
			expected:     false,
		},
		{
			name:         "Invalid - 1",
			controlTower: "Invalid",
			expected:     false,
		},
		{
			name:         "Invalid - 2",
			controlTower: "123123123123",
			expected:     false,
		},
		{
			name:         "Empty",
			controlTower: "",
			expected:     false,
		},
		{
			name:         "Nil",
			controlTower: "",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aviation_dto.ToControlTower(tt.controlTower)
			assert.IsType(t, false, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAirportOwnership(t *testing.T) {
	tests := []struct {
		name      string
		ownership string
		expected  enum.OwnershipEnum
	}{
		{
			name:      "Public",
			ownership: "PU",
			expected:  enum.OWN_PUBLIC,
		},
		{
			name:      "Private",
			ownership: "PR",
			expected:  enum.OWN_PRIVATE,
		},
		{
			name:      "Invalid 1",
			ownership: "Invalid",
			expected:  enum.OWN_NIL,
		},
		{
			name:      "Invalid 2",
			ownership: "alsdhajsh",
			expected:  enum.OWN_NIL,
		},
		{
			name:      "Empty",
			ownership: "",
			expected:  enum.OWN_NIL,
		},
		{
			name:      "Nil",
			ownership: "",
			expected:  enum.OWN_NIL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aviation_dto.ToAirportOwnership(tt.ownership)
			assert.IsType(t, enum.OWN_NIL, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAirportUse(t *testing.T) {
	tests := []struct {
		name     string
		use      string
		expected enum.UseTypeEnum
	}{
		{
			name:     "Public",
			use:      "PU",
			expected: enum.USE_PUBLIC,
		},
		{
			name:     "Private",
			use:      "PR",
			expected: enum.USE_PRIVATE,
		},
		{
			name:     "Invalid",
			use:      "Invalid",
			expected: enum.USE_NIL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aviation_dto.ToAirportUse(tt.use)
			assert.IsType(t, enum.USE_NIL, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAirportEffectiveDate(t *testing.T) {
	tests := []struct {
		name          string
		effectiveDate string
		expected      string
	}{
		{
			name:          "Valid date 1",
			effectiveDate: "01/02/2006",
			expected:      "2006-02-01T00:00:00Z",
		},
		{
			name:          "Valid date 2",
			effectiveDate: "01/01/2025",
			expected:      "2025-01-01T00:00:00Z",
		},
		{
			name:          "Empty date string",
			effectiveDate: "",
			expected:      "",
		},
		{
			name:          "Invalid date format",
			effectiveDate: "2006-01-02",
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aviation_dto.ToAirportEffectiveDate(tt.effectiveDate)
			if tt.expected == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result.Format(time.RFC3339Nano))
			}
		})
	}
}
