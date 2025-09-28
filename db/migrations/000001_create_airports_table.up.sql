-- Tipe enum native PostgreSQL (opsional). Bisa juga pakai CHECK.
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'facility_type') THEN
        CREATE TYPE facility_type AS ENUM ('AIRPORT','HELIPORT');
    END IF;

    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'ownership_type') THEN
        CREATE TYPE ownership_type AS ENUM ('PU','PR');
    END IF;

    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'use_type') THEN
        CREATE TYPE use_type AS ENUM ('PU','PR');
    END IF;
END
$$ LANGUAGE plpgsql;

-- UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.airports (
    id                          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_number                 VARCHAR(255) UNIQUE,                         -- FAA site number
    icao_id                     VARCHAR(10)  NOT NULL UNIQUE,
    faa_id                      VARCHAR(10)  UNIQUE,
    iata_id                     VARCHAR(10)  UNIQUE,
    name                        VARCHAR(255),
    type                        facility_type,                               -- 'AIRPORT'/'HELIPORT'
    status                      BOOLEAN,                                     -- true=open (O), false=closed (C)
    country                     VARCHAR(64),
    state                       VARCHAR(10),
    state_full                  VARCHAR(64),
    county                      VARCHAR(64),
    city                        VARCHAR(64),
    ownership                   ownership_type,                              -- 'PU'/'PR'
    "use"                       use_type,                                    -- 'PU'/'PR' (pakai quotes karena keyword)
    manager                     VARCHAR(255),
    manager_phone               VARCHAR(32),                                 -- bisa mengandung dash/alfanumerik
    latitude                    VARCHAR(64),                                 -- "35-26-04.0000N"
    latitude_sec                VARCHAR(64),                                 -- "127564.0000N"
    longitude                   VARCHAR(64),                                 -- "082-32-33.8240W"
    longitude_sec               VARCHAR(64),                                 -- "297153.8240W"
    elevation                   INTEGER,                                     -- akan diformat string di VIEW
    magnetic_variation          VARCHAR(16),
    tpa                         INTEGER,
    vfr_sectional               VARCHAR(64),
    district_office             VARCHAR(16),                                 -- FSDO (ex: MEM)
    boundary_artcc              VARCHAR(8),
    boundary_artcc_name         VARCHAR(64),
    responsible_artcc           VARCHAR(8),
    responsible_artcc_name      VARCHAR(64),
    fss_phone_number            VARCHAR(32),
    fss_phone_number_tollfree   VARCHAR(32),
    notam_facility_ident        VARCHAR(16),
    certification_typedate      VARCHAR(64),
    customs_airport_of_entry    BOOLEAN,                                     -- true->'Y'
    military_join_use           BOOLEAN,
    military_landing            BOOLEAN,
    lighting_schedule           VARCHAR(64),
    beacon_schedule             VARCHAR(64),
    control_tower               BOOLEAN,
    unicom                      VARCHAR(64),
    ctaf                        VARCHAR(64),
    effective_date              DATE,
    sync_status                 BOOLEAN,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS trg_airports_updated_at ON public.airports;

CREATE TRIGGER trg_airports_updated_at
BEFORE UPDATE ON public.airports
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();