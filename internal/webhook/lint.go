package webhook

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// LintExecutor executes commit linting
type LintExecutor struct {
	config CommitlintConfig
}

// NewLintExecutor creates a new lint executor
func NewLintExecutor(config CommitlintConfig) *LintExecutor {
	return &LintExecutor{
		config: config,
	}
}

// LintResult represents the result of linting
type LintResult struct {
	Valid  bool
	Issues []LintIssue
}

// LintIssue represents a single lint issue
type LintIssue struct {
	Level   string
	Message string
}

// LintCommit lints a single commit message
func (e *LintExecutor) LintCommit(message string) (*LintResult, error) {
	// Load lint configuration
	lintConfig, err := config.Parse(e.config.ConfigPath)
	if err != nil {
		// If config file doesn't exist, use defaults
		lintConfig = config.NewDefault()
	}

	// Create linter
	linter, err := config.NewLinter(lintConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create linter: %w", err)
	}

	// Lint the message
	report, err := linter.ParseAndLint(message)
	if err != nil {
		return nil, fmt.Errorf("failed to lint message: %w", err)
	}

	// Convert report to result
	issues := report.Issues()
	result := &LintResult{
		Valid:  len(issues) == 0,
		Issues: make([]LintIssue, 0, len(issues)),
	}

	// Collect issues
	for _, issue := range issues {
		level := "error"
		if issue.Severity() == lint.SeverityWarn {
			level = "warning"
		}

		result.Issues = append(result.Issues, LintIssue{
			Level:   level,
			Message: fmt.Sprintf("%s: %s", issue.RuleName(), issue.Description()),
		})
	}

	return result, nil
}

// LintCommits lints multiple commit messages
func (e *LintExecutor) LintCommits(messages []string) ([]*LintResult, error) {
	results := make([]*LintResult, len(messages))

	for i, message := range messages {
		result, err := e.LintCommit(message)
		if err != nil {
			return nil, fmt.Errorf("failed to lint commit %d: %w", i, err)
		}
		results[i] = result
	}

	return results, nil
}