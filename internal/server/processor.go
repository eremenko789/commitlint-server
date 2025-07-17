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
		result := s.linter.Lint(commit.Commit.Message, commit.SHA)
		
		if result.Valid {
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
			// Commit is invalid
			allPassed = false
			failedCommits = append(failedCommits, commit.SHA[:8])
			totalErrors += len(result.Errors)
			
			// Build error description
			var errorMessages []string
			for _, err := range result.Errors {
				errorMessages = append(errorMessages, err.Message)
			}
			
			description := fmt.Sprintf("Найдено ошибок: %d", len(result.Errors))
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