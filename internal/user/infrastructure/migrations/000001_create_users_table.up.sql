-- Create schema for user module
CREATE SCHEMA IF NOT EXISTS "user";

-- Create enum for role column
CREATE TYPE "user".user_role AS ENUM ('user', 'admin');

-- Create users table
CREATE TABLE "user".users (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR(32) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT username_length CHECK (LENGTH(username) > 0)
);

-- Create indexes
CREATE INDEX idx_users_id ON "user".users(id);
CREATE INDEX idx_users_username ON "user".users(username);
