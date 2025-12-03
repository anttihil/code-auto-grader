package submissions

import (
	"github.com/anttihil/code-auto-grader/models"
)

// CalculateScore computes the total score based on test results and AI evaluation
// Uses the grading rubric weights (defaults: 60% tests, 40% AI)
func CalculateScore(testScore, aiScore int, rubric models.GradingRubric) int {
	testWeight := rubric.TestWeight
	aiWeight := rubric.AIWeight

	// Normalize weights if they don't sum to 100
	totalWeight := testWeight + aiWeight
	if totalWeight == 0 {
		totalWeight = 100
		testWeight = 60
		aiWeight = 40
	}

	return (testScore*testWeight + aiScore*aiWeight) / totalWeight
}

// CalculateTestScore computes test score as a percentage (0-100)
func CalculateTestScore(testResults *models.TestResults) int {
	if testResults == nil || testResults.TotalTests == 0 {
		return 0
	}
	return (testResults.PassedTests * 100) / testResults.TotalTests
}

// DeterminePassStatus returns whether a submission passes based on the total score
func DeterminePassStatus(totalScore int) bool {
	return totalScore >= MasteryThreshold
}

// ScaleScoreToWeight scales a 0-100 score to a specific weight
// For example, if AI score is 85/100 and AI weight is 40, returns 34
func ScaleScoreToWeight(score, weight int) int {
	return (score * weight) / 100
}
