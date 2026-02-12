-- Create schema for user module
SET statement_timeout = '5s';
SET lock_timeout = '1s';
CREATE SCHEMA IF NOT EXISTS "user";

-- Create enum for role column
CREATE TYPE "user"."USER_ROLE" AS ENUM ('user', 'admin');

-- Create users table
CREATE TABLE IF NOT EXISTS "user"."users" (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR NOT NULL UNIQUE,
    display_name VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    user_role "user"."USER_ROLE" NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT username_length CHECK (LENGTH(username) > 0)
);

-- Create indexes
-- squawk-ignore require-concurrent-index-creation
CREATE INDEX IF NOT EXISTS
idx_users_id ON "user".users (id);

-- squawk-ignore require-concurrent-index-creation
CREATE INDEX IF NOT EXISTS
idx_users_username ON "user".users (username);
