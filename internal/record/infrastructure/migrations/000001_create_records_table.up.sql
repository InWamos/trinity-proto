-- Create schema for user module
SET statement_timeout = '5s';
SET lock_timeout = '1s';
CREATE SCHEMA IF NOT EXISTS "records";

CREATE TABLE IF NOT EXISTS "records"."telegram_records" (
    id UUID PRIMARY KEY NOT NULL,
    from_user_telegram_id BIGINT
    REFERENCES "records"."telegram_users" (telegram_id),
    in_telegram_chat_id BIGINT NOT NULL,
    message_text TEXT,
    posted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS "records"."telegram_users" (
    id UUID PRIMARY KEY NOT NULL,
    telegram_id BIGINT NOT NULL,
    telegram_user_identity_id
    UUID REFERENCES "records"."telegram_user_identities" (id),
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL,
    UNIQUE (telegram_id, added_by_user)
);

CREATE TABLE IF NOT EXISTS "records"."telegram_user_identities" (
    id UUID PRIMARY KEY NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT,
    username TEXT,
    phone_number TEXT,
    bio TEXT,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL
);
-- squawk-ignore require-concurrent-index-creation
CREATE INDEX IF NOT EXISTS
idx_from_user_telegram_id ON "records".telegram_record (from_user_telegram_id);
