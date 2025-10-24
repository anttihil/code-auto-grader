-- Migration: 000002_seed_dogfooding_data
-- Description: Insert initial data for dogfooding MVP

-- ============================================
-- Seed User (Hard-coded for MVP)
-- ============================================
INSERT INTO users (
    id,
    username,
    email,
    display_name
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'dogfood_user',
    'antti@example.com',
    'Antti (Dogfooding)'
);

-- ============================================
-- Seed Tech Tree: Rust Fundamentals
-- ============================================
INSERT INTO tech_trees (
    id,
    name,
    slug,
    description,
    language,
    difficulty_level,
    estimated_hours,
    is_published
) VALUES (
    '10000000-0000-0000-0000-000000000001',
    'Rust Fundamentals',
    'rust-fundamentals',
    'Learn Rust from basics through ownership, structs, and enums. Build a solid foundation in Rust''s unique memory management and type system.',
    'rust',
    'beginner',
    5,
    true
);

-- ============================================
-- Seed Nodes
-- ============================================

-- Node 1: Rust Basics
INSERT INTO nodes (
    id,
    tech_tree_id,
    name,
    slug,
    description,
    learning_objectives,
    order_index,
    estimated_minutes
) VALUES (
    '20000000-0000-0000-0000-000000000001',
    '10000000-0000-0000-0000-000000000001',
    'Rust Basics',
    'rust-basics',
    'Master Rust fundamentals: variables, mutability, functions, and control flow. Learn the syntax and basic patterns that form the foundation of Rust programming.',
    ARRAY[
        'Understand Rust variable binding and mutability',
        'Write functions with type annotations',
        'Use control flow constructs (if/else, match)',
        'Apply pattern matching to simple types'
    ],
    1,
    90
);

-- Node 2: Ownership & Borrowing
INSERT INTO nodes (
    id,
    tech_tree_id,
    name,
    slug,
    description,
    learning_objectives,
    order_index,
    estimated_minutes
) VALUES (
    '20000000-0000-0000-0000-000000000002',
    '10000000-0000-0000-0000-000000000001',
    'Ownership & Borrowing',
    'ownership-borrowing',
    'Dive into Rust''s unique ownership system. Understand memory management without garbage collection through ownership, borrowing, and lifetimes.',
    ARRAY[
        'Understand Rust''s ownership rules',
        'Recognize when ownership moves occur',
        'Use references to borrow data',
        'Apply borrowing rules (mutable XOR immutable)',
        'Prevent data races at compile time'
    ],
    2,
    120
);

-- Node 3: Structs & Enums
INSERT INTO nodes (
    id,
    tech_tree_id,
    name,
    slug,
    description,
    learning_objectives,
    order_index,
    estimated_minutes
) VALUES (
    '20000000-0000-0000-0000-000000000003',
    '10000000-0000-0000-0000-000000000001',
    'Structs & Enums',
    'structs-enums',
    'Create custom data types with structs and enums. Learn how to model domain concepts, implement methods, and use pattern matching effectively.',
    ARRAY[
        'Define custom structs with named fields',
        'Implement methods and associated functions',
        'Create enums with multiple variants',
        'Use pattern matching with match expressions',
        'Work with Option<T> for nullable values'
    ],
    3,
    105
);

-- ============================================
-- Seed Node Prerequisites (DAG Structure)
-- ============================================

-- Node 2 (Ownership) requires Node 1 (Basics)
INSERT INTO node_prerequisites (
    node_id,
    prerequisite_node_id
) VALUES (
    '20000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000001'
);

-- Node 3 (Structs) requires Node 2 (Ownership)
INSERT INTO node_prerequisites (
    node_id,
    prerequisite_node_id
) VALUES (
    '20000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000002'
);

-- ============================================
-- Initialize User Progress for All Nodes
-- ============================================

-- Node 1: Ready to start (no prerequisites)
INSERT INTO user_progress (
    user_id,
    node_id,
    status,
    exercises_completed,
    exercises_total,
    mastery_threshold
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    'not_started',
    0,
    3,  -- 3 exercises in Node 1
    80.00
);

-- Node 2: Locked (requires Node 1)
INSERT INTO user_progress (
    user_id,
    node_id,
    status,
    exercises_completed,
    exercises_total,
    mastery_threshold
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000002',
    'not_started',
    0,
    3,  -- 3 exercises in Node 2
    80.00
);

-- Node 3: Locked (requires Node 2)
INSERT INTO user_progress (
    user_id,
    node_id,
    status,
    exercises_completed,
    exercises_total,
    mastery_threshold
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000003',
    'not_started',
    0,
    3,  -- 3 exercises in Node 3
    80.00
);

-- ============================================
-- Note: Exercise data will be seeded separately
-- ============================================
-- Run `go run scripts/seed_exercises.go` to parse
-- markdown files and insert exercise content
