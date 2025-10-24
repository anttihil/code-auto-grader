# Database Schema Design

This document describes the PostgreSQL database schema for the Code Auto-Grader platform.

## Design Principles

1. **Single source of truth**: All learning content stored in database
2. **Flexible content**: Use JSONB for exercise details (easy to evolve)
3. **Performance**: Index foreign keys and frequently queried fields
4. **Audit trail**: Track when entities are created/updated
5. **DAG support**: Node prerequisites form directed acyclic graph
6. **User progress**: Track both exercise-level and node-level mastery

## Entity-Relationship Diagram

```
┌─────────────┐
│   Users     │
└──────┬──────┘
       │
       │ 1:N
       │
       ├─────────────────────────────────────┐
       │                                     │
       ↓                                     ↓
┌──────────────┐                    ┌─────────────────┐
│ User Progress│                    │  Submissions    │
│  (Node-level)│                    │  (Exercise-level)│
└──────┬───────┘                    └────────┬────────┘
       │                                     │
       │ N:1                                 │ N:1
       │                                     │
       ↓                                     ↓
┌──────────────┐                    ┌──────────────┐
│    Nodes     │                    │  Exercises   │
└──────┬───────┘                    └──────┬───────┘
       │                                   │
       │ N:1                               │ N:1
       │                                   │
       ↓                                   ↓
┌──────────────┐                          │
│  Tech Trees  │←─────────────────────────┘
└──────────────┘

      ┌──────────────────┐
      │ Node Prerequisites│  (Many-to-many)
      │  (DAG structure) │
      └──────────────────┘
             ↑    ↓
             │    │
        ┌────┴────┴────┐
        │    Nodes     │
        └──────────────┘
```

## Table Definitions

### 1. users

Stores user accounts (even though MVP has one hard-coded user, design for future).

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),  -- NULL for MVP (no auth yet)
    display_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

**Notes**:
- `password_hash`: NULL for MVP, populate when adding authentication
- Use UUID for distributed-friendly IDs
- `display_name`: For showing in UI (can differ from username)

---

### 2. tech_trees

Represents different learning paths (e.g., "Rust Fundamentals", "Python Basics").

```sql
CREATE TABLE tech_trees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    language VARCHAR(50) NOT NULL,  -- 'rust', 'python', 'javascript'
    difficulty_level VARCHAR(20),   -- 'beginner', 'intermediate', 'advanced'
    estimated_hours INTEGER,        -- Total estimated completion time
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tech_trees_slug ON tech_trees(slug);
CREATE INDEX idx_tech_trees_language ON tech_trees(language);
CREATE INDEX idx_tech_trees_published ON tech_trees(is_published);
```

**Notes**:
- `slug`: URL-friendly identifier (e.g., "rust-fundamentals")
- `language`: Programming language for code execution
- `is_published`: Hide incomplete tech trees from users

---

### 3. nodes

Individual topics within a tech tree (e.g., "Ownership & Borrowing").

```sql
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tech_tree_id UUID NOT NULL REFERENCES tech_trees(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    learning_objectives TEXT[],  -- Array of learning objectives
    order_index INTEGER NOT NULL,  -- Display order within tech tree
    estimated_minutes INTEGER,     -- Time to complete this node
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tech_tree_id, slug)
);

CREATE INDEX idx_nodes_tech_tree ON nodes(tech_tree_id);
CREATE INDEX idx_nodes_order ON nodes(tech_tree_id, order_index);
```

**Notes**:
- `order_index`: For display ordering (Node 1, Node 2, etc.)
- `learning_objectives`: Array of strings for goals
- Each node belongs to exactly one tech tree

---

### 4. node_prerequisites

Defines the DAG structure (which nodes must be completed before others).

```sql
CREATE TABLE node_prerequisites (
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    prerequisite_node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (node_id, prerequisite_node_id),
    CHECK (node_id != prerequisite_node_id)  -- Prevent self-reference
);

CREATE INDEX idx_node_prereqs_node ON node_prerequisites(node_id);
CREATE INDEX idx_node_prereqs_prereq ON node_prerequisites(prerequisite_node_id);
```

**Notes**:
- Many-to-many relationship
- `CHECK` constraint prevents a node from being its own prerequisite
- Application logic should prevent cycles in the DAG

**Example**:
```
Node 3 (Structs) requires:
  - Node 2 (Ownership)
Node 2 (Ownership) requires:
  - Node 1 (Basics)
```

---

### 5. exercises

Individual exercises within a node.

```sql
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20) NOT NULL,  -- 'easy', 'medium', 'hard'
    estimated_minutes INTEGER,
    order_index INTEGER NOT NULL,  -- Order within node (1, 2, 3)

    -- Exercise content (JSONB for flexibility)
    content JSONB NOT NULL,  -- See structure below

    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(node_id, slug)
);

CREATE INDEX idx_exercises_node ON exercises(node_id);
CREATE INDEX idx_exercises_order ON exercises(node_id, order_index);
CREATE INDEX idx_exercises_difficulty ON exercises(difficulty);
```

**Exercise Content Structure (JSONB)**:

```json
{
  "starter_code": "fn main() {\n    // TODO\n}",
  "reference_solution": "fn main() {\n    // Solution\n}",
  "test_cases": [
    {
      "name": "test_basic",
      "code": "#[test]\nfn test_basic() { ... }",
      "visible": true
    }
  ],
  "hints": [
    {
      "order": 1,
      "title": "Hint 1: Basic Syntax",
      "content": "Remember to use..."
    }
  ],
  "grading_rubric": {
    "test_weight": 60,
    "ai_weight": 40,
    "ai_criteria": [
      {
        "name": "Code quality",
        "points": 20,
        "description": "Clean, idiomatic code"
      }
    ]
  },
  "ai_evaluation_prompt": "You are evaluating...",
  "common_mistakes": [
    {
      "pattern": "Using &mut when & would work",
      "guidance": "Explain principle of least privilege"
    }
  ],
  "expected_outcomes": [
    "Understand ownership rules",
    "Use borrowing appropriately"
  ]
}
```

**Why JSONB?**
- Flexible schema (exercise types may evolve)
- Easy to query specific fields (e.g., `content->'test_cases'`)
- No need to create new tables for hints, test cases, etc.

---

### 6. submissions

User attempts at exercises (stores code, results, AI feedback).

```sql
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,

    -- Submission content
    submitted_code TEXT NOT NULL,
    language VARCHAR(50) NOT NULL,  -- 'rust', 'python', 'javascript'

    -- Execution results
    execution_output TEXT,
    execution_error TEXT,
    execution_time_ms INTEGER,
    test_results JSONB,  -- Array of test pass/fail results

    -- AI evaluation
    ai_feedback TEXT,
    ai_score INTEGER,  -- 0-40 (AI portion of grade)
    ai_evaluation_time_ms INTEGER,

    -- Overall score
    test_score INTEGER,  -- 0-60 (test portion of grade)
    total_score INTEGER,  -- 0-100 (test_score + ai_score)

    -- Metadata
    attempt_number INTEGER NOT NULL,  -- 1st, 2nd, 3rd attempt, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_user ON submissions(user_id);
CREATE INDEX idx_submissions_exercise ON submissions(exercise_id);
CREATE INDEX idx_submissions_user_exercise ON submissions(user_id, exercise_id);
CREATE INDEX idx_submissions_created ON submissions(created_at DESC);
```

**Test Results Structure (JSONB)**:

```json
{
  "total_tests": 8,
  "passed_tests": 6,
  "failed_tests": 2,
  "results": [
    {
      "name": "test_basic",
      "status": "passed",
      "duration_ms": 5
    },
    {
      "name": "test_edge_case",
      "status": "failed",
      "error": "assertion failed: expected 42, got 0",
      "duration_ms": 3
    }
  ]
}
```

