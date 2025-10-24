-- Migration: 000001_init_schema
-- Description: Create core database tables for Code Auto-Grader

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================
-- Table: users
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),  -- null for MVP (no auth yet)
    display_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_username ON users (username);

COMMENT ON TABLE users IS
    'User accounts (single hard-coded user for MVP)';
COMMENT ON COLUMN users.password_hash IS
    'null for MVP, populate when adding authentication';

-- ============================================
-- Table: tech_trees
-- ============================================
CREATE TABLE tech_trees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    language VARCHAR(50) NOT NULL,
    difficulty_level VARCHAR(20),
    estimated_hours INTEGER,
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE INDEX idx_tech_trees_slug ON tech_trees (slug);
CREATE INDEX idx_tech_trees_language ON tech_trees (language);
CREATE INDEX idx_tech_trees_published ON tech_trees (is_published);

COMMENT ON TABLE tech_trees IS
    'Learning paths (e.g., Rust Fundamentals, Python Basics)';
COMMENT ON COLUMN tech_trees.slug IS
    'URL-friendly identifier';
COMMENT ON COLUMN tech_trees.is_published IS
    'Hide incomplete tech trees from users';

-- ============================================
-- Table: nodes
-- ============================================
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tech_tree_id UUID NOT NULL
        REFERENCES tech_trees (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    learning_objectives TEXT [],
    order_index INTEGER NOT NULL,
    estimated_minutes INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    UNIQUE (tech_tree_id, slug)
);

CREATE INDEX idx_nodes_tech_tree ON nodes (tech_tree_id);
CREATE INDEX idx_nodes_order ON nodes (tech_tree_id, order_index);

COMMENT ON TABLE nodes IS
    'Individual topics within a tech tree';
COMMENT ON COLUMN nodes.learning_objectives IS
    'Array of learning goals for this node';
COMMENT ON COLUMN nodes.order_index IS
    'Display order within tech tree';

-- ============================================
-- Table: node_prerequisites
-- ============================================
CREATE TABLE node_prerequisites (
    node_id UUID NOT NULL
        REFERENCES nodes (id) ON DELETE CASCADE,
    prerequisite_node_id UUID NOT NULL
        REFERENCES nodes (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    PRIMARY KEY (node_id, prerequisite_node_id),
    CHECK (node_id != prerequisite_node_id)
);

CREATE INDEX idx_node_prereqs_node ON node_prerequisites (node_id);
CREATE INDEX idx_node_prereqs_prereq
    ON node_prerequisites (prerequisite_node_id);

COMMENT ON TABLE node_prerequisites IS
    'DAG structure: which nodes must be completed before others';
COMMENT ON CONSTRAINT node_prerequisites_check
    ON node_prerequisites IS 'Prevent self-reference';

-- ============================================
-- Table: exercises
-- ============================================
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    estimated_minutes INTEGER,
    order_index INTEGER NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    UNIQUE (node_id, slug),
    CHECK (difficulty IN ('easy', 'medium', 'hard'))
);

CREATE INDEX idx_exercises_node ON exercises (node_id);
CREATE INDEX idx_exercises_order ON exercises (node_id, order_index);
CREATE INDEX idx_exercises_difficulty ON exercises (difficulty);

COMMENT ON TABLE exercises IS
    'Individual exercises within a node';
COMMENT ON COLUMN exercises.content IS
    'JSONB: starter_code, reference_solution, test_cases, hints';
COMMENT ON COLUMN exercises.order_index IS
    'Order within node (1, 2, 3)';

-- ============================================
-- Table: submissions
-- ============================================
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL
        REFERENCES exercises (id) ON DELETE CASCADE,

    -- Submission content
    submitted_code TEXT NOT NULL,
    language VARCHAR(50) NOT NULL,

    -- Execution results
    execution_output TEXT,
    execution_error TEXT,
    execution_time_ms INTEGER,
    test_results JSONB,

    -- AI evaluation
    ai_feedback TEXT,
    ai_score INTEGER,
    ai_evaluation_time_ms INTEGER,

    -- Overall score
    test_score INTEGER,
    total_score INTEGER,

    -- Metadata
    attempt_number INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    CHECK (ai_score IS null OR
        (ai_score >= 0 AND ai_score <= 40)),
    CHECK (test_score IS null OR
        (test_score >= 0 AND test_score <= 60)),
    CHECK (total_score IS null OR
        (total_score >= 0 AND total_score <= 100)),
    CHECK (attempt_number > 0)
);

CREATE INDEX idx_submissions_user ON submissions (user_id);
CREATE INDEX idx_submissions_exercise ON submissions (exercise_id);
CREATE INDEX idx_submissions_user_exercise
    ON submissions (user_id, exercise_id);
CREATE INDEX idx_submissions_created ON submissions (created_at DESC);

COMMENT ON TABLE submissions IS
    'User attempts at exercises (code, results, AI feedback)';
COMMENT ON COLUMN submissions.test_results IS
    'JSONB: structured test pass/fail results';
COMMENT ON COLUMN submissions.attempt_number IS
    '1st, 2nd, 3rd attempt, etc.';
COMMENT ON COLUMN submissions.total_score IS
    'test_score (60%) + ai_score (40%)';

-- ============================================
-- Table: user_progress
-- ============================================
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    node_id UUID NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,

    -- Progress tracking
    status VARCHAR(20) NOT NULL DEFAULT 'not_started',
    exercises_completed INTEGER DEFAULT 0,
    exercises_total INTEGER NOT NULL,

    -- Mastery calculation
    mastery_score DECIMAL(5, 2),
    mastery_threshold DECIMAL(5, 2) DEFAULT 80.00,
    is_mastered BOOLEAN DEFAULT false,

    -- Timestamps
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    UNIQUE (user_id, node_id),
    CHECK (status IN ('not_started', 'in_progress', 'completed')),
    CHECK (exercises_completed >= 0),
    CHECK (exercises_completed <= exercises_total),
    CHECK (mastery_score IS null OR
        (mastery_score >= 0 AND mastery_score <= 100)),
    CHECK (mastery_threshold >= 0 AND mastery_threshold <= 100)
);

CREATE INDEX idx_user_progress_user ON user_progress (user_id);
CREATE INDEX idx_user_progress_node ON user_progress (node_id);
CREATE INDEX idx_user_progress_status ON user_progress (user_id, status);
CREATE INDEX idx_user_progress_mastered
    ON user_progress (user_id, is_mastered);

COMMENT ON TABLE user_progress IS
    'Track user progress at node level (mastery, completion)';
COMMENT ON COLUMN user_progress.status IS
    'not_started, in_progress, completed';
COMMENT ON COLUMN user_progress.mastery_score IS
    'Average score across all exercises (0.00-100.00)';
COMMENT ON COLUMN user_progress.is_mastered IS
    'true when all exercises completed above threshold';

-- ============================================
-- Function: Update updated_at timestamp
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables with updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tech_trees_updated_at
    BEFORE UPDATE ON tech_trees
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_nodes_updated_at
    BEFORE UPDATE ON nodes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exercises_updated_at
    BEFORE UPDATE ON exercises
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_progress_updated_at
    BEFORE UPDATE ON user_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON FUNCTION update_updated_at_column() IS
    'Automatically update updated_at timestamp on row modification';
