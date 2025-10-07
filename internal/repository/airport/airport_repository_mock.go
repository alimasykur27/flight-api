package repository_airport

import (
	"context"
	"database/sql"
	"flight-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type AirportRepositoryMock struct {
	Mock mock.Mock
}

func (r *AirportRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error) {
	args := r.Mock.Called(ctx, tx, airport)
	var out model.Airport
	if v, ok := args.Get(0).(model.Airport); ok {
		out = v
	}
	return out, args.Error(1)
}

func (r *AirportRepositoryMock) SyncAirport(ctx context.Context, tx *sql.Tx, airport model.Airport) (model.Airport, error) {
	args := r.Mock.Called(ctx, tx, airport)
	var out model.Airport
	if v, ok := args.Get(0).(model.Airport); ok {
		out = v
	}
	return out, args.Error(1)
}

func (r *AirportRepositoryMock) FindAll(ctx context.Context, tx *sql.Tx, args map[string]interface{}) ([]model.Airport, int, error) {
	call := r.Mock.Called(ctx, tx, args)

	var list []model.Airport
	if v, ok := call.Get(0).([]model.Airport); ok {
		list = v
	}
	total := 0
	if v, ok := call.Get(1).(int); ok {
		total = v
	}
	return list, total, call.Error(2)
}

func (r *AirportRepositoryMock) FindBySearchName(ctx context.Context, tx *sql.Tx, name string, args map[string]interface{}) ([]model.Airport, int, error) {
	call := r.Mock.Called(ctx, tx, name, args)

	var list []model.Airport
	if v, ok := call.Get(0).([]model.Airport); ok {
		list = v
	}

	total := 0
	if v, ok := call.Get(1).(int); ok {
		total = v
	}

	return list, total, call.Error(2)
}

func (r *AirportRepositoryMock) FindByID(ctx context.Context, tx *sql.Tx, id string) (model.Airport, error) {
	call := r.Mock.Called(ctx, tx, id)
	var out model.Airport
	if v, ok := call.Get(0).(model.Airport); ok {
		out = v
	}
	return out, call.Error(1)
}

func (r *AirportRepositoryMock) FindExistsByICAOID(ctx context.Context, tx *sql.Tx, icaoId string) (bool, error) {
	args := r.Mock.Called(ctx, tx, icaoId)
	var exists bool
	if v, ok := args.Get(0).(bool); ok {
		exists = v
	}
	return exists, args.Error(1)
}

func (r *AirportRepositoryMock) FindByICAOID(ctx context.Context, tx *sql.Tx, icaoId string) (model.Airport, error) {
	call := r.Mock.Called(ctx, tx, icaoId)
	var out model.Airport
	if v, ok := call.Get(0).(model.Airport); ok {
		out = v
	}
	return out, call.Error(1)
}

func (r *AirportRepositoryMock) Update(ctx context.Context, tx *sql.Tx, id string, airport model.Airport) (model.Airport, error) {
	call := r.Mock.Called(ctx, tx, id, airport)
	var out model.Airport
	if v, ok := call.Get(0).(model.Airport); ok {
		out = v
	}
	return out, call.Error(1)
}

func (r *AirportRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	call := r.Mock.Called(ctx, tx, id)
	return call.Error(0)
}
