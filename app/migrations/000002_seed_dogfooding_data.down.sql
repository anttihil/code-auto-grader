-- Migration Down: 000002_seed_dogfooding_data
-- Description: Remove seed data for dogfooding MVP

-- Delete in reverse order to respect foreign keys

-- Delete user progress
DELETE FROM user_progress
WHERE user_id = '00000000-0000-0000-0000-000000000001';

-- Delete node prerequisites
DELETE FROM node_prerequisites
WHERE node_id IN (
    '20000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003'
);

-- Delete exercises (if any were created)
DELETE FROM exercises
WHERE node_id IN (
    '20000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003'
);

-- Delete nodes
DELETE FROM nodes
WHERE tech_tree_id = '10000000-0000-0000-0000-000000000001';

-- Delete tech tree
DELETE FROM tech_trees
WHERE id = '10000000-0000-0000-0000-000000000001';

-- Delete user
DELETE FROM users
WHERE id = '00000000-0000-0000-0000-000000000001';
