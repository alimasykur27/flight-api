package airport_dto_test

import (
	"encoding/json"
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAirportRequestToAirport(t *testing.T) {
	test := map[string]interface{}{
		"site_number":   "A",
		"icao_id":       "A",
		"faa_id":        "A",
		"iata_id":       "A",
		"name":          "A",
		"type":          "A",
		"status":        true,
		"country":       "A",
		"state":         "A",
		"state_full":    "A",
		"county":        "A",
		"city":          "A",
		"owership":      "A",
		"use":           "A",
		"manager":       "A",
		"manager_phone": "A",
		"latitude":      "A",
		"latitude_sec":  "A",
		"longitude":     "A",
		"longitude_sec": "A",
		"elevation":     int64(10),
		"control_tower": true,
		"unicom":        "A",
		"ctaf":          "A",
	}

	// read json
	jsonBytes, _ := json.Marshal(test)
	var request airport_dto.AirportRequestDto
	err := json.Unmarshal(jsonBytes, &request)
	if err != nil {
		t.Fatal(err)
	}

	airport := airport_dto.AirportRequestToAirport(request)

	assert.IsType(t, model.Airport{}, airport)
	assert.Equal(t, "A", *airport.SiteNumber)
	assert.Equal(t, "A", *airport.ICAOID)
	assert.Equal(t, "A", *airport.FAAID)
	assert.Equal(t, "A", *airport.IATAID)
	assert.Equal(t, "A", *airport.Name)
	assert.Equal(t, "A", *airport.Type)
	assert.Equal(t, true, *airport.Status)
	assert.Equal(t, "A", *airport.Country)
	assert.Equal(t, "A", *airport.State)
	assert.Equal(t, "A", *airport.StateFull)
	assert.Equal(t, "A", *airport.County)
	assert.Equal(t, "A", *airport.City)
	assert.Equal(t, "A", *airport.Ownership)
	assert.Equal(t, "A", *airport.Use)
	assert.Equal(t, "A", *airport.Manager)
	assert.Equal(t, "A", *airport.ManagerPhone)
	assert.Equal(t, "A", *airport.Latitude)
	assert.Equal(t, "A", *airport.LatitudeSec)
	assert.Equal(t, "A", *airport.Longitude)
	assert.Equal(t, "A", *airport.LongitudeSec)
	assert.Equal(t, int64(10), *airport.Elevation)
	assert.Equal(t, true, *airport.ControlTower)
	assert.Equal(t, "A", *airport.Unicom)
	assert.Equal(t, "A", *airport.CTAF)
}
