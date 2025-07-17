package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Server represents the webhook server
type Server struct {
	config       *Config
	giteaClient  *GiteaClient
	lintExecutor *LintExecutor
}

// NewServer creates a new webhook server
func NewServer(config *Config) *Server {
	return &Server{
		config:       config,
		giteaClient:  NewGiteaClient(config.Gitea),
		lintExecutor: NewLintExecutor(config.Commitlint),
	}
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Route requests
	switch r.URL.Path {
	case s.config.Webhook.Path:
		s.handleWebhook(w, r)
	case "/health":
		s.handleHealth(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handleWebhook handles incoming webhook requests
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	// Verify signature if secret is configured
	if s.config.Webhook.Secret != "" {
		signature := r.Header.Get("X-Gitea-Signature")
		if !s.verifySignature(body, signature) {
			log.Printf("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Parse event type
	eventType := r.Header.Get("X-Gitea-Event")
	if eventType == "" {
		log.Printf("Missing X-Gitea-Event header")
		http.Error(w, "Missing event header", http.StatusBadRequest)
		return
	}

	// Check if event is configured
	eventConfigured := false
	for _, e := range s.config.Webhook.Events {
		if e == eventType {
			eventConfigured = true
			break
		}
	}
	if !eventConfigured {
		log.Printf("Event type %s not configured", eventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle pull request events
	if eventType == "pull_request" || eventType == "pull_request_sync" {
		var payload PullRequestPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Failed to parse pull request payload: %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		// Process only opened, synchronize, and reopened actions
		if payload.Action == "opened" || payload.Action == "synchronize" || payload.Action == "reopened" {
			go s.processPullRequest(&payload)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// verifySignature verifies the webhook signature
func (s *Server) verifySignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(s.config.Webhook.Secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// processPullRequest processes a pull request event
func (s *Server) processPullRequest(payload *PullRequestPayload) {
	log.Printf("Processing pull request #%d: %s", payload.PullRequest.Number, payload.PullRequest.Title)

	// Set pending status
	status := &CommitStatus{
		State:       "pending",
		TargetURL:   "",
		Description: "Running commit lint...",
		Context:     "commitlint",
	}

	if err := s.giteaClient.CreateCommitStatus(
		payload.Repository.Owner.Username,
		payload.Repository.Name,
		payload.PullRequest.Head.Sha,
		status,
	); err != nil {
		log.Printf("Failed to set pending status: %v", err)
	}

	// Get commits
	commits, err := s.giteaClient.GetPullRequestCommits(
		payload.Repository.Owner.Username,
		payload.Repository.Name,
		payload.PullRequest.Number,
	)
	if err != nil {
		log.Printf("Failed to get commits: %v", err)
		s.setErrorStatus(payload, "Failed to fetch commits")
		return
	}

	// Lint commits
	var errors []string
	for _, commit := range commits {
		result, err := s.lintExecutor.LintCommit(commit.Message)
		if err != nil {
			log.Printf("Failed to lint commit %s: %v", commit.SHA[:7], err)
			errors = append(errors, fmt.Sprintf("Failed to lint commit %s", commit.SHA[:7]))
			continue
		}

		if !result.Valid {
			for _, issue := range result.Issues {
				errors = append(errors, fmt.Sprintf("[%s] %s: %s", commit.SHA[:7], issue.Level, issue.Message))
			}
		}
	}

	// Set final status
	if len(errors) > 0 {
		description := fmt.Sprintf("Found %d issue(s)", len(errors))
		if len(errors) == 1 {
			description = errors[0]
		}
		status.State = "failure"
		status.Description = description
	} else {
		status.State = "success"
		status.Description = "All commits pass lint checks"
	}

	if err := s.giteaClient.CreateCommitStatus(
		payload.Repository.Owner.Username,
		payload.Repository.Name,
		payload.PullRequest.Head.Sha,
		status,
	); err != nil {
		log.Printf("Failed to set final status: %v", err)
	}

	// Add comment with details if there are errors
	if len(errors) > 0 {
		comment := "## Commit Lint Results\n\n"
		comment += "The following issues were found:\n\n"
		for _, err := range errors {
			comment += fmt.Sprintf("- %s\n", err)
		}
		comment += "\nPlease fix these issues and push new commits."

		if err := s.giteaClient.CreatePullRequestComment(
			payload.Repository.Owner.Username,
			payload.Repository.Name,
			payload.PullRequest.Number,
			comment,
		); err != nil {
			log.Printf("Failed to create comment: %v", err)
		}
	}
}

// setErrorStatus sets an error status on the commit
func (s *Server) setErrorStatus(payload *PullRequestPayload, message string) {
	status := &CommitStatus{
		State:       "error",
		TargetURL:   "",
		Description: message,
		Context:     "commitlint",
	}

	if err := s.giteaClient.CreateCommitStatus(
		payload.Repository.Owner.Username,
		payload.Repository.Name,
		payload.PullRequest.Head.Sha,
		status,
	); err != nil {
		log.Printf("Failed to set error status: %v", err)
	}
}