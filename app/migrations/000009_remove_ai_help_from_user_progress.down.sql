-- Migration: 000009_remove_ai_help_from_user_progress (rollback)
-- Description: Restore ai_help_count to user_progress

ALTER TABLE user_progress ADD COLUMN ai_help_count INTEGER DEFAULT 0;

ALTER TABLE user_progress ADD CONSTRAINT user_progress_ai_help_count_check
    CHECK (ai_help_count >= 0);

COMMENT ON COLUMN user_progress.ai_help_count IS
    'Number of times user requested AI assistance for exercises in this node';
