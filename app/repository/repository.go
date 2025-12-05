package repository

import (
	"context"

	"github.com/anttihil/code-auto-grader/models"
)

// TechTreeRepository defines operations for tech tree data access
type TechTreeRepository interface {
	// GetByID retrieves a tech tree by its ID
	GetByID(ctx context.Context, id string) (*models.TechTree, error)

	// GetBySlug retrieves a tech tree by its slug
	GetBySlug(ctx context.Context, slug string) (*models.TechTree, error)

	// List retrieves all published tech trees
	List(ctx context.Context) ([]*models.TechTree, error)

	// Create creates a new tech tree
	Create(ctx context.Context, tree *models.TechTree) error

	// Update updates an existing tech tree
	Update(ctx context.Context, tree *models.TechTree) error
}

// NodeRepository defines operations for node data access
type NodeRepository interface {
	// GetByID retrieves a node by its ID
	GetByID(ctx context.Context, id string) (*models.Node, error)

	// GetByTechTree retrieves all nodes for a tech tree
	GetByTechTree(ctx context.Context, treeID string) ([]*models.Node, error)

	// GetPrerequisites retrieves prerequisite node IDs for a node
	GetPrerequisites(ctx context.Context, nodeID string) ([]string, error)

	// GetPrerequisitesBatch retrieves prerequisites for multiple nodes in a single query
	// Returns a map of nodeID -> []prerequisiteNodeIDs
	GetPrerequisitesBatch(ctx context.Context, nodeIDs []string) (map[string][]string, error)

	// Create creates a new node
	Create(ctx context.Context, node *models.Node) error

	// Update updates an existing node
	Update(ctx context.Context, node *models.Node) error
}

// ExerciseRepository defines operations for exercise data access
type ExerciseRepository interface {
	// GetByID retrieves an exercise by its ID
	GetByID(ctx context.Context, id string) (*models.Exercise, error)

	// GetByNode retrieves all exercises for a node
	GetByNode(ctx context.Context, nodeID string) ([]*models.Exercise, error)

	// GetBySlug retrieves an exercise by node ID and slug
	GetBySlug(ctx context.Context, nodeID, slug string) (*models.Exercise, error)

	// Create creates a new exercise
	Create(ctx context.Context, exercise *models.Exercise) error

	// Update updates an existing exercise
	Update(ctx context.Context, exercise *models.Exercise) error
}

// SubmissionRepository defines operations for submission data access
type SubmissionRepository interface {
	// GetByID retrieves a submission by its ID
	GetByID(ctx context.Context, id string) (*models.Submission, error)

	// GetByExercise retrieves all submissions for an exercise by a user
	GetByExercise(ctx context.Context, userID, exerciseID string) ([]*models.Submission, error)

	// GetByExercisePaginated retrieves submissions for an exercise with pagination
	// Returns submissions, total count, and error
	GetByExercisePaginated(ctx context.Context, userID, exerciseID string, limit, offset int) ([]*models.Submission, int, error)

	// GetLatestAttempt gets the most recent submission for a user+exercise
	GetLatestAttempt(ctx context.Context, userID, exerciseID string) (*models.Submission, error)

	// Create creates a new submission
	Create(ctx context.Context, submission *models.Submission) error

	// CountAttempts counts how many attempts a user has made for an exercise
	CountAttempts(ctx context.Context, userID, exerciseID string) (int, error)
}

// UserProgressRepository defines operations for user progress tracking
type UserProgressRepository interface {
	// GetByUserAndNode retrieves progress for a specific user and node
	GetByUserAndNode(ctx context.Context, userID, nodeID string) (*models.UserProgress, error)

	// GetByUser retrieves all progress records for a user
	GetByUser(ctx context.Context, userID string) ([]*models.UserProgress, error)

	// Upsert creates or updates a user progress record
	Upsert(ctx context.Context, progress *models.UserProgress) error

	// UpdateMastery updates mastery-related fields
	UpdateMastery(ctx context.Context, userID, nodeID string, masteryScore float64, isMastered bool) error
}

// UserRepository defines operations for user data access
type UserRepository interface {
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*models.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error
}

// UserExerciseProgressRepository defines operations for per-exercise progress tracking
type UserExerciseProgressRepository interface {
	// GetByUserAndExercise retrieves the progress by unique user+exercise pair
	GetByUserAndExercise(ctx context.Context, userID string, exerciseID string) (*models.UserExerciseProgress, error)

	// GetByUserAndNode retrieves all exercise progress for a user within a node
	GetByUserAndNode(ctx context.Context, userID string, nodeID string) ([]*models.UserExerciseProgress, error)

	// UpsertProgress creates or updates a user exercise progress record
	UpsertProgress(ctx context.Context, userProgress *models.UserExerciseProgress) error

	// UpdateHelpCounts atomically increments AI help and attempts counters
	UpdateHelpCounts(ctx context.Context, userID string, exerciseID string, addedAICount int, addedAttempts int) error
}

// LearningMaterialRepository defines operations for learning material data access
type LearningMaterialRepository interface {
	// GetByNodeID retrieves all learning materials for a node, ordered by order_index
	GetByNodeID(ctx context.Context, nodeID string) ([]*models.LearningMaterial, error)

	// GetByID retrieves a learning material by its ID
	GetByID(ctx context.Context, id string) (*models.LearningMaterial, error)

	// Create creates a new learning material
	Create(ctx context.Context, material *models.LearningMaterial) error

	// Update updates an existing learning material
	Update(ctx context.Context, material *models.LearningMaterial) error

	// Delete deletes a learning material by its ID
	Delete(ctx context.Context, id string) error
}