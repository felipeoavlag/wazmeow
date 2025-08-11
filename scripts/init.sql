-- WazMeow Database Initialization
-- This file initializes the PostgreSQL database for WazMeow

-- Create database if it doesn't exist (this is handled by POSTGRES_DB env var)
-- The database 'wazmeow' will be created automatically by the postgres container

-- Set timezone
SET timezone = 'UTC';

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Note: Tables will be created automatically by Bun ORM migrations
-- when the application starts up. This file is just for any initial
-- database setup that needs to happen before the application runs.

-- You can add any additional database initialization here if needed
-- For example: creating additional users, setting permissions, etc.

SELECT 'WazMeow database initialized successfully' as message;
