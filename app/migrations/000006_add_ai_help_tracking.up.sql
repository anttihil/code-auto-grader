-- Migration: 000006_add_ai_help_tracking
-- Description: Add column to track AI help requests per user+exercise

-- Add ai_help_count to user_progress table
-- This tracks how many times a user has requested AI explanations (error help, hints, etc.)
ALTER TABLE user_progress ADD COLUMN ai_help_count INTEGER DEFAULT 0;

-- Add constraint to ensure non-negative count
ALTER TABLE user_progress ADD CONSTRAINT user_progress_ai_help_count_check
    CHECK (ai_help_count >= 0);

COMMENT ON COLUMN user_progress.ai_help_count IS
    'Number of times user requested AI assistance for exercises in this node';
