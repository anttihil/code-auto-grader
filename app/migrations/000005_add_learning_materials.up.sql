-- Migration: 000005_add_learning_materials
-- Description: Create learning_materials table to separate learning content from exercise instructions

-- ============================================
-- Table: learning_materials
-- ============================================
CREATE TABLE learning_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE INDEX idx_learning_materials_node_id ON learning_materials (node_id);
CREATE INDEX idx_learning_materials_order ON learning_materials (node_id, order_index);

-- Apply update trigger
CREATE TRIGGER update_learning_materials_updated_at
    BEFORE UPDATE ON learning_materials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE learning_materials IS
    'Learning content for tech tree nodes, separate from exercise instructions';
COMMENT ON COLUMN learning_materials.title IS
    'Section title (e.g., "Ownership Basics", "Borrowing Rules")';
COMMENT ON COLUMN learning_materials.content IS
    'Markdown content explaining the concept';
COMMENT ON COLUMN learning_materials.order_index IS
    'Display order within the node (0, 1, 2, ...)';