**Notes**:
- `attempt_number`: Increment for each new submission by same user/exercise
- Stores both raw execution output and structured test results
- AI feedback stored as text (markdown)
- Separate scores for tests and AI evaluation

---

### 7. user_progress

Tracks user progress at the node level (mastery, completion).

```sql
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,

    -- Progress tracking
    status VARCHAR(20) NOT NULL DEFAULT 'not_started',  -- 'not_started', 'in_progress', 'completed'
    exercises_completed INTEGER DEFAULT 0,
    exercises_total INTEGER NOT NULL,  -- Cached from node's exercise count

    -- Mastery calculation
    mastery_score DECIMAL(5,2),  -- Average score across all exercises (0.00-100.00)
    mastery_threshold DECIMAL(5,2) DEFAULT 80.00,  -- Required to mark as "completed"
    is_mastered BOOLEAN DEFAULT false,

    -- Timestamps
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, node_id)
);

CREATE INDEX idx_user_progress_user ON user_progress(user_id);
CREATE INDEX idx_user_progress_node ON user_progress(node_id);
CREATE INDEX idx_user_progress_status ON user_progress(user_id, status);
CREATE INDEX idx_user_progress_mastered ON user_progress(user_id, is_mastered);
```

**Status Transitions**:
- `not_started` → `in_progress`: When user views first exercise
- `in_progress` → `completed`: When `exercises_completed == exercises_total` AND `mastery_score >= mastery_threshold`

**Notes**:
- `exercises_total`: Denormalized for performance (avoid join)
- `mastery_score`: Recalculated when submissions change
- `is_mastered`: TRUE when all exercises completed above threshold

---

## Indexes Summary

**Primary Performance Indexes**:
1. Foreign keys (automatically indexed)
2. Lookup fields (email, username, slug)
3. Filter fields (language, difficulty, status)
4. Sort fields (order_index, created_at)

**Query Patterns Optimized**:
- Get user's progress across all nodes
- Fetch exercises for a node in order
- Find submissions by user+exercise
- Check node prerequisites
- List tech trees by language

---

## Sample Queries

### 1. Get user's current node with next available exercises

```sql
SELECT
    n.id, n.name, n.slug,
    up.status, up.exercises_completed, up.exercises_total,
    up.mastery_score
FROM nodes n
JOIN user_progress up ON up.node_id = n.id
WHERE up.user_id = $1
  AND up.status = 'in_progress'
ORDER BY up.last_activity_at DESC
LIMIT 1;
```

### 2. Get available nodes (prerequisites met)

```sql
WITH user_completed_nodes AS (
    SELECT node_id
    FROM user_progress
    WHERE user_id = $1 AND is_mastered = true
)
SELECT DISTINCT n.*
FROM nodes n
WHERE n.tech_tree_id = $2
  AND NOT EXISTS (
      SELECT 1 FROM node_prerequisites np
      WHERE np.node_id = n.id
        AND np.prerequisite_node_id NOT IN (
            SELECT node_id FROM user_completed_nodes
        )
  );
```

### 3. Get exercise with user's submission history

```sql
SELECT
    e.*,
    json_agg(
        json_build_object(
            'id', s.id,
            'score', s.total_score,
            'created_at', s.created_at,
            'attempt_number', s.attempt_number
        ) ORDER BY s.created_at DESC
    ) FILTER (WHERE s.id IS NOT NULL) as submission_history
FROM exercises e
LEFT JOIN submissions s ON s.exercise_id = e.id AND s.user_id = $1
WHERE e.id = $2
GROUP BY e.id;
```

### 4. Calculate node mastery score

```sql
SELECT
    AVG(s.total_score) as mastery_score,
    COUNT(DISTINCT s.exercise_id) as exercises_completed
FROM submissions s
WHERE s.user_id = $1
  AND s.exercise_id IN (
      SELECT id FROM exercises WHERE node_id = $2
  )
  AND s.id IN (
      -- Get only the best submission per exercise
      SELECT id FROM (
          SELECT id, exercise_id,
                 ROW_NUMBER() OVER (
                     PARTITION BY exercise_id
                     ORDER BY total_score DESC, created_at DESC
                 ) as rn
          FROM submissions
          WHERE user_id = $1
      ) ranked WHERE rn = 1
  );
```

