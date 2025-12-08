-- Rollback the whole migration
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS user_role;