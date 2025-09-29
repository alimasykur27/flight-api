-- Drop trigger
DROP TRIGGER IF EXISTS trg_weather_updated_at ON public.weather;

-- Drop function
DROP FUNCTION IF EXISTS set_weather_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_weather_station_id;
DROP INDEX IF EXISTS idx_weather_obs_time;

-- Drop table
DROP TABLE IF EXISTS public.weather;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";