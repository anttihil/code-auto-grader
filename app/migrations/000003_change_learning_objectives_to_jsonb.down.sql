-- Revert learning_objectives from JSONB back to TEXT[]
ALTER TABLE nodes
ALTER COLUMN learning_objectives TYPE TEXT[] USING (
    SELECT ARRAY(SELECT jsonb_array_elements_text(learning_objectives))
);

COMMENT ON COLUMN nodes.learning_objectives IS
    'Array of learning goals for this node';
