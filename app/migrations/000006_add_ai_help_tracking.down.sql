-- Migration: 000006_add_ai_help_tracking (rollback)
-- Description: Remove ai_help_count column from user_progress

ALTER TABLE user_progress DROP CONSTRAINT IF EXISTS user_progress_ai_help_count_check;
ALTER TABLE user_progress DROP COLUMN IF EXISTS ai_help_count;
