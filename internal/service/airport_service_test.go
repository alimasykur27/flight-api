package service

import (
	"context"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/enum"
	"flight-api/internal/repository"
	"flight-api/pkg/database"
	"flight-api/pkg/logger"
	"flight-api/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAirport(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	ctx := context.Background()

	cfg, err := config.Load()
	util.PanicIfError(err)

	db, err := database.Connect(cfg.DatabaseURL)
	util.PanicIfError(err)

	validate := util.NewValidator()

	repo := repository.NewAirportRepository(logger)
	service := NewAirportService(repo, db, validate, logger)

	airportReq := airport_dto.AirportRequestDto{
		SiteNumber: util.Ptr("askdj"),
		ICAOID:     util.Ptr("AAAA"),
		Name:       util.Ptr("Ali International Airport"),
		Type:       enum.AIRPORT,
		Status:     util.Ptr(true),
	}

	response := service.Create(ctx, airportReq)

	// Check
	assert.NotNil(t, response)
	assert.Equal(t, response.ICAOID, *airportReq.ICAOID)
	assert.Equal(t, response.Name, *airportReq.Name)
	assert.Equal(t, response.Type, airportReq.Type)
	assert.Equal(t, response.Status, *airportReq.Status)

	// Clean
	// err = repo.Delete(ctx, db, response.ID)
	// assert.Nil(t, err)
}
