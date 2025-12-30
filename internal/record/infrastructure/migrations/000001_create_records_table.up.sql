-- Create schema for user module
SET statement_timeout = '5s';
SET lock_timeout = '1s';
CREATE SCHEMA IF NOT EXISTS "records";

CREATE TABLE IF NOT EXISTS "records"."telegram_records" (
    id UUID PRIMARY KEY NOT NULL,
    from_user_telegram_id BIGINT NOT NULL,
    in_telegram_chat_id BIGINT NOT NULL,
    message_text TEXT,
    attachments UUID,
    posted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by_user UUID
);

-- squawk-ignore require-concurrent-index-creation
CREATE INDEX IF NOT EXISTS
idx_from_user_telegram_id ON "records".telegram_record (from_user_telegram_id);
