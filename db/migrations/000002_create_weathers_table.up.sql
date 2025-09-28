-- UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.weather (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id       VARCHAR(10) NOT NULL,
    type             VARCHAR(8),
    raw              TEXT,
    temperature      NUMERIC(5,1),
    dewpoint         NUMERIC(5,1),
    wind             INTEGER,
    wind_velocity    INTEGER,
    visibility       NUMERIC(5,1),
    alt_hg           NUMERIC(5,2),
    alt_mb           NUMERIC(6,2),
    wx               VARCHAR(64),
    auto_report      BOOLEAN,
    category         VARCHAR(8),
    observation_time TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Optional indexes for faster query
CREATE INDEX IF NOT EXISTS idx_weather_station_id ON public.weather(station_id);
CREATE INDEX IF NOT EXISTS idx_weather_obs_time   ON public.weather(observation_time);

-- Trigger auto-update updated_at
CREATE OR REPLACE FUNCTION set_weather_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_weather_updated_at ON public.weather;
CREATE TRIGGER trg_weather_updated_at
BEFORE UPDATE ON public.weather
FOR EACH ROW EXECUTE FUNCTION set_weather_updated_at();
