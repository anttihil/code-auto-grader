-- Migration: 000005_add_learning_materials (rollback)
-- Description: Remove learning_materials table

DROP TRIGGER IF EXISTS update_learning_materials_updated_at ON learning_materials;
DROP TABLE IF EXISTS learning_materials;