---

## Migration Strategy

### Phase 1: Core Tables (MVP)
1. users
2. tech_trees
3. nodes
4. node_prerequisites
5. exercises

### Phase 2: Progress Tracking
6. submissions
7. user_progress

### Phase 3: Optional (Post-MVP)
- achievements
- user_settings
- feedback_votes (upvote/downvote exercises)
- exercise_analytics (aggregate stats)

---

## Data Integrity

### Constraints

1. **Foreign Key Cascades**:
   - Delete tech tree → delete all nodes, exercises, progress
   - Delete user → delete all submissions, progress
   - Delete node → delete prerequisites, exercises

2. **Check Constraints**:
   - `node_prerequisites`: No self-references
   - `submissions.total_score`: Between 0 and 100
   - `user_progress.mastery_score`: Between 0 and 100

3. **Unique Constraints**:
   - `(user_id, node_id)` in user_progress
   - `(tech_tree_id, slug)` in nodes
   - `(node_id, slug)` in exercises

### Application-Level Validations

1. **DAG Validation**: Prevent cycles in node prerequisites
2. **Score Calculation**: `total_score = test_score + ai_score`
3. **Progress Updates**: Auto-update `user_progress` when submissions change
4. **Attempt Numbering**: Auto-increment `attempt_number` per exercise

---

## Seed Data for Dogfooding MVP

### 1. Create Hard-Coded User

```sql
INSERT INTO users (id, username, email, display_name) VALUES
('00000000-0000-0000-0000-000000000001', 'dogfood_user', 'antti@example.com', 'Antti (Dogfooding)');
```

### 2. Create Rust Tech Tree

```sql
INSERT INTO tech_trees (id, name, slug, description, language, difficulty_level, estimated_hours, is_published) VALUES
('10000000-0000-0000-0000-000000000001',
 'Rust Fundamentals',
 'rust-fundamentals',
 'Learn Rust from basics to ownership, structs, and enums',
 'rust',
 'beginner',
 5,
 true);
```

### 3. Create Nodes

```sql
INSERT INTO nodes (id, tech_tree_id, name, slug, description, order_index, estimated_minutes) VALUES
('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001',
 'Rust Basics', 'rust-basics',
 'Variables, functions, control flow', 1, 90),

('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001',
 'Ownership & Borrowing', 'ownership-borrowing',
 'Rust''s unique memory management', 2, 120),

('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001',
 'Structs & Enums', 'structs-enums',
 'Custom data types and pattern matching', 3, 105);
```

### 4. Create Prerequisites

```sql
-- Node 2 requires Node 1
INSERT INTO node_prerequisites (node_id, prerequisite_node_id) VALUES
('20000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001');

-- Node 3 requires Node 2
INSERT INTO node_prerequisites (node_id, prerequisite_node_id) VALUES
('20000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000002');
```

### 5. Create Exercises

See separate seed script for parsing markdown files and inserting exercises with JSONB content.

---

## Performance Considerations

### Indexes
- All foreign keys indexed
- Composite index on `(user_id, exercise_id)` for submission lookups
- Index on `created_at DESC` for recent submissions

### Query Optimization
- Use `EXPLAIN ANALYZE` to validate query plans
- Consider materialized view for user dashboard (aggregates)
- Cache user_progress in application layer (Redis)

### JSONB Performance
- Index on `content` if querying specific fields: `CREATE INDEX idx_exercises_hints ON exercises USING GIN ((content->'hints'));`
- Keep JSONB documents reasonable size (<1MB)

---

## Next Steps

1. Create migration files (up/down) for each table
2. Implement seed script to parse exercise markdown files
3. Create Go models with `sqlx` tags
4. Write CRUD functions for each entity
5. Test schema with sample data

See `/app/migrations/` for SQL migration files.
