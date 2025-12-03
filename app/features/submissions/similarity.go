package submissions

import (
	"strings"
)

// checkCodeSimilarity compares submitted code to starter code and returns similarity level
func (s *GradingService) checkCodeSimilarity(submitted, starter, language string) CodeSimilarityLevel {
	// 1. Normalize both codes (trim whitespace, normalize formatting)
	normalizedSubmitted := normalizeCode(submitted)
	normalizedStarter := normalizeCode(starter)

	// 2. Check exact match
	if normalizedSubmitted == normalizedStarter {
		return CodeIdentical
	}

	// 3. Remove comments and check again
	strippedSubmitted := stripComments(normalizedSubmitted, language)
	strippedStarter := stripComments(normalizedStarter, language)

	if strippedSubmitted == strippedStarter {
		return CodeTrivialChanges
	}

	// 4. Calculate similarity percentage
	similarity := calculateSimilarity(strippedSubmitted, strippedStarter)

	// 5. Classify based on similarity threshold
	if similarity >= 90.0 {
		return CodeMinimalChanges
	}

	return CodeSubstantial
}

// normalizeCode removes leading/trailing whitespace and normalizes formatting
func normalizeCode(code string) string {
	// Trim overall whitespace
	code = strings.TrimSpace(code)

	// Normalize line endings
	code = strings.ReplaceAll(code, "\r\n", "\n")

	// Remove multiple consecutive blank lines
	lines := strings.Split(code, "\n")
	var normalized []string
	previousBlank := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !previousBlank {
				normalized = append(normalized, "")
				previousBlank = true
			}
		} else {
			// Normalize indentation (convert tabs to spaces)
			line = strings.ReplaceAll(line, "\t", "    ")
			normalized = append(normalized, strings.TrimRight(line, " "))
			previousBlank = false
		}
	}

	return strings.Join(normalized, "\n")
}

// stripComments removes comments from code (language-aware)
func stripComments(code, language string) string {
	language = strings.ToLower(language)

	switch language {
	case "javascript", "typescript", "go", "rust", "java", "c", "cpp", "c++":
		return stripCStyleComments(code)
	case "python":
		return stripPythonComments(code)
	default:
		return code
	}
}

// stripCStyleComments removes // and /* */ style comments
func stripCStyleComments(code string) string {
	lines := strings.Split(code, "\n")
	var stripped []string
	inBlockComment := false

	for _, line := range lines {
		result := line

		// Handle block comments
		for {
			if inBlockComment {
				if idx := strings.Index(result, "*/"); idx != -1 {
					result = result[idx+2:]
					inBlockComment = false
				} else {
					result = ""
					break
				}
			}

			if idx := strings.Index(result, "/*"); idx != -1 {
				beforeComment := result[:idx]
				result = result[idx+2:]
				inBlockComment = true

				if endIdx := strings.Index(result, "*/"); endIdx != -1 {
					result = beforeComment + result[endIdx+2:]
					inBlockComment = false
				} else {
					result = beforeComment
					break
				}
			} else {
				break
			}
		}

		// Remove single-line comments
		if idx := strings.Index(result, "//"); idx != -1 {
			result = result[:idx]
		}

		if strings.TrimSpace(result) != "" {
			stripped = append(stripped, result)
		}
	}

	return strings.Join(stripped, "\n")
}

// stripPythonComments removes # style comments
func stripPythonComments(code string) string {
	lines := strings.Split(code, "\n")
	var stripped []string

	for _, line := range lines {
		// Simple approach: remove everything after #
		// Note: doesn't handle # in strings, but good enough for similarity check
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		if strings.TrimSpace(line) != "" {
			stripped = append(stripped, line)
		}
	}

	return strings.Join(stripped, "\n")
}

// calculateSimilarity returns similarity percentage (0-100) using line-based comparison
func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 100.0
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Simple line-based similarity
	lines1 := strings.Split(s1, "\n")
	lines2 := strings.Split(s2, "\n")

	// Count matching lines
	matches := 0
	totalLines := max(len(lines1), len(lines2))

	// Create map of lines in s2 for efficient lookup
	lineMap := make(map[string]int)
	for _, line := range lines2 {
		normalized := strings.TrimSpace(line)
		if normalized != "" {
			lineMap[normalized]++
		}
	}

	// Count matches from s1
	for _, line := range lines1 {
		normalized := strings.TrimSpace(line)
		if normalized != "" && lineMap[normalized] > 0 {
			matches++
			lineMap[normalized]--
		}
	}

	if totalLines == 0 {
		return 100.0
	}

	similarity := (float64(matches) / float64(totalLines)) * 100.0
	return similarity
}
