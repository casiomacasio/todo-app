CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE refresh_tokens
ALTER COLUMN id SET DEFAULT uuid_generate_v4();