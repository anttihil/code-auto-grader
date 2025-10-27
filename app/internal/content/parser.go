package content

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ExerciseFrontmatter represents the JSON frontmatter in exercise markdown files
type ExerciseFrontmatter struct {
	ID               string `json:"id"`
	NodeID           string `json:"node_id"`
	Title            string `json:"title"`
	Slug             string `json:"slug"`
	Difficulty       string `json:"difficulty"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	OrderIndex       int    `json:"order_index"`
}

// ParsedExercise represents a fully parsed exercise with all sections
type ParsedExercise struct {
	Frontmatter ExerciseFrontmatter
	Description string
	StarterCode string
	Solution    string
	TestCases   string
	Hints       []string
	Outcomes    []string
	Mistakes    []string
	AIPrompt    string
}

// ParseExerciseMarkdown parses an exercise markdown file with JSON frontmatter
func ParseExerciseMarkdown(filePath string) (*ParsedExercise, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// Split frontmatter and body
	parts := strings.SplitN(string(content), "---", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid frontmatter format in %s", filePath)
	}

	// Parse JSON frontmatter
	var frontmatter ExerciseFrontmatter
	if err := json.Unmarshal([]byte(strings.TrimSpace(parts[0])), &frontmatter); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// Parse markdown body
	body := strings.TrimSpace(parts[1])
	sections := extractSections(body)

	// Extract individual components
	exercise := &ParsedExercise{
		Frontmatter: frontmatter,
		Description: extractDescription(sections),
		StarterCode: extractCodeBlock(sections["Starting Code"]),
		Solution:    extractCodeBlock(sections["Reference Solution"]),
		TestCases:   extractCodeBlock(sections["Test Cases"]),
		Hints:       extractHints(sections),
		Outcomes:    extractOutcomes(sections),
		Mistakes:    extractMistakes(sections),
		AIPrompt:    extractAIPrompt(sections),
	}

	return exercise, nil
}

// extractSections splits markdown by ## headers
func extractSections(markdown string) map[string]string {
	sections := make(map[string]string)

	// Regex to match ## Section Name
	re := regexp.MustCompile(`(?m)^## (.+)$`)
	matches := re.FindAllStringSubmatchIndex(markdown, -1)

	for i, match := range matches {
		sectionName := markdown[match[2]:match[3]]

		// Content starts after the header line
		start := match[1]

		// Content ends at next header or EOF
		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(markdown)
		}

		content := strings.TrimSpace(markdown[start:end])
		sections[sectionName] = content
	}

	return sections
}

// extractDescription gets the Problem Description section
func extractDescription(sections map[string]string) string {
	if desc, ok := sections["Problem Description"]; ok {
		return desc
	}
	if desc, ok := sections["Description"]; ok {
		return desc
	}
	return ""
}

// extractCodeBlock extracts first code block from text
func extractCodeBlock(text string) string {
	// Match ```language\ncode\n``` pattern
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractHints parses hints from sections
func extractHints(sections map[string]string) []string {
	var hints []string

	hintsText, ok := sections["Hints"]
	if !ok {
		return hints
	}

	// Try to extract from ### Hint N format
	re := regexp.MustCompile(`(?m)^### Hint \d+:?\s*(.+?)$`)
	matches := re.FindAllStringSubmatchIndex(hintsText, -1)

	if len(matches) > 0 {
		// Structured hints
		for i, match := range matches {
			start := match[1]
			var end int
			if i+1 < len(matches) {
				end = matches[i+1][0]
			} else {
				end = len(hintsText)
			}
			hintText := strings.TrimSpace(hintsText[start:end])
			hints = append(hints, hintText)
		}
	} else {
		// Fallback: look for numbered list or details tags
		// Extract from <details><summary> tags
		detailsRe := regexp.MustCompile(`(?s)<details>.*?<summary>(.*?)</summary>(.*?)</details>`)
		detailsMatches := detailsRe.FindAllStringSubmatch(hintsText, -1)

		for _, dm := range detailsMatches {
			if len(dm) > 2 {
				title := strings.TrimSpace(dm[1])
				content := strings.TrimSpace(dm[2])
				hints = append(hints, title+"\n\n"+content)
			}
		}
	}

	return hints
}

// extractOutcomes parses expected outcomes
func extractOutcomes(sections map[string]string) []string {
	var outcomes []string

	text, ok := sections["Expected Outcomes"]
	if !ok {
		text, ok = sections["Expected Student Outcomes"]
		if !ok {
			return outcomes
		}
	}

	// Extract bullet points (- or ✅)
	re := regexp.MustCompile(`(?m)^[-✅]\s+(.+?)$`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) > 1 {
			outcomes = append(outcomes, match[1])
		}
	}

	return outcomes
}

