package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Airport struct {
	ID            uuid.UUID      `db:"id"`
	SiteNumber    sql.NullString `db:"site_number"`
	ICAOID        string         `db:"icao_id"`
	FAAID         sql.NullString `db:"faa_id"`
	IATAID        sql.NullString `db:"iata_id"`
	Name          sql.NullString `db:"name"`
	Type          sql.NullString `db:"type"`
	Status        sql.NullBool   `db:"status"`
	Country       sql.NullString `db:"country"`
	State         sql.NullString `db:"state"`
	StateFull     sql.NullString `db:"state_full"`
	County        sql.NullString `db:"county"`
	City          sql.NullString `db:"city"`
	Ownership     sql.NullString `db:"ownership"`
	Use           sql.NullString `db:"use"`
	Manager       sql.NullString `db:"manager"`
	ManagerPhone  sql.NullString `db:"manager_phone"`
	Latitude      sql.NullString `db:"latitude"`
	LatitudeSec   sql.NullString `db:"latitude_sec"`
	Longitude     sql.NullString `db:"longitude"`
	LongitudeSec  sql.NullString `db:"longitude_sec"`
	Elevation     sql.NullInt64  `db:"elevation"`
	ControlTower  sql.NullBool   `db:"control_tower"`
	Unicom        sql.NullString `db:"unicom"`
	CTAF          sql.NullString `db:"ctaf"`
	EffectiveDate sql.NullTime   `db:"effective_date"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}
