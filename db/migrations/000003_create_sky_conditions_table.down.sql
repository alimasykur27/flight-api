-- Drop trigger
DROP TRIGGER IF EXISTS trg_sky_condition_updated_at ON public.sky_condition;

-- Drop function
DROP FUNCTION IF EXISTS set_sky_condition_updated_at();

-- Drop index
DROP INDEX IF EXISTS sky_condition_weather_id_idx;

-- Drop table
DROP TABLE IF EXISTS public.sky_condition;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";