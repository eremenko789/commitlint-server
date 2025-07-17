package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// processPullRequest processes a pull request webhook
func (s *Server) processPullRequest(webhook *PullRequestWebhook) {
	client := NewGiteaClient(s.config.Server.GiteaURL, s.config.Server.GiteaToken)
	
	owner := webhook.Repository.Owner.Login
	repo := webhook.Repository.Name
	prNumber := webhook.PullRequest.Number
	
	if s.config.Server.Debug {
		log.Printf("Начинаем обработку PR #%d в %s/%s", prNumber, owner, repo)
	}

	// Get all commits in the pull request
	commits, err := client.GetPullRequestCommits(owner, repo, prNumber)
	if err != nil {
		log.Printf("Ошибка получения коммитов для PR #%d: %v", prNumber, err)
		return
	}

	if len(commits) == 0 {
		log.Printf("Нет коммитов в PR #%d", prNumber)
		return
	}

	// Process each commit
	var allPassed = true
	var failedCommits []string
	var totalErrors = 0

	for _, commit := range commits {
		if s.config.Server.Debug {
			log.Printf("Проверяем коммит %s: %s", commit.SHA[:8], commit.Commit.Message)
		}

		// Set status to pending
		status := CommitStatus{
			State:       "pending",
			Context:     "commitlint",
			Description: "Проверка сообщения коммита...",
		}
		
		if err := client.SetCommitStatus(owner, repo, commit.SHA, status); err != nil {
			log.Printf("Ошибка установки статуса pending для коммита %s: %v", commit.SHA, err)
		}

		// Lint the commit message
		result, err := s.linter.ParseAndLint(commit.Commit.Message)
		if err != nil {
			log.Printf("Ошибка линтинга коммита %s: %v", commit.SHA, err)
			continue
		}
		
		issues := result.Issues()
		if len(issues) == 0 {
			// Commit is valid
			status = CommitStatus{
				State:       "success",
				Context:     "commitlint",
				Description: "Сообщение коммита соответствует правилам",
			}
			if s.config.Server.Debug {
				log.Printf("✓ Коммит %s прошел проверку", commit.SHA[:8])
			}
		} else {
			// Check if there are any errors (not just warnings)
			hasErrors := false
			errorCount := 0
			for _, issue := range issues {
				if issue.Severity() == lint.SeverityError {
					hasErrors = true
					errorCount++
				}
			}
			
			if hasErrors {
				// Commit has errors
				allPassed = false
				failedCommits = append(failedCommits, commit.SHA[:8])
				totalErrors += errorCount
				
				// Build error description
				var errorMessages []string
				for _, issue := range issues {
					if issue.Severity() == lint.SeverityError {
						errorMessages = append(errorMessages, issue.Description())
					}
				}
				
				description := fmt.Sprintf("Найдено ошибок: %d", errorCount)
				if len(errorMessages) > 0 {
					description += " - " + strings.Join(errorMessages, "; ")
				}
				
				// Limit description length for API
				if len(description) > 140 {
					description = description[:137] + "..."
				}
				
				status = CommitStatus{
					State:       "failure",
					Context:     "commitlint",
					Description: description,
				}
				
				log.Printf("✗ Коммит %s не прошел проверку: %s", commit.SHA[:8], description)
			} else {
				// Only warnings, treat as success
				status = CommitStatus{
					State:       "success",
					Context:     "commitlint",
					Description: fmt.Sprintf("Сообщение коммита соответствует правилам (предупреждений: %d)", len(issues)),
				}
				if s.config.Server.Debug {
					log.Printf("✓ Коммит %s прошел проверку с предупреждениями", commit.SHA[:8])
				}
			}
		}

		// Set final status
		if err := client.SetCommitStatus(owner, repo, commit.SHA, status); err != nil {
			log.Printf("Ошибка установки финального статуса для коммита %s: %v", commit.SHA, err)
		}
	}

	// Log summary
	if allPassed {
		log.Printf("✓ Все %d коммитов в PR #%d прошли проверку", len(commits), prNumber)
	} else {
		log.Printf("✗ %d из %d коммитов в PR #%d не прошли проверку (всего ошибок: %d)", 
			len(failedCommits), len(commits), prNumber, totalErrors)
	}
}