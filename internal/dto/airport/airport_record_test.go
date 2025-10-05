package airport_dto_test

import (
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/model"
	"flight-api/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToAirportRecordDto(t *testing.T) {
	ID := uuid.New()
	m := model.Airport{
		ID:            &ID,
		SiteNumber:    util.Ptr("A"),
		ICAOID:        util.Ptr("A"),
		FAAID:         util.Ptr("A"),
		IATAID:        util.Ptr("A"),
		Name:          util.Ptr("A"),
		Type:          util.Ptr("A"),
		Status:        util.Ptr(true),
		Country:       util.Ptr("A"),
		State:         util.Ptr("A"),
		StateFull:     util.Ptr("A"),
		County:        util.Ptr("A"),
		City:          util.Ptr("A"),
		Ownership:     util.Ptr("A"),
		Use:           util.Ptr("A"),
		Manager:       util.Ptr("A"),
		ManagerPhone:  util.Ptr("A"),
		Latitude:      util.Ptr("A"),
		LatitudeSec:   util.Ptr("A"),
		Longitude:     util.Ptr("A"),
		LongitudeSec:  util.Ptr("A"),
		Elevation:     util.Ptr(int64(10)),
		ControlTower:  util.Ptr(true),
		Unicom:        util.Ptr("A"),
		CTAF:          util.Ptr("A"),
		EffectiveDate: util.Ptr(time.Date(2025, 10, 1, 15, 30, 0, 0, time.UTC)),
		CreatedAt:     util.Ptr(time.Now()),
		UpdatedAt:     util.Ptr(time.Now()),
	}

	result := airport_dto.ToAirportRecordDto(m)

	assert.IsType(t, airport_dto.AirportRecordDto{}, result)
	assert.Equal(t, ID, *result.ID)
	assert.Equal(t, "airport", *result.Object)
	assert.Equal(t, "A", *result.SiteNumber)
	assert.Equal(t, "A", *result.ICAOID)
	assert.Equal(t, "A", *result.FAAID)
	assert.Equal(t, "A", *result.IATAID)
	assert.Equal(t, "A", *result.Name)
	assert.Equal(t, "A", *result.Type)
	assert.Equal(t, true, *result.Status)
	assert.Equal(t, *m.CreatedAt, result.CreatedAt)
	assert.Equal(t, *m.UpdatedAt, result.UpdatedAt)
}

func TestToAirportRecordDtos(t *testing.T) {
	ID1 := uuid.New()
	ID2 := uuid.New()

	m := []model.Airport{
		{
			ID:            &ID1,
			SiteNumber:    util.Ptr("A"),
			ICAOID:        util.Ptr("A"),
			FAAID:         util.Ptr("A"),
			IATAID:        util.Ptr("A"),
			Name:          util.Ptr("A"),
			Type:          util.Ptr("A"),
			Status:        util.Ptr(true),
			Country:       util.Ptr("A"),
			State:         util.Ptr("A"),
			StateFull:     util.Ptr("A"),
			County:        util.Ptr("A"),
			City:          util.Ptr("A"),
			Ownership:     util.Ptr("A"),
			Use:           util.Ptr("A"),
			Manager:       util.Ptr("A"),
			ManagerPhone:  util.Ptr("A"),
			Latitude:      util.Ptr("A"),
			LatitudeSec:   util.Ptr("A"),
			Longitude:     util.Ptr("A"),
			LongitudeSec:  util.Ptr("A"),
			Elevation:     util.Ptr(int64(10)),
			ControlTower:  util.Ptr(true),
			Unicom:        util.Ptr("A"),
			CTAF:          util.Ptr("A"),
			EffectiveDate: util.Ptr(time.Date(2025, 10, 1, 15, 30, 0, 0, time.UTC)),
			CreatedAt:     util.Ptr(time.Now()),
			UpdatedAt:     util.Ptr(time.Now()),
		},
		{
			ID:            &ID2,
			SiteNumber:    util.Ptr("B"),
			ICAOID:        util.Ptr("B"),
			FAAID:         util.Ptr("B"),
			IATAID:        util.Ptr("B"),
			Name:          util.Ptr("B"),
			Type:          util.Ptr("B"),
			Status:        util.Ptr(false),
			Country:       util.Ptr("B"),
			State:         util.Ptr("B"),
			StateFull:     util.Ptr("B"),
			County:        util.Ptr("B"),
			City:          util.Ptr("B"),
			Ownership:     util.Ptr("B"),
			Use:           util.Ptr("B"),
			Manager:       util.Ptr("B"),
			ManagerPhone:  util.Ptr("B"),
			Latitude:      util.Ptr("B"),
			LatitudeSec:   util.Ptr("B"),
			Longitude:     util.Ptr("B"),
			LongitudeSec:  util.Ptr("B"),
			Elevation:     util.Ptr(int64(11123)),
			ControlTower:  util.Ptr(true),
			Unicom:        util.Ptr("A"),
			CTAF:          util.Ptr("A"),
			EffectiveDate: util.Ptr(time.Date(2025, 10, 2, 11, 30, 0, 0, time.UTC)),
			CreatedAt:     util.Ptr(time.Now().Add(5 * time.Hour)),
			UpdatedAt:     util.Ptr(time.Now().Add(5 * time.Hour)),
		},
	}

	results := airport_dto.ToAirportRecordDtos(m)

	assert.IsType(t, []airport_dto.AirportRecordDto{}, results)

	for i, result := range results {
		assert.IsType(t, airport_dto.AirportRecordDto{}, result)
		assert.Equal(t, *m[i].ID, *result.ID)
		assert.Equal(t, "airport", *result.Object)
		assert.Equal(t, *m[i].SiteNumber, *result.SiteNumber)
		assert.Equal(t, *m[i].ICAOID, *result.ICAOID)
		assert.Equal(t, *m[i].FAAID, *result.FAAID)
		assert.Equal(t, *m[i].IATAID, *result.IATAID)
		assert.Equal(t, *m[i].Name, *result.Name)
		assert.Equal(t, *m[i].Type, *result.Type)
		assert.Equal(t, *m[i].Status, *result.Status)
		assert.Equal(t, *m[i].CreatedAt, result.CreatedAt)
		assert.Equal(t, *m[i].UpdatedAt, result.UpdatedAt)
	}
}
