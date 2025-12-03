# Progress Feature

**User progress tracking and mastery calculations**

## Overview

The progress feature tracks user advancement through learning nodes, calculates mastery scores, and determines when users have completed nodes. It's a vertical slice containing models, repository, and service logic.

## Responsibilities

- Track user progress through nodes (exercises completed, mastery score)
- Update progress after successful submissions (≥80% score)
- Calculate and persist mastery status
- Provide progress data for tech tree and node views

## Architecture

```
progress/
├── models.go      # UserProgress type, status constants
├── repository.go  # Database operations (GetByUser, GetByUserAndNode, Upsert)
├── service.go     # Business logic (UpdateAfterSubmission, completion stats)
└── README.md      # This file
```

## Key Concepts

### Mastery Threshold

- **Default**: 80% score required to achieve mastery
- **Configurable**: Per-node via `mastery_threshold` field
- **Purpose**: Ensures users demonstrate sufficient understanding before moving on

### Progress Status

- `not_started` - User hasn't attempted any exercises in the node
- `in_progress` - User has started but not completed the node
- `completed` - User has achieved mastery (score ≥ threshold)

### Exercise Tracking

- `exercises_completed` increments when user scores ≥ mastery threshold
- `exercises_total` is the total number of exercises in the node (TODO: implement counting)
- `mastery_score` is the user's highest score for the node

## API

### Service Methods

```go
// Update progress after a submission
func (s *Service) UpdateAfterSubmission(ctx context.Context,
    userID, nodeID string, score int, requestID string) error

// Retrieve progress for a user and node
func (s *Service) GetByUserAndNode(ctx context.Context,
    userID, nodeID string) (*UserProgress, error)

// Retrieve all progress for a user
func (s *Service) GetByUser(ctx context.Context,
    userID string) ([]*UserProgress, error)

// Calculate completion statistics
func CalculateCompletionStats(progressList []*UserProgress,
    totalNodes int) (completedCount, totalCount int)
```

### Repository Interface

```go
type Repository interface {
    GetByUserAndNode(ctx context.Context, userID, nodeID string) (*UserProgress, error)
    GetByUser(ctx context.Context, userID string) ([]*UserProgress, error)
    Upsert(ctx context.Context, progress *UserProgress) error
    UpdateMastery(ctx context.Context, userID, nodeID string,
        masteryScore float64, isMastered bool) error
}
```

## Usage Example

### From Grading Feature (after evaluating a submission)

```go
import "github.com/anttihil/code-auto-grader/features/progress"

// After calculating submission score
progressSvc := progress.NewService(progressRepo)
err := progressSvc.UpdateAfterSubmission(ctx, userID, nodeID, totalScore, requestID)
if err != nil {
    // Handle error - progress update failure shouldn't fail the submission
    logger.Logger.Warn("failed to update progress", slog.String("error", err.Error()))
}
```

### From Learning Feature (displaying tech tree)

```go
import "github.com/anttihil/code-auto-grader/features/progress"

// Fetch all progress for user
progressSvc := progress.NewService(progressRepo)
progressList, err := progressSvc.GetByUser(ctx, userID)
if err != nil {
    return fmt.Errorf("failed to get user progress: %w", err)
}

// Calculate completion stats
completed, total := progress.CalculateCompletionStats(progressList, len(nodes))
```

## Database Schema

The `user_progress` table structure:

```sql
CREATE TABLE user_progress (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    node_id UUID NOT NULL REFERENCES nodes(id),
    status TEXT NOT NULL,                    -- not_started, in_progress, completed
    exercises_completed INTEGER NOT NULL,
    exercises_total INTEGER NOT NULL,
    mastery_score DOUBLE PRECISION,
    mastery_threshold DOUBLE PRECISION NOT NULL DEFAULT 80.0,
    is_mastered BOOLEAN NOT NULL DEFAULT FALSE,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    last_activity_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, node_id)
);
```

## Dependencies

- **Internal**: `internal/logger` for structured logging
- **External**: `github.com/google/uuid`, `github.com/jackc/pgx/v5`
- **Features**: None (completely isolated)

## Future Enhancements

- [ ] Implement `exercises_total` counting from database
- [ ] Add streak tracking (consecutive days with activity)
- [ ] Add time-to-mastery metrics
- [ ] Support for prerequisite checking logic
- [ ] Batch progress updates for multiple nodes
