package grading

import (
	"testing"

	"github.com/anttihil/code-auto-grader/models"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name      string
		testScore int
		aiScore   int
		rubric    models.GradingRubric
		expected  int
	}{
		{
			name:      "default weights 60/40",
			testScore: 75,
			aiScore:   85,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  79, // (75*60 + 85*40) / 100 = (4500 + 3400) / 100 = 79
		},
		{
			name:      "custom weights 80/20",
			testScore: 60,
			aiScore:   90,
			rubric:    models.GradingRubric{TestWeight: 80, AIWeight: 20},
			expected:  66, // (60*80 + 90*20) / 100 = (4800 + 1800) / 100 = 66
		},
		{
			name:      "equal weights 50/50",
			testScore: 70,
			aiScore:   90,
			rubric:    models.GradingRubric{TestWeight: 50, AIWeight: 50},
			expected:  80, // (70*50 + 90*50) / 100 = (3500 + 4500) / 100 = 80
		},
		{
			name:      "zero test score",
			testScore: 0,
			aiScore:   100,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  40, // (0*60 + 100*40) / 100 = 40
		},
		{
			name:      "zero AI score",
			testScore: 100,
			aiScore:   0,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  60, // (100*60 + 0*40) / 100 = 60
		},
		{
			name:      "both scores zero",
			testScore: 0,
			aiScore:   0,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  0,
		},
		{
			name:      "both scores 100",
			testScore: 100,
			aiScore:   100,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  100,
		},
		{
			name:      "zero weights defaults to 60/40",
			testScore: 80,
			aiScore:   90,
			rubric:    models.GradingRubric{TestWeight: 0, AIWeight: 0},
			expected:  84, // Should default to (80*60 + 90*40) / 100 = 84
		},
		{
			name:      "weights sum to non-100 gets normalized",
			testScore: 80,
			aiScore:   60,
			rubric:    models.GradingRubric{TestWeight: 30, AIWeight: 20}, // sum = 50
			expected:  72, // (80*30 + 60*20) / 50 = (2400 + 1200) / 50 = 72
		},
		{
			name:      "100% test weight",
			testScore: 85,
			aiScore:   50,
			rubric:    models.GradingRubric{TestWeight: 100, AIWeight: 0},
			expected:  85, // AI score ignored
		},
		{
			name:      "100% AI weight",
			testScore: 50,
			aiScore:   95,
			rubric:    models.GradingRubric{TestWeight: 0, AIWeight: 100},
			expected:  95, // Test score ignored (but weight is 0, so should use default)
		},
		{
			name:      "mastery threshold edge case - just below",
			testScore: 79,
			aiScore:   79,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  79,
		},
		{
			name:      "mastery threshold edge case - exactly at",
			testScore: 80,
			aiScore:   80,
			rubric:    models.GradingRubric{TestWeight: 60, AIWeight: 40},
			expected:  80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateScore(tt.testScore, tt.aiScore, tt.rubric)
			if result != tt.expected {
				t.Errorf("CalculateScore(%d, %d, %+v) = %d, want %d",
					tt.testScore, tt.aiScore, tt.rubric, result, tt.expected)
			}
		})
	}
}

func TestCalculateTestScore(t *testing.T) {
	tests := []struct {
		name        string
		testResults *models.TestResults
		expected    int
	}{
		{
			name: "all tests passed",
			testResults: &models.TestResults{
				TotalTests:  10,
				PassedTests: 10,
				FailedTests: 0,
			},
			expected: 100,
		},
		{
			name: "half tests passed",
			testResults: &models.TestResults{
				TotalTests:  10,
				PassedTests: 5,
				FailedTests: 5,
			},
			expected: 50,
		},
		{
			name: "no tests passed",
			testResults: &models.TestResults{
				TotalTests:  10,
				PassedTests: 0,
				FailedTests: 10,
			},
			expected: 0,
		},
		{
			name: "75% tests passed",
			testResults: &models.TestResults{
				TotalTests:  8,
				PassedTests: 6,
				FailedTests: 2,
			},
			expected: 75, // (6 * 100) / 8 = 75
		},
		{
			name: "rounding down - 2 of 3",
			testResults: &models.TestResults{
				TotalTests:  3,
				PassedTests: 2,
				FailedTests: 1,
			},
			expected: 66, // (2 * 100) / 3 = 66.666... rounds to 66
		},
		{
			name: "single test passed",
			testResults: &models.TestResults{
				TotalTests:  1,
				PassedTests: 1,
				FailedTests: 0,
			},
			expected: 100,
		},
		{
			name: "single test failed",
			testResults: &models.TestResults{
				TotalTests:  1,
				PassedTests: 0,
				FailedTests: 1,
			},
			expected: 0,
		},
		{
			name:        "nil test results",
			testResults: nil,
			expected:    0,
		},
		{
			name: "zero total tests",
			testResults: &models.TestResults{
				TotalTests:  0,
				PassedTests: 0,
				FailedTests: 0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTestScore(tt.testResults)
			if result != tt.expected {
				t.Errorf("CalculateTestScore(%+v) = %d, want %d",
					tt.testResults, result, tt.expected)
			}
		})
	}
}

func TestDeterminePassStatus(t *testing.T) {
	tests := []struct {
		name       string
		totalScore int
		expected   bool
	}{
		{
			name:       "score above threshold",
			totalScore: 90,
			expected:   true,
		},
		{
			name:       "score exactly at threshold (80)",
			totalScore: MasteryThreshold,
			expected:   true,
		},
		{
			name:       "score just below threshold",
			totalScore: 79,
			expected:   false,
		},
		{
			name:       "score at zero",
			totalScore: 0,
			expected:   false,
		},
		{
			name:       "perfect score",
			totalScore: 100,
			expected:   true,
		},
		{
			name:       "score at 50",
			totalScore: 50,
			expected:   false,
		},
		{
			name:       "score at 81 (just above threshold)",
			totalScore: 81,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeterminePassStatus(tt.totalScore)
			if result != tt.expected {
				t.Errorf("DeterminePassStatus(%d) = %v, want %v (threshold=%d)",
					tt.totalScore, result, tt.expected, MasteryThreshold)
			}
		})
	}
}

func TestScaleScoreToWeight(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		weight   int
		expected int
	}{
		{
			name:     "85% score with 40 weight",
			score:    85,
			weight:   40,
			expected: 34, // (85 * 40) / 100 = 34
		},
		{
			name:     "100% score with 60 weight",
			score:    100,
			weight:   60,
			expected: 60,
		},
		{
			name:     "0% score with any weight",
			score:    0,
			weight:   50,
			expected: 0,
		},
		{
			name:     "50% score with 50 weight",
			score:    50,
			weight:   50,
			expected: 25,
		},
		{
			name:     "75% score with 60 weight",
			score:    75,
			weight:   60,
			expected: 45, // (75 * 60) / 100 = 45
		},
		{
			name:     "100% score with 100 weight",
			score:    100,
			weight:   100,
			expected: 100,
		},
		{
			name:     "score with 0 weight",
			score:    80,
			weight:   0,
			expected: 0,
		},
		{
			name:     "rounding down - 66.666...",
			score:    66,
			weight:   50,
			expected: 33, // (66 * 50) / 100 = 33
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScaleScoreToWeight(tt.score, tt.weight)
			if result != tt.expected {
				t.Errorf("ScaleScoreToWeight(%d, %d) = %d, want %d",
					tt.score, tt.weight, result, tt.expected)
			}
		})
	}
}

// TestMasteryThresholdConstant verifies the mastery threshold is set correctly
func TestMasteryThresholdConstant(t *testing.T) {
	if MasteryThreshold != 80 {
		t.Errorf("MasteryThreshold = %d, want 80", MasteryThreshold)
	}
}
