-- Convert ai_feedback from TEXT to JSONB
ALTER TABLE submissions
ADD COLUMN ai_feedback_new JSONB;

-- Try to convert existing JSON strings to JSONB
-- For MVP: existing data is markdown text (not JSON), so we'll set to NULL
-- New submissions will have proper JSON structure
UPDATE submissions
SET ai_feedback_new = CASE
    WHEN ai_feedback IS NULL OR ai_feedback = '' THEN NULL
    WHEN ai_feedback::text ~ '^\s*\{' THEN ai_feedback::jsonb  -- Looks like JSON
    ELSE NULL  -- Markdown text from old format, set to NULL
END;

-- Drop the old column and rename the new one
ALTER TABLE submissions DROP COLUMN ai_feedback;
ALTER TABLE submissions RENAME COLUMN ai_feedback_new TO ai_feedback;

-- Add GIN index for efficient JSONB queries
CREATE INDEX idx_submissions_ai_feedback_gin ON submissions USING GIN (ai_feedback);

-- Add comment to document the JSONB structure
COMMENT ON COLUMN submissions.ai_feedback IS 'Structured AI feedback as JSONB: {passed, score, feedback, strengths[], improvements[], correctness, codeQuality, suggestions}';
