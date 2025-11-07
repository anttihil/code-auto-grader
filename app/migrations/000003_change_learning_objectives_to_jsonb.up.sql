-- Change learning_objectives from TEXT[] to JSONB
ALTER TABLE nodes ALTER COLUMN learning_objectives DROP DEFAULT;
ALTER TABLE nodes
ALTER COLUMN learning_objectives TYPE JSONB USING learning_objectives::TEXT::JSONB;
ALTER TABLE nodes ALTER COLUMN learning_objectives SET DEFAULT '[]'::JSONB;

COMMENT ON COLUMN nodes.learning_objectives IS
    'JSONB array of learning goals for this node';
