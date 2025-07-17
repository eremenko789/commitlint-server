package webhook

import "time"

// PullRequestPayload represents a pull request webhook payload
type PullRequestPayload struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

// PullRequest represents a pull request
type PullRequest struct {
	ID          int64          `json:"id"`
	Number      int            `json:"number"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	State       string         `json:"state"`
	Base        PullRequestRef `json:"base"`
	Head        PullRequestRef `json:"head"`
	MergeBase   string         `json:"merge_base"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PullRequestRef represents a pull request reference
type PullRequestRef struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	Sha   string     `json:"sha"`
	Repo  Repository `json:"repo"`
}

// Repository represents a repository
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
	Private  bool   `json:"private"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

// User represents a user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
}