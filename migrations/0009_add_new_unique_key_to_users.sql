-- +goose Up
-- Drop the global unique constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_mobile_no_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

-- Add app-scoped unique indices
-- We use UNIQUE INDEX instead of CONSTRAINT to easily name them and for consistency with typical PG patterns for composite uniqueness
CREATE UNIQUE INDEX uq_users_app_mobile_no ON users (app_id, mobile_no);
CREATE UNIQUE INDEX uq_users_app_email ON users (app_id, email);

-- +goose Down
-- Remove the app-scoped unique indices
DROP INDEX IF EXISTS uq_users_app_mobile_no;
DROP INDEX IF EXISTS uq_users_app_email;

-- Restore the global unique constraints
-- Note: This will fail if there are duplicate mobile numbers or emails across different apps for different users
ALTER TABLE users ADD CONSTRAINT users_mobile_no_key UNIQUE (mobile_no);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
