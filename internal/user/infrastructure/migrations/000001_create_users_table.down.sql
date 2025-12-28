-- Rollback the whole migration
SET statement_timeout = '30s';
SET lock_timeout = '1s';

-- squawk-ignore ban-drop-table
DROP TABLE IF EXISTS "user".users CASCADE;
DROP TYPE IF EXISTS "user".USER_ROLE;
DROP SCHEMA IF EXISTS "user" CASCADE;
