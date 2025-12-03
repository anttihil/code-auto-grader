package techtrees

import "fmt"

// PrerequisiteError indicates that a node's prerequisites are not met
type PrerequisiteError struct {
	NodeID       string
	MissingNodes []string
}

func (e *PrerequisiteError) Error() string {
	return fmt.Sprintf("prerequisites not met for node %s: missing nodes %v", e.NodeID, e.MissingNodes)
}

// NodeNotFoundError indicates that a requested node was not found
type NodeNotFoundError struct {
	NodeID string
}

func (e *NodeNotFoundError) Error() string {
	return fmt.Sprintf("node not found: %s", e.NodeID)
}

// TechTreeNotFoundError indicates that a requested tech tree was not found
type TechTreeNotFoundError struct {
	TreeID string
}

func (e *TechTreeNotFoundError) Error() string {
	return fmt.Sprintf("tech tree not found: %s", e.TreeID)
}
