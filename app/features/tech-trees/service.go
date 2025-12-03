// Package learning contains business logic for tech tree navigation and user progress tracking
package techtrees

import (
	"context"
	"log/slog"

	"github.com/anttihil/code-auto-grader/features/progress"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
)

// LearningService handles tech tree navigation, node access, and progress tracking
type LearningService struct {
	techTreeRepo    repository.TechTreeRepository
	nodeRepo        repository.NodeRepository
	exerciseRepo    repository.ExerciseRepository
	progressService *progress.Service
	submissionRepo  repository.SubmissionRepository
}

// NewLearningService creates a new learning service with required dependencies
func NewLearningService(
	techTreeRepo repository.TechTreeRepository,
	nodeRepo repository.NodeRepository,
	exerciseRepo repository.ExerciseRepository,
	progressService *progress.Service,
	submissionRepo repository.SubmissionRepository,
) *LearningService {
	if techTreeRepo == nil {
		panic("tech tree repository is required")
	}
	if nodeRepo == nil {
		panic("node repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if progressService == nil {
		panic("progress service is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}

	return &LearningService{
		techTreeRepo:    techTreeRepo,
		nodeRepo:        nodeRepo,
		exerciseRepo:    exerciseRepo,
		progressService: progressService,
		submissionRepo:  submissionRepo,
	}
}

// GetTechTree fetches a tech tree with all nodes enriched with user progress
func (s *LearningService) GetTechTree(ctx context.Context, treeID, userID, requestID string) (*TechTreeView, error) {
	logger.Logger.Info("fetching tech tree",
		slog.String("request_id", requestID),
		slog.String("tree_id", treeID),
		slog.String("user_id", userID),
	)

	// Fetch tech tree
	tree, err := s.techTreeRepo.GetByID(ctx, treeID)
	if err != nil {
		logger.Logger.Error("failed to fetch tech tree",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, &TechTreeNotFoundError{TreeID: treeID}
	}

	// Fetch all nodes for this tree
	nodes, err := s.nodeRepo.GetByTechTree(ctx, treeID)
	if err != nil {
		logger.Logger.Error("failed to fetch nodes",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// Fetch user progress for all nodes
	userProgress, err := s.progressService.GetByUser(ctx, userID)
	if err != nil {
		logger.Logger.Warn("failed to fetch user progress, continuing without it",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		userProgress = nil
	}

	// Create progress map for efficient lookup
	progressMap := make(map[string]*models.UserProgress)
	if userProgress != nil {
		for i := range userProgress {
			progressMap[userProgress[i].NodeID.String()] = userProgress[i]
		}
	}

	// Fetch exercise counts for each node
	exerciseCounts := make(map[string]int)
	for _, node := range nodes {
		exercises, err := s.exerciseRepo.GetByNode(ctx, node.ID.String())
		if err != nil {
			logger.Logger.Warn("failed to count exercises for node",
				slog.String("request_id", requestID),
				slog.String("node_id", node.ID.String()),
				slog.String("error", err.Error()),
			)
			exerciseCounts[node.ID.String()] = 0
		} else {
			exerciseCounts[node.ID.String()] = len(exercises)
		}
	}

	// Build nodes with progress
	nodesWithProgress := make([]*NodeWithProgress, 0, len(nodes))
	completedNodes := 0

	for i, node := range nodes {
		var progress *NodeProgress

		if userProg, exists := progressMap[node.ID.String()]; exists {
			// Update exercises_total if not set
			totalExercises := userProg.ExercisesTotal
			if totalExercises == 0 {
				totalExercises = exerciseCounts[node.ID.String()]
			}

			progress = &NodeProgress{
				NodeID:             node.ID.String(),
				ExercisesCompleted: userProg.ExercisesCompleted,
				TotalExercises:     totalExercises,
				MasteryAchieved:    userProg.IsMastered,
				IsUnlocked:         true, // If progress exists, node is unlocked
				BestScores:         make(map[string]int),
			}

			if userProg.IsMastered {
				completedNodes++
			}
		} else {
			// Create default progress for nodes without progress
			exerciseCount := exerciseCounts[node.ID.String()]
			progress = &NodeProgress{
				NodeID:             node.ID.String(),
				ExercisesCompleted: 0,
				TotalExercises:     exerciseCount,
				MasteryAchieved:    false,
				IsUnlocked:         i == 0, // First node is always unlocked
				BestScores:         make(map[string]int),
			}
		}

		nodesWithProgress = append(nodesWithProgress, &NodeWithProgress{
			Node:     node,
			Progress: progress,
		})
	}

	logger.Logger.Info("tech tree fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("node_count", len(nodes)),
		slog.Int("completed_nodes", completedNodes),
	)

	return &TechTreeView{
		TechTree:       tree,
		Nodes:          nodesWithProgress,
		TotalNodes:     len(nodes),
		CompletedNodes: completedNodes,
	}, nil
}

// GetNode fetches a node with its exercises and user progress
func (s *LearningService) GetNode(ctx context.Context, nodeID, userID, requestID string) (*NodeView, error) {
	logger.Logger.Info("fetching node",
		slog.String("request_id", requestID),
		slog.String("node_id", nodeID),
		slog.String("user_id", userID),
	)

	// Fetch node
	node, err := s.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		logger.Logger.Error("failed to fetch node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, &NodeNotFoundError{NodeID: nodeID}
	}

	// Fetch exercises for this node
	exercises, err := s.exerciseRepo.GetByNode(ctx, nodeID)
	if err != nil {
		logger.Logger.Error("failed to fetch exercises",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// Fetch user progress for this node
	userProgress, err := s.progressService.GetByUserAndNode(ctx, userID, nodeID)
	if err != nil {
		logger.Logger.Warn("no progress found for user and node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		userProgress = nil
	}

	// Fetch best scores for each exercise from submissions
	bestScores := make(map[string]int)
	for _, exercise := range exercises {
		submissions, _, err := s.submissionRepo.GetByExercisePaginated(ctx, userID, exercise.ID.String(), 100, 0)
		if err == nil && len(submissions) > 0 {
			// Find best score
			for _, sub := range submissions {
				if sub.TotalScore != nil {
					if bestScore, exists := bestScores[exercise.ID.String()]; !exists || *sub.TotalScore > bestScore {
						bestScores[exercise.ID.String()] = *sub.TotalScore
					}
				}
			}
		}
	}

	// Build progress object
	var progress *NodeProgress
	if userProgress != nil {
		// Update exercises_total if not set
		totalExercises := userProgress.ExercisesTotal
		if totalExercises == 0 {
			totalExercises = len(exercises)
		}

		progress = &NodeProgress{
			NodeID:             nodeID,
			ExercisesCompleted: userProgress.ExercisesCompleted,
			TotalExercises:     totalExercises,
			MasteryAchieved:    userProgress.IsMastered,
			IsUnlocked:         true, // If progress exists, node is unlocked
			BestScores:         bestScores,
		}
	} else {
		// Create default progress
		progress = &NodeProgress{
			NodeID:             nodeID,
			ExercisesCompleted: 0,
			TotalExercises:     len(exercises),
			MasteryAchieved:    false,
			IsUnlocked:         true, // Assume unlocked if viewing it
			BestScores:         bestScores,
		}
	}

	logger.Logger.Info("node fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("exercise_count", len(exercises)),
	)

	return &NodeView{
		Node:      node,
		Exercises: exercises,
		Progress:  progress,
	}, nil
}

// CalculateProgress calculates overall progress for a user across a tech tree
func (s *LearningService) CalculateProgress(ctx context.Context, userID, treeID string) (completedNodes int, totalNodes int, err error) {
	// Fetch all nodes for the tree
	nodes, err := s.nodeRepo.GetByTechTree(ctx, treeID)
	if err != nil {
		return 0, 0, err
	}

	// Fetch user progress
	userProgress, err := s.progressService.GetByUser(ctx, userID)
	if err != nil {
		return 0, len(nodes), nil // Return total but no completed
	}

	// Count completed nodes (mastered)
	completedCount := 0
	for _, progress := range userProgress {
		if progress.IsMastered {
			completedCount++
		}
	}

	return completedCount, len(nodes), nil
}
