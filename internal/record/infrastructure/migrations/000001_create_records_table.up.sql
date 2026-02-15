-- Create schema for user module
SET statement_timeout = '5s';
SET lock_timeout = '1s';
CREATE SCHEMA IF NOT EXISTS "records";

CREATE TABLE IF NOT EXISTS "records"."telegram_users" (
    id UUID PRIMARY KEY NOT NULL,
    telegram_id BIGINT NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL,
    CONSTRAINT "unique_telegram_id_per_user" UNIQUE (telegram_id, added_by_user)
);

CREATE TABLE IF NOT EXISTS "records"."telegram_records" (
    id UUID PRIMARY KEY NOT NULL,
    message_id BIGINT NOT NULL CONSTRAINT "unique_telegram_message_id" UNIQUE,
    from_telegram_user_id UUID
    REFERENCES "records".telegram_users (id),
    in_telegram_chat_id BIGINT NOT NULL,
    message_text TEXT,
    posted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS "records"."telegram_identities" (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID CONSTRAINT "fk_telegram_identities_user"
    REFERENCES "records".telegram_users (
        id
    ),
    first_name TEXT NOT NULL,
    last_name TEXT,
    username TEXT,
    phone_number TEXT,
    bio TEXT,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID NOT NULL,
    CONSTRAINT "unique_telegram_identity_per_user"
    UNIQUE (user_id, first_name, last_name, username, phone_number, bio)
);
-- squawk-ignore require-concurrent-index-creation
CREATE INDEX IF NOT EXISTS
idx_from_user_telegram_id ON "records"."telegram_records" (message_id);
