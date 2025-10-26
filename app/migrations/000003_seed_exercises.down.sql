-- Migration: 000003_seed_exercises (DOWN)
-- Description: Remove all seeded exercises

-- Delete exercises in reverse order (Node 3, then Node 2, then Node 1)

-- Node 3 exercises
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000009'; -- Enums and Pattern Matching
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000008'; -- Methods and Self
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000007'; -- Structs Basics

-- Node 2 exercises
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000006'; -- Borrowing Rules
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000005'; -- References and Borrowing
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000004'; -- Ownership and Moves

-- Node 1 exercises
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000003'; -- Control Flow
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000002'; -- Functions
DELETE FROM exercises WHERE id = '30000000-0000-0000-0000-000000000001'; -- Variables
