-- Rollback the whole migration
DROP TABLE IF EXISTS "user".users CASCADE;
DROP TYPE IF EXISTS "user".user_role;
DROP SCHEMA IF EXISTS "user" CASCADE;