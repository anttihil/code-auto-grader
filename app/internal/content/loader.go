package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Loader loads content from the content directory
type Loader struct {
	contentDir string
}

// NewLoader creates a new content loader
func NewLoader(contentDir string) *Loader {
	return &Loader{
		contentDir: contentDir,
	}
}

// LoadUsers loads users from users.json
func (l *Loader) LoadUsers() ([]User, error) {
	path := filepath.Join(l.contentDir, "users.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading users.json: %w", err)
	}

	var config struct {
		Users []User `json:"users"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing users.json: %w", err)
	}

	return config.Users, nil
}

// LoadTechTrees loads tech trees from tech-trees.json
func (l *Loader) LoadTechTrees() ([]TechTree, error) {
	path := filepath.Join(l.contentDir, "tech-trees.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tech-trees.json: %w", err)
	}

	var config struct {
		TechTrees []TechTree `json:"tech_trees"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing tech-trees.json: %w", err)
	}

	return config.TechTrees, nil
}

// LoadNodes loads nodes from nodes.json
func (l *Loader) LoadNodes() ([]Node, error) {
	path := filepath.Join(l.contentDir, "nodes.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading nodes.json: %w", err)
	}

	var config struct {
		Nodes []Node `json:"nodes"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing nodes.json: %w", err)
	}

	return config.Nodes, nil
}

// LoadExercises walks the content directory and loads all exercise markdown files
func (l *Loader) LoadExercises() ([]*Exercise, error) {
	var exercises []*Exercise

	exercisesDir := filepath.Join(l.contentDir, "rust") // TODO: Support multiple languages

	err := filepath.Walk(exercisesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Parse exercise
		parsed, err := ParseExerciseMarkdown(path)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		// Convert to Exercise model
		exercise, err := parsed.ToExercise()
		if err != nil {
			return fmt.Errorf("converting %s: %w", path, err)
		}

		exercises = append(exercises, exercise)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking exercises directory: %w", err)
	}

	return exercises, nil
}

// LoadAll loads all content (users, tech trees, nodes, exercises)
func (l *Loader) LoadAll() (*SeedData, []*Exercise, error) {
	users, err := l.LoadUsers()
	if err != nil {
		return nil, nil, fmt.Errorf("loading users: %w", err)
	}

	techTrees, err := l.LoadTechTrees()
	if err != nil {
		return nil, nil, fmt.Errorf("loading tech trees: %w", err)
	}

	nodes, err := l.LoadNodes()
	if err != nil {
		return nil, nil, fmt.Errorf("loading nodes: %w", err)
	}

	exercises, err := l.LoadExercises()
	if err != nil {
		return nil, nil, fmt.Errorf("loading exercises: %w", err)
	}

	seedData := &SeedData{
		Users:     users,
		TechTrees: techTrees,
		Nodes:     nodes,
	}

	return seedData, exercises, nil
}
