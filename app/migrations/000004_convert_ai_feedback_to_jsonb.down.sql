-- Revert ai_feedback from JSONB back to TEXT
DROP INDEX IF EXISTS idx_submissions_ai_feedback_gin;

ALTER TABLE submissions
ADD COLUMN ai_feedback_new TEXT;

-- Convert JSONB back to TEXT (as JSON string)
UPDATE submissions
SET ai_feedback_new = ai_feedback::text
WHERE ai_feedback IS NOT NULL;

-- Drop the JSONB column and rename the TEXT column
ALTER TABLE submissions DROP COLUMN ai_feedback;
ALTER TABLE submissions RENAME COLUMN ai_feedback_new TO ai_feedback;
