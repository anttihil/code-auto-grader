-- Migration: 000009_remove_ai_help_from_user_progress
-- Description: Remove ai_help_count from user_progress (now tracked per-exercise)

ALTER TABLE user_progress DROP CONSTRAINT IF EXISTS user_progress_ai_help_count_check;
ALTER TABLE user_progress DROP COLUMN IF EXISTS ai_help_count;
