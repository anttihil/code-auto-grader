-- Migration: 000002_add_performance_indexes
-- Description: Add missing indexes for common query patterns identified in code review

-- Index for submission scoring queries (top scores per exercise)
CREATE INDEX idx_submissions_exercise_score
    ON submissions (exercise_id, total_score DESC, created_at DESC);

-- Index for user dashboard queries (show in-progress nodes)
CREATE INDEX idx_user_progress_user_activity
    ON user_progress (user_id, status, last_activity_at DESC)
    WHERE status IN ('in_progress', 'not_started');

-- Index for recent submission history (paginated queries)
CREATE INDEX idx_submissions_user_recent
    ON submissions (user_id, created_at DESC);

-- Partial index for active user queries (excluding completed progress)
CREATE INDEX idx_user_progress_active
    ON user_progress (user_id, last_activity_at DESC)
    WHERE status != 'completed';

COMMENT ON INDEX idx_submissions_exercise_score IS
    'Optimizes leaderboard and top score queries per exercise';

COMMENT ON INDEX idx_user_progress_user_activity IS
    'Optimizes dashboard queries for showing active learning nodes';

COMMENT ON INDEX idx_submissions_user_recent IS
    'Optimizes submission history pagination queries';

COMMENT ON INDEX idx_user_progress_active IS
    'Optimizes queries for active learning sessions';
