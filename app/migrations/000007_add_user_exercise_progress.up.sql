-- Migration: 000007_add_user_exercise_progress
-- Description: Add per-exercise progress tracking with AI help count
--              This replaces the node-level ai_help_count for more granular insights

-- Create user_exercise_progress table for per-exercise tracking
CREATE TABLE user_exercise_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises (id) ON DELETE CASCADE,

    -- AI Help tracking
    ai_help_count INTEGER DEFAULT 0,
    attempts_count INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    -- Ensure one record per user+exercise pair
    UNIQUE (user_id, exercise_id),

    -- Ensure non-negative counts
    CONSTRAINT user_exercise_progress_ai_help_count_check CHECK (ai_help_count >= 0),
    CONSTRAINT user_exercise_progress_attempts_count_check CHECK (attempts_count >= 0)
);

-- Indexes for common query patterns
CREATE INDEX idx_user_exercise_progress_user ON user_exercise_progress (user_id);
CREATE INDEX idx_user_exercise_progress_exercise ON user_exercise_progress (exercise_id);
CREATE INDEX idx_user_exercise_progress_user_exercise ON user_exercise_progress (user_id, exercise_id);

COMMENT ON TABLE user_exercise_progress IS
    'Tracks per-exercise progress and AI help usage for each user';
COMMENT ON COLUMN user_exercise_progress.ai_help_count IS
    'Number of times user requested AI assistance for this specific exercise';
