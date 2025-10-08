package model

import (
	"time"

	"github.com/google/uuid"
)

type Airport struct {
	ID            *uuid.UUID `db:"id"`
	SiteNumber    *string    `db:"site_number"`
	ICAOID        *string    `db:"icao_id"`
	FAAID         *string    `db:"faa_id"`
	IATAID        *string    `db:"iata_id"`
	Name          *string    `db:"name"`
	Type          *string    `db:"type"`
	Status        *bool      `db:"status"`
	Country       *string    `db:"country"`
	State         *string    `db:"state"`
	StateFull     *string    `db:"state_full"`
	County        *string    `db:"county"`
	City          *string    `db:"city"`
	Ownership     *string    `db:"ownership"`
	Use           *string    `db:"use"`
	Manager       *string    `db:"manager"`
	ManagerPhone  *string    `db:"manager_phone"`
	Latitude      *string    `db:"latitude"`
	LatitudeSec   *string    `db:"latitude_sec"`
	Longitude     *string    `db:"longitude"`
	LongitudeSec  *string    `db:"longitude_sec"`
	Elevation     *int64     `db:"elevation"`
	ControlTower  *bool      `db:"control_tower"`
	Unicom        *string    `db:"unicom"`
	CTAF          *string    `db:"ctaf"`
	EffectiveDate *time.Time `db:"effective_date"`
	SyncStatus    *int64     `db:"sync_status"`
	SyncMessage   *string    `db:"sync_message"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
