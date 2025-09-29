-- UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- SKY CONDITION TABLE
CREATE TABLE public.sky_condition (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    weather_id  UUID NOT NULL REFERENCES public.weather(id) ON DELETE CASCADE,
    coverage    VARCHAR(8) NOT NULL,   -- contoh: FEW, SCT, BKN, OVC
    base_agl    INTEGER,               -- altitude dalam feet AGL
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index
CREATE INDEX IF NOT EXISTS sky_condition_weather_id_idx ON public.sky_condition(weather_id);

-- Trigger auto-update updated_at
CREATE OR REPLACE FUNCTION set_sky_condition_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_sky_condition_updated_at ON public.sky_condition;
CREATE TRIGGER trg_sky_condition_updated_at
BEFORE UPDATE ON public.sky_condition
FOR EACH ROW EXECUTE FUNCTION set_sky_condition_updated_at();
