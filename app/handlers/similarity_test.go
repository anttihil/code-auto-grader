package handlers

import (
	"testing"
)

func TestNormalizeCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes trailing whitespace",
			input:    "function add(a, b) {  \n  return a + b;  \n}  ",
			expected: "function add(a, b) {\n  return a + b;\n}",
		},
		{
			name:     "normalizes tabs to spaces",
			input:    "function add(a, b) {\n\treturn a + b;\n}",
			expected: "function add(a, b) {\n    return a + b;\n}",
		},
		{
			name:     "removes multiple blank lines",
			input:    "function add(a, b) {\n\n\n  return a + b;\n}",
			expected: "function add(a, b) {\n\n  return a + b;\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCode(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeCode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStripComments(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		language string
		expected string
	}{
		{
			name:     "strips single-line JS comments",
			code:     "function add(a, b) {\n  // This is a comment\n  return a + b;\n}",
			language: "javascript",
			expected: "function add(a, b) {\n  return a + b;\n}",
		},
		{
			name:     "strips block JS comments",
			code:     "function add(a, b) {\n  /* This is a block comment */\n  return a + b;\n}",
			language: "javascript",
			expected: "function add(a, b) {\n  return a + b;\n}",
		},
		{
			name:     "strips Python comments",
			code:     "def add(a, b):\n    # This is a comment\n    return a + b",
			language: "python",
			expected: "def add(a, b):\n    return a + b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripComments(tt.code, tt.language)
			if result != tt.expected {
				t.Errorf("stripComments() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		minScore float64 // Minimum expected similarity
		maxScore float64 // Maximum expected similarity
	}{
		{
			name:     "identical strings",
			s1:       "function add(a, b) { return a + b; }",
			s2:       "function add(a, b) { return a + b; }",
			minScore: 100.0,
			maxScore: 100.0,
		},
		{
			name:     "completely different strings",
			s1:       "function add(a, b) { return a + b; }",
			s2:       "def multiply(x, y): return x * y",
			minScore: 0.0,
			maxScore: 30.0,
		},
		{
			name:     "minor differences",
			s1:       "function add(a, b) {\n  return a + b;\n}",
			s2:       "function add(a, b) {\n  return a + b + 0;\n}",
			minScore: 60.0,
			maxScore: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSimilarity(tt.s1, tt.s2)
			if result < tt.minScore || result > tt.maxScore {
				t.Errorf("calculateSimilarity() = %f, want between %f and %f", result, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestCheckCodeSimilarity(t *testing.T) {
	h := &SubmissionHandler{}

	tests := []struct {
		name      string
		submitted string
		starter   string
		language  string
		expected  CodeSimilarityLevel
	}{
		{
			name:      "identical code",
			submitted: "function add(a, b) { return a + b; }",
			starter:   "function add(a, b) { return a + b; }",
			language:  "javascript",
			expected:  CodeIdentical,
		},
		{
			name:      "only comments added",
			submitted: "// This adds two numbers\nfunction add(a, b) { return a + b; }",
			starter:   "function add(a, b) { return a + b; }",
			language:  "javascript",
			expected:  CodeTrivialChanges,
		},
		{
			name:      "minor whitespace changes",
			submitted: "function add(a, b) {\n  return a + b;\n}",
			starter:   "function add(a, b) {\n return a + b;\n}",
			language:  "javascript",
			expected:  CodeMinimalChanges, // Indentation changes detected as minimal
		},
		{
			name:      "substantial changes",
			submitted: "function add(a, b) {\n  const result = a + b;\n  console.log('Result:', result);\n  return result;\n}",
			starter:   "function add(a, b) { return a + b; }",
			language:  "javascript",
			expected:  CodeSubstantial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.checkCodeSimilarity(tt.submitted, tt.starter, tt.language)
			if result != tt.expected {
				t.Errorf("checkCodeSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}
