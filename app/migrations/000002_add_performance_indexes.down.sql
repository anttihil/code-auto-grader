-- Migration Down: 000002_add_performance_indexes
-- Description: Remove performance indexes

DROP INDEX IF EXISTS idx_user_progress_active;
DROP INDEX IF EXISTS idx_submissions_user_recent;
DROP INDEX IF EXISTS idx_user_progress_user_activity;
DROP INDEX IF EXISTS idx_submissions_exercise_score;
