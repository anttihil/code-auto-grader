-- Migration: 000008_cleanup_submissions_constraints
-- Description: Remove attempt_number (now tracked in user_exercise_progress)
--              and remove hardcoded 40/60 score constraints (business logic, not schema)

-- Drop attempt_number column (also removes its CHECK constraint)
ALTER TABLE submissions DROP COLUMN IF EXISTS attempt_number;

-- Drop hardcoded score range constraints
-- These constraints couple the schema to business logic (60% tests, 40% AI)
-- The total_score 0-100 constraint is kept as that's a reasonable invariant

-- Find and drop the ai_score constraint (ai_score <= 40)
DO $$
DECLARE
    constraint_name TEXT;
BEGIN
    SELECT conname INTO constraint_name
    FROM pg_constraint
    WHERE conrelid = 'submissions'::regclass
      AND contype = 'c'
      AND pg_get_constraintdef(oid) LIKE '%ai_score%40%';

    IF constraint_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE submissions DROP CONSTRAINT %I', constraint_name);
    END IF;
END $$;

-- Find and drop the test_score constraint (test_score <= 60)
DO $$
DECLARE
    constraint_name TEXT;
BEGIN
    SELECT conname INTO constraint_name
    FROM pg_constraint
    WHERE conrelid = 'submissions'::regclass
      AND contype = 'c'
      AND pg_get_constraintdef(oid) LIKE '%test_score%60%';

    IF constraint_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE submissions DROP CONSTRAINT %I', constraint_name);
    END IF;
END $$;

-- Add more flexible constraints (just ensure non-negative and <= 100)
ALTER TABLE submissions ADD CONSTRAINT submissions_ai_score_check
    CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 100));

ALTER TABLE submissions ADD CONSTRAINT submissions_test_score_check
    CHECK (test_score IS NULL OR (test_score >= 0 AND test_score <= 100));
