-- Migration: 000007_add_user_exercise_progress (rollback)
-- Description: Remove user_exercise_progress table

DROP TABLE IF EXISTS user_exercise_progress;
