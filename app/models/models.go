package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ============================================
// User Models
// ============================================

type User struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash *string    `db:"password_hash" json:"-"` // Omit from JSON
	DisplayName  *string    `db:"display_name" json:"display_name,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
}

// ============================================
// Tech Tree Models
// ============================================

type TechTree struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Slug            string    `db:"slug" json:"slug"`
	Description     *string   `db:"description" json:"description,omitempty"`
	Language        string    `db:"language" json:"language"`
	DifficultyLevel *string   `db:"difficulty_level" json:"difficulty_level,omitempty"`
	EstimatedHours  *int      `db:"estimated_hours" json:"estimated_hours,omitempty"`
	IsPublished     bool      `db:"is_published" json:"is_published"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// ============================================
// Node Models
// ============================================

type Node struct {
	ID                 uuid.UUID    `db:"id" json:"id"`
	TechTreeID         uuid.UUID    `db:"tech_tree_id" json:"tech_tree_id"`
	Name               string   `db:"name" json:"name"`
	Slug               string   `db:"slug" json:"slug"`
	Description        *string  `db:"description" json:"description,omitempty"`
	LearningObjectives []string `db:"learning_objectives" json:"learning_objectives"`
	OrderIndex         int      `db:"order_index" json:"order_index"`
	EstimatedMinutes   *int     `db:"estimated_minutes" json:"estimated_minutes,omitempty"`
	CreatedAt          time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time    `db:"updated_at" json:"updated_at"`
}

type NodePrerequisite struct {
	NodeID             uuid.UUID `db:"node_id" json:"node_id"`
	PrerequisiteNodeID uuid.UUID `db:"prerequisite_node_id" json:"prerequisite_node_id"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

// ============================================
// Exercise Models
// ============================================

type Exercise struct {
	ID               uuid.UUID       `db:"id" json:"id"`
	NodeID           uuid.UUID       `db:"node_id" json:"node_id"`
	Title            string          `db:"title" json:"title"`
	Slug             string          `db:"slug" json:"slug"`
	Description      string          `db:"description" json:"description"`
	Difficulty       string          `db:"difficulty" json:"difficulty"`
	EstimatedMinutes *int            `db:"estimated_minutes" json:"estimated_minutes,omitempty"`
	OrderIndex       int             `db:"order_index" json:"order_index"`
	Content          ExerciseContent `db:"content" json:"content"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
}

// ExerciseContent represents the JSONB content field
type ExerciseContent struct {
	StarterCode          string              `json:"starter_code"`
	ReferenceSolution    string              `json:"reference_solution"`
	TestCases            []TestCase          `json:"test_cases"`
	Hints                []Hint              `json:"hints"`
	GradingRubric        GradingRubric       `json:"grading_rubric"`
	AIEvaluationPrompt   string              `json:"ai_evaluation_prompt"`
	CommonMistakes       []CommonMistake     `json:"common_mistakes"`
	ExpectedOutcomes     []string            `json:"expected_outcomes"`
}

type TestCase struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Visible bool   `json:"visible"`
}

type Hint struct {
	Order   int    `json:"order"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GradingRubric struct {
	TestWeight  int                 `json:"test_weight"`    // e.g., 60
	AIWeight    int                 `json:"ai_weight"`      // e.g., 40
	AICriteria  []AICriterion       `json:"ai_criteria"`
}

type AICriterion struct {
	Name        string `json:"name"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

type CommonMistake struct {
	Pattern  string `json:"pattern"`
	Guidance string `json:"guidance"`
}

// Implement sql.Scanner and driver.Valuer for ExerciseContent (JSONB)
func (ec *ExerciseContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into ExerciseContent", value)
	}

	return json.Unmarshal(bytes, ec)
}

func (ec ExerciseContent) Value() (driver.Value, error) {
	return json.Marshal(ec)
}

// ============================================
// Submission Models
// ============================================

// AIFeedback represents structured AI evaluation feedback
type AIFeedback struct {
	Passed       bool     `json:"passed"`
	Score        int      `json:"score"`
	Feedback     string   `json:"feedback"`
	Strengths    []string `json:"strengths"`
	Improvements []string `json:"improvements"`
	Correctness  string   `json:"correctness"`
	CodeQuality  string   `json:"codeQuality"`
	Suggestions  string   `json:"suggestions"`
}

// Implement sql.Scanner and driver.Valuer for AIFeedback (JSONB)
func (af *AIFeedback) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into AIFeedback", value)
	}

	return json.Unmarshal(bytes, af)
}

func (af AIFeedback) Value() (driver.Value, error) {
	return json.Marshal(af)
}

type Submission struct {
	ID                   uuid.UUID    `db:"id" json:"id"`
	UserID               uuid.UUID    `db:"user_id" json:"user_id"`
	ExerciseID           uuid.UUID    `db:"exercise_id" json:"exercise_id"`
	SubmittedCode        string       `db:"submitted_code" json:"submitted_code"`
	Language             string       `db:"language" json:"language"`
	ExecutionOutput      *string      `db:"execution_output" json:"execution_output,omitempty"`
	ExecutionError       *string      `db:"execution_error" json:"execution_error,omitempty"`
	ExecutionTimeMs      *int         `db:"execution_time_ms" json:"execution_time_ms,omitempty"`
	TestResults          *TestResults `db:"test_results" json:"test_results,omitempty"`
	AIFeedback           *AIFeedback  `db:"ai_feedback" json:"ai_feedback,omitempty"`
	AIScore              *int         `db:"ai_score" json:"ai_score,omitempty"`
	AIEvaluationTimeMs   *int         `db:"ai_evaluation_time_ms" json:"ai_evaluation_time_ms,omitempty"`
	TestScore            *int         `db:"test_score" json:"test_score,omitempty"`
	TotalScore           *int         `db:"total_score" json:"total_score,omitempty"`
	CreatedAt            time.Time    `db:"created_at" json:"created_at"`
}

// TestResults represents the JSONB test_results field
type TestResults struct {
	TotalTests  int          `json:"total_tests"`
	PassedTests int          `json:"passed_tests"`
	FailedTests int          `json:"failed_tests"`
	Results     []TestResult `json:"results"`
}

type TestResult struct {
	Name       string  `json:"name"`
	Status     string  `json:"status"` // "passed", "failed"
	Error      *string `json:"error,omitempty"`
	DurationMs *int    `json:"duration_ms,omitempty"`
}

// Implement sql.Scanner and driver.Valuer for TestResults (JSONB)
func (tr *TestResults) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into TestResults", value)
	}

	return json.Unmarshal(bytes, tr)
}

func (tr TestResults) Value() (driver.Value, error) {
	return json.Marshal(tr)
}

// ============================================
// User Progress Models
// ============================================

type UserProgress struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	UserID             uuid.UUID  `db:"user_id" json:"user_id"`
	NodeID             uuid.UUID  `db:"node_id" json:"node_id"`
	Status             string     `db:"status" json:"status"` // "not_started", "in_progress", "completed"
	ExercisesCompleted int        `db:"exercises_completed" json:"exercises_completed"`
	ExercisesTotal     int        `db:"exercises_total" json:"exercises_total"`
	MasteryScore       *float64   `db:"mastery_score" json:"mastery_score,omitempty"`
	MasteryThreshold   float64    `db:"mastery_threshold" json:"mastery_threshold"`
	IsMastered         bool       `db:"is_mastered" json:"is_mastered"`
	StartedAt          *time.Time `db:"started_at" json:"started_at,omitempty"`
	CompletedAt        *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	LastActivityAt     *time.Time `db:"last_activity_at" json:"last_activity_at,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

// UserExerciseProgress tracks per-exercise progress and AI help usage
type UserExerciseProgress struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	ExerciseID    uuid.UUID `db:"exercise_id" json:"exercise_id"`
	AIHelpCount   int       `db:"ai_help_count" json:"ai_help_count"`
	AttemptsCount int       `db:"attempts_count" json:"attempts_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ============================================
// Learning Material Models
// ============================================

// LearningMaterial represents educational content for a tech tree node
type LearningMaterial struct {
	ID         uuid.UUID `db:"id" json:"id"`
	NodeID     uuid.UUID `db:"node_id" json:"node_id"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content" json:"content"`
	OrderIndex int       `db:"order_index" json:"order_index"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ============================================
// Helper Types
// ============================================

// StringArray is a custom type for PostgreSQL TEXT[] arrays
type StringArray []string

// Implement sql.Scanner for StringArray
func (sa *StringArray) Scan(src interface{}) error {
	if src == nil {
		*sa = nil
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return sa.parsePostgresArray(string(v))
	case string:
		return sa.parsePostgresArray(v)
	case []string:
		*sa = v
		return nil
	default:
		return fmt.Errorf("unsupported type for StringArray: %T", src)
	}
}

// parsePostgresArray parses PostgreSQL TEXT[] format: {elem1,elem2,elem3}
func (sa *StringArray) parsePostgresArray(s string) error {
	s = strings.TrimSpace(s)
	if s == "" || s == "{}" {
		*sa = []string{}
		return nil
	}

	// Remove outer braces
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	if s == "" {
		*sa = []string{}
		return nil
	}

	// Split by comma and trim quotes/spaces from each element
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Remove quotes if present
		part = strings.Trim(part, "\"")
		if part != "" {
			result = append(result, part)
		}
	}

	*sa = result
	return nil
}

// Implement driver.Valuer for StringArray
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}

	// Format as PostgreSQL array: {elem1,elem2,elem3}
	// Escape double quotes and wrap elements in quotes
	quoted := make([]string, len(sa))
	for i, s := range sa {
		// Escape quotes and backslashes
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		quoted[i] = `"` + s + `"`
	}

	return "{" + strings.Join(quoted, ",") + "}", nil
}

// ============================================
// View Models (for API responses)
// ============================================

// NodeWithProgress combines Node and UserProgress
type NodeWithProgress struct {
	Node
	Progress *UserProgress `json:"progress,omitempty"`
	IsLocked bool          `json:"is_locked"` // True if prerequisites not met
}

// ExerciseWithHistory includes submission history
type ExerciseWithHistory struct {
	Exercise
	SubmissionHistory []Submission `json:"submission_history"`
	BestScore         *int         `json:"best_score,omitempty"`
}

// TechTreeWithNodes includes all nodes
type TechTreeWithNodes struct {
	TechTree
	Nodes []NodeWithProgress `json:"nodes"`
}

// SubmissionRequest represents a code submission from the client
type SubmissionRequest struct {
	ExerciseID string `json:"exercise_id"`
	Code       string `json:"code"`
	Language   string `json:"language"`
}

// SubmissionResponse is returned after evaluating a submission
type SubmissionResponse struct {
	SubmissionID    uuid.UUID    `json:"submission_id"`
	ExecutionOutput *string      `json:"execution_output,omitempty"`
	ExecutionError  *string      `json:"execution_error,omitempty"`
	TestResults     *TestResults `json:"test_results,omitempty"`
	AIFeedback      *AIFeedback  `json:"ai_feedback,omitempty"`
	TestScore       *int         `json:"test_score,omitempty"`
	AIScore         *int         `json:"ai_score,omitempty"`
	TotalScore      *int         `json:"total_score,omitempty"`
}
