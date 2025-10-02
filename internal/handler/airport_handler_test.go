package handler

import (
	"bytes"
	"context"
	"encoding/json"
	dto "flight-api/internal/dto/airport"
	pagination_dto "flight-api/internal/dto/pagination"
	queryparams "flight-api/internal/dto/query_params"
	response_dto "flight-api/internal/dto/response"
	"flight-api/pkg/logger"
	"flight-api/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- Mock Service ---
type mockAirportService struct {
}

func newMockAiportService() *mockAirportService {
	return &mockAirportService{}
}

func (m *mockAirportService) Create(ctx context.Context, r dto.AirportRequestDto) dto.AirportDto { //nolint:ireturn
	id, _ := uuid.NewRandom()

	respnse := dto.AirportDto{
		ID:         &id,
		SiteNumber: util.Ptr("SN001"),
		ICAOID:     util.Ptr("WIII"),
		Name:       util.Ptr("Soekarno-Hatta International Airport"),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return respnse
}

func (m *mockAirportService) FindAll(ctx context.Context, p queryparams.QueryParams) pagination_dto.PaginationDto { //nolint:ireturn
	id1, _ := uuid.NewRandom()
	id2, _ := uuid.NewRandom()

	data := []dto.AirportDto{
		{ID: &id1, ICAOID: util.Ptr("WIII")},
		{ID: &id2, ICAOID: util.Ptr("WADD")},
	}

	records := make([]interface{}, len(data))
	for i, v := range data {
		records[i] = v
	}

	return pagination_dto.PaginationDto{
		Object:  "pagination",
		Records: records,
		Total:   2,
		Meta: &pagination_dto.PaginationMetaDto{
			Limit: 2,
			Page:  1,
			Next:  false,
		},
	}
}

func (m *mockAirportService) FindByID(ctx context.Context, id string) (dto.AirportDto, error) {
	return dto.AirportDto{}, nil
}

func (m *mockAirportService) Update(ctx context.Context, id string, u dto.AirportUpdateDto) (dto.AirportDto, error) {
	return dto.AirportDto{}, nil
}

func (m *mockAirportService) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockAirportService) GetWeatherCondition(ctx context.Context, code string, name string, query queryparams.QueryParams) (*pagination_dto.PaginationDto, error) {
	return &pagination_dto.PaginationDto{}, nil
}

// --- Tests ---
func TestAirportHandlerCreateMock(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	mockService := newMockAiportService()
	h := NewAirportHandler(mockService, logger)

	reqBody := `{
		"site_number": "SN002",
		"icao_id": "WIII",
		"name": "Soekarno-Hatta International Airport"
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/airports", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Test Create
	h.Create(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response response_dto.ResponseDto
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response.Status)

	data, _ := json.Marshal(response.Data)

	var airport dto.AirportDto
	err = json.Unmarshal(data, &airport)
	assert.NoError(t, err)

	// Check ID (UUID)
	_, err = uuid.Parse(airport.ID.String())
	assert.NoError(t, err)

	// Check data type for the other
	assert.IsType(t, "string", airport.SiteNumber)
	assert.IsType(t, "string", airport.ICAOID)
	assert.IsType(t, "string", airport.Name)
	assert.IsType(t, time.Time{}, airport.CreatedAt)
	assert.IsType(t, time.Time{}, airport.UpdatedAt)

	// Check value
	assert.Equal(t, "SN001", airport.SiteNumber)
	assert.Equal(t, "WIII", airport.ICAOID)
	assert.Equal(t, "Soekarno-Hatta International Airport", airport.Name)
}

func TestAirportHandlerFindAllMock(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	mockService := &mockAirportService{}
	h := NewAirportHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/v1/airports?limit=10&page=1", nil)
	rr := httptest.NewRecorder()

	// Test FindAll
	h.FindAll(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response response_dto.ResponseDto
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response.Status)

	data, _ := json.Marshal(response.Data)
	var pagination pagination_dto.PaginationDto
	err = json.Unmarshal(data, &pagination)
	assert.NoError(t, err)

	// Check Pagination
	assert.Equal(t, "pagination", pagination.Object)
	assert.Equal(t, 2, pagination.Total)
	assert.Equal(t, 2, len(pagination.Records))

	// Check Meta
	assert.Equal(t, 2, pagination.Meta.Limit)
	assert.Equal(t, 1, pagination.Meta.Page)
	assert.Equal(t, false, pagination.Meta.Next)

	// Check Records
	for _, record := range pagination.Records {
		data, _ := json.Marshal(record)
		var airport dto.AirportDto
		err := json.Unmarshal(data, &airport)
		assert.NoError(t, err)

		// Check ID (UUID)
		_, err = uuid.Parse(airport.ID.String())
		assert.NoError(t, err)

		// Check data type for the other
		assert.IsType(t, "string", airport.ICAOID)
	}
}
