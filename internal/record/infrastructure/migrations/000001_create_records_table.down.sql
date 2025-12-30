-- squawk-ignore-file ban-drop-table
-- Drop schema for user module
SET statement_timeout = '5s';
SET lock_timeout = '1s';
-- squawk-ignore require-concurrent-index-deletion
DROP INDEX IF EXISTS "records".idx_from_user_telegram_id;
DROP TABLE IF EXISTS "records".telegram_records;
DROP SCHEMA IF EXISTS "records";
