-- Drop trigger 
DROP TRIGGER IF EXISTS trg_airports_updated_at ON public.airports;

-- Drop function
DROP FUNCTION IF EXISTS set_updated_at();

-- Drop table
DROP TABLE IF EXISTS public.airports;

-- Drop enum types
DROP TYPE IF EXISTS facility_type;
DROP TYPE IF EXISTS ownership_type;
DROP TYPE IF EXISTS use_type;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";