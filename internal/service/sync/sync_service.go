package service_sync

import (
	"context"
	sync_dto "flight-api/internal/dto/sync"
)

type ISyncService interface {
	SyncAirports(ctx context.Context, req sync_dto.SyncAirportRequest) ([]sync_dto.SyncAirportResponse, error)
}