// extractMistakes parses common mistakes
func extractMistakes(sections map[string]string) []string {
	var mistakes []string

	text, ok := sections["Common Mistakes"]
	if !ok {
		text, ok = sections["Common Mistakes to Watch For"]
		if !ok {
			return mistakes
		}
	}

	// Extract numbered list items
	re := regexp.MustCompile(`(?m)^\d+\.\s+\*\*(.+?)\*\*:?\s*(.+?)$`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) > 2 {
			pattern := match[1]
			guidance := match[2]
			mistakes = append(mistakes, pattern+": "+guidance)
		}
	}

	return mistakes
}

// extractAIPrompt gets the AI evaluation prompt
func extractAIPrompt(sections map[string]string) string {
	if prompt, ok := sections["AI Evaluation Prompt"]; ok {
		// Extract from code block if present
		code := extractCodeBlock(prompt)
		if code != "" {
			return code
		}
		return prompt
	}
	if prompt, ok := sections["Sample AI Feedback Prompt"]; ok {
		code := extractCodeBlock(prompt)
		if code != "" {
			return code
		}
		return prompt
	}
	return ""
}

// ToExercise converts ParsedExercise to database Exercise model
func (p *ParsedExercise) ToExercise() (*Exercise, error) {
	id, err := uuid.Parse(p.Frontmatter.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid exercise ID: %w", err)
	}

	nodeID, err := uuid.Parse(p.Frontmatter.NodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	// Build ExerciseContent
	content := ExerciseContent{
		StarterCode:       p.StarterCode,
		ReferenceSolution: p.Solution,
		TestCases:         parseTestCases(p.TestCases),
		Hints:             buildHints(p.Hints),
		CommonMistakes:    buildMistakes(p.Mistakes),
		ExpectedOutcomes:  p.Outcomes,
		AIEvaluationPrompt: p.AIPrompt,
		GradingRubric: GradingRubric{
			TestWeight: 60,
			AIWeight:   40,
		},
	}

	exercise := &Exercise{
		ID:               id,
		NodeID:           nodeID,
		Title:            p.Frontmatter.Title,
		Slug:             p.Frontmatter.Slug,
		Description:      p.Description,
		Difficulty:       p.Frontmatter.Difficulty,
		EstimatedMinutes: &p.Frontmatter.EstimatedMinutes,
		OrderIndex:       p.Frontmatter.OrderIndex,
		Content:          content,
	}

	return exercise, nil
}

// parseTestCases extracts individual test functions
func parseTestCases(testCode string) []TestCase {
	var testCases []TestCase

	// Extract #[test] functions
	re := regexp.MustCompile(`(?s)#\[test\]\s+fn\s+(\w+)\s*\(\)\s*\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(testCode, -1)

	for _, match := range matches {
		if len(match) > 2 {
			name := match[1]
			body := strings.TrimSpace(match[2])

			testCases = append(testCases, TestCase{
				Name:    name,
				Code:    fmt.Sprintf("#[test]\nfn %s() {\n%s\n}", name, body),
				Visible: true, // All visible by default
			})
		}
	}

	return testCases
}

// buildHints converts hint strings to Hint structs
func buildHints(hints []string) []Hint {
	var result []Hint

	for i, hintText := range hints {
		// Try to extract title from first line
		lines := strings.SplitN(hintText, "\n", 2)
		title := lines[0]
		content := hintText
		if len(lines) > 1 {
			content = lines[1]
		}

		result = append(result, Hint{
			Order:   i + 1,
			Title:   title,
			Content: content,
		})
	}

	return result
}

// buildMistakes converts mistake strings to CommonMistake structs
func buildMistakes(mistakes []string) []CommonMistake {
	var result []CommonMistake

	for _, mistake := range mistakes {
		// Split on first colon
		parts := strings.SplitN(mistake, ":", 2)
		pattern := mistake
		guidance := ""

		if len(parts) == 2 {
			pattern = strings.TrimSpace(parts[0])
			guidance = strings.TrimSpace(parts[1])
		}

		result = append(result, CommonMistake{
			Pattern:  pattern,
			Guidance: guidance,
		})
	}

	return result
}
