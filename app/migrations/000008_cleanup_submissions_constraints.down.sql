-- Migration: 000008_cleanup_submissions_constraints (rollback)
-- Description: Restore attempt_number and hardcoded score constraints

-- Drop the flexible constraints
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS submissions_ai_score_check;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS submissions_test_score_check;

-- Restore attempt_number column
ALTER TABLE submissions ADD COLUMN attempt_number INTEGER;

-- Set default value for existing rows (based on submission order per user+exercise)
UPDATE submissions s
SET attempt_number = sub.rn
FROM (
    SELECT id, ROW_NUMBER() OVER (
        PARTITION BY user_id, exercise_id
        ORDER BY created_at
    ) as rn
    FROM submissions
) sub
WHERE s.id = sub.id;

-- Make it NOT NULL after populating
ALTER TABLE submissions ALTER COLUMN attempt_number SET NOT NULL;

-- Restore the original hardcoded constraints
ALTER TABLE submissions ADD CONSTRAINT submissions_attempt_number_check
    CHECK (attempt_number > 0);

-- Note: Original unnamed constraints had 40/60 limits
-- Adding them back with explicit names for clarity
ALTER TABLE submissions ADD CONSTRAINT submissions_ai_score_range_check
    CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 40));

ALTER TABLE submissions ADD CONSTRAINT submissions_test_score_range_check
    CHECK (test_score IS NULL OR (test_score >= 0 AND test_score <= 60));
