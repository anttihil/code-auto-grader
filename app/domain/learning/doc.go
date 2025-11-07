// Package learning provides the domain logic for tech tree navigation and progress tracking.
//
// This package manages the learning path through a technology tree structure:
//   - Retrieving tech trees with enriched progress data for users
//   - Fetching individual nodes with exercises and best scores
//   - Calculating progress percentages across entire trees
//   - Determining node unlock status based on prerequisites
//   - Aggregating completion statistics and mastery achievements
//
// The package is framework-agnostic and works with domain-specific view models that
// combine data from multiple sources (tech trees, nodes, exercises, progress, submissions).
//
// # Key Components
//
// Service: LearningService provides tech tree navigation functionality
//   - GetTechTree: Returns a complete tech tree with progress for all nodes
//   - GetNode: Returns a single node with exercises and submission history
//   - CalculateProgress: Computes completion percentage for a tree
//
// Models: View models that aggregate data for presentation
//   - TechTreeView: Tech tree with enriched node progress
//   - NodeView: Node with exercises and best scores
//   - NodeProgress: Progress tracking for individual nodes
//   - ExerciseProgress: Exercise-level completion and scores
//
// # Example Usage
//
//	service := learning.NewLearningService(treeRepo, nodeRepo, exerciseRepo, progressRepo, submissionRepo)
//
//	tree, err := service.GetTechTree(ctx, "tree-id", "user-id", "request-id")
//	if err != nil {
//	    // Handle error (e.g., ErrTechTreeNotFound)
//	}
//
//	fmt.Printf("Tree: %s\n", tree.Name)
//	for _, node := range tree.Nodes {
//	    fmt.Printf("  Node: %s - %d/%d exercises completed\n",
//	        node.Title,
//	        node.Progress.ExercisesCompleted,
//	        node.Progress.TotalExercises)
//	}
//
// # Progress Calculation
//
// Progress is tracked at multiple levels:
//   - Exercise level: Best score achieved for each exercise
//   - Node level: Count of mastered exercises (score ≥ 80%)
//   - Tree level: Aggregate of all node progress
//
// Nodes can have prerequisites that must be completed before they become unlocked.
// The IsUnlocked field in NodeProgress indicates whether a user has access to a node.
//
// # Testing
//
// The package uses repository interfaces for all data access, enabling easy testing
// with mocks. See service_test.go for examples of testing navigation and progress
// calculation without a real database.
package learning
