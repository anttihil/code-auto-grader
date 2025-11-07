-- Revert learning_objectives from JSONB back to TEXT[]
-- Note: This drops and recreates the column, losing data
-- This is acceptable for test/dev environments
ALTER TABLE nodes DROP COLUMN learning_objectives;
ALTER TABLE nodes ADD COLUMN learning_objectives TEXT[] NOT NULL DEFAULT '{}';

COMMENT ON COLUMN nodes.learning_objectives IS
    'Array of learning goals for this node';
