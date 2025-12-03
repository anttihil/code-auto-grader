package techtrees

import (
	"github.com/anttihil/code-auto-grader/models"
)

// TechTreeView represents a tech tree with embedded user progress
type TechTreeView struct {
	TechTree       *models.TechTree
	Nodes          []*NodeWithProgress
	TotalNodes     int
	CompletedNodes int
}

// NodeWithProgress represents a node enriched with user progress data
type NodeWithProgress struct {
	Node     *models.Node
	Progress *NodeProgress
}

// NodeProgress represents user progress for a specific node
type NodeProgress struct {
	NodeID             string
	ExercisesCompleted int
	TotalExercises     int
	MasteryAchieved    bool
	IsUnlocked         bool
	BestScores         map[string]int // exerciseID -> best score
}

// NodeView represents a detailed view of a node with exercises and progress
type NodeView struct {
	Node      *models.Node
	Exercises []*models.Exercise
	Progress  *NodeProgress
}
