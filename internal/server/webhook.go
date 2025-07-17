package server

import "time"

// PullRequestWebhook represents a Gitea pull request webhook payload
type PullRequestWebhook struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

// PullRequest represents a pull request in Gitea webhook
type PullRequest struct {
	ID          int       `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	Head        Branch    `json:"head"`
	Base        Branch    `json:"base"`
	User        User      `json:"user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MergeCommitSHA string `json:"merge_commit_sha"`
	HTMLURL     string    `json:"html_url"`
}

// Branch represents a git branch
type Branch struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	SHA   string     `json:"sha"`
	Repo  Repository `json:"repo"`
}

// Repository represents a git repository
type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
	Private  bool   `json:"private"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

// User represents a Gitea user
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// Commit represents a git commit
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	URL       string    `json:"url"`
	Author    CommitUser `json:"author"`
	Committer CommitUser `json:"committer"`
	Timestamp time.Time `json:"timestamp"`
}

// CommitUser represents the author/committer of a commit
type CommitUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}