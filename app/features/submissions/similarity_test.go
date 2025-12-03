package submissions

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
		{
			name:     "handles CRLF line endings",
			input:    "function add(a, b) {\r\n  return a + b;\r\n}",
			expected: "function add(a, b) {\n  return a + b;\n}",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \n\n   \t\t   \n   ",
			expected: "",
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
		{
			name:     "handles multi-line block comments",
			code:     "function test() {\n  /* First line\n     Second line\n     Third line */\n  return true;\n}",
			language: "javascript",
			expected: "function test() {\n  return true;\n}",
		},
		{
			name:     "handles inline block comments",
			code:     "function test() { /* inline comment */ return true; }",
			language: "javascript",
			expected: "function test() {  return true; }",
		},
		{
			name:     "Go language comments",
			code:     "func add(a, b int) int {\n  // Go comment\n  return a + b\n}",
			language: "go",
			expected: "func add(a, b int) int {\n  return a + b\n}",
		},
		{
			name:     "Rust language comments",
			code:     "fn add(a: i32, b: i32) -> i32 {\n  // Rust comment\n  a + b\n}",
			language: "rust",
			expected: "fn add(a: i32, b: i32) -> i32 {\n  a + b\n}",
		},
		{
			name:     "unknown language - no stripping",
			code:     "# This stays\nsome code",
			language: "unknown",
			expected: "# This stays\nsome code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripComments(tt.code, tt.language)
			if result != tt.expected {
				t.Errorf("stripComments()\ngot:  %q\nwant: %q", result, tt.expected)
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
		{
			name:     "empty strings",
			s1:       "",
			s2:       "",
			minScore: 100.0,
			maxScore: 100.0,
		},
		{
			name:     "one empty string",
			s1:       "some code",
			s2:       "",
			minScore: 0.0,
			maxScore: 0.0,
		},
		{
			name:     "partial match",
			s1:       "line1\nline2\nline3",
			s2:       "line1\nline2\ndifferent",
			minScore: 60.0,
			maxScore: 70.0,
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
	// Create a service instance (checkCodeSimilarity is a method on GradingService)
	service := &GradingService{}

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
		{
			name:      "only block comments added",
			submitted: "function add(a, b) {\n  /* This adds two numbers */\n  return a + b;\n}",
			starter:   "function add(a, b) {\n  return a + b;\n}",
			language:  "javascript",
			expected:  CodeTrivialChanges,
		},
		{
			name:      "Python - only comments added",
			submitted: "def add(a, b):\n    # This adds two numbers\n    return a + b",
			starter:   "def add(a, b):\n    return a + b",
			language:  "python",
			expected:  CodeTrivialChanges,
		},
		{
			name:      "Python - substantial changes",
			submitted: "def add(a, b):\n    result = a + b\n    print(f'Result: {result}')\n    return result",
			starter:   "def add(a, b):\n    return a + b",
			language:  "python",
			expected:  CodeSubstantial,
		},
		{
			name:      "Go - substantial changes (different implementation)",
			submitted: "func add(a, b int) int {\n    return a + b\n}",
			starter:   "func add(a, b int) int {\n    return a + b + 0\n}",
			language:  "go",
			expected:  CodeSubstantial,
		},
		{
			name:      "Rust - minimal whitespace only",
			submitted: "fn add(a: i32, b: i32) -> i32 {\n    a + b\n}",
			starter:   "fn add(a: i32, b: i32) -> i32 {\n  a + b\n}",
			language:  "rust",
			expected:  CodeMinimalChanges,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.checkCodeSimilarity(tt.submitted, tt.starter, tt.language)
			if result != tt.expected {
				t.Errorf("checkCodeSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCodeSimilarityLevelConstants verifies the similarity level constants
func TestCodeSimilarityLevelConstants(t *testing.T) {
	// Verify enum values are sequential
	if CodeIdentical != 0 {
		t.Errorf("CodeIdentical = %d, want 0", CodeIdentical)
	}
	if CodeTrivialChanges != 1 {
		t.Errorf("CodeTrivialChanges = %d, want 1", CodeTrivialChanges)
	}
	if CodeMinimalChanges != 2 {
		t.Errorf("CodeMinimalChanges = %d, want 2", CodeMinimalChanges)
	}
	if CodeSubstantial != 3 {
		t.Errorf("CodeSubstantial = %d, want 3", CodeSubstantial)
	}
}
