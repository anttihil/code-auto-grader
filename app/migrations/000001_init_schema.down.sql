-- Migration Down: 000001_init_schema
-- Description: Drop all tables created in 000001_init_schema.up.sql

-- Drop triggers first
DROP TRIGGER IF EXISTS update_user_progress_updated_at ON user_progress;
DROP TRIGGER IF EXISTS update_exercises_updated_at ON exercises;
DROP TRIGGER IF EXISTS update_nodes_updated_at ON nodes;
DROP TRIGGER IF EXISTS update_tech_trees_updated_at ON tech_trees;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS user_progress;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS node_prerequisites;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS tech_trees;
DROP TABLE IF EXISTS users;

-- Optionally drop extension (be careful if other schemas use it)
-- DROP EXTENSION IF EXISTS "pgcrypto";
