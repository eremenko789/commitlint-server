package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GiteaClient represents a client for Gitea API
type GiteaClient struct {
	config     GiteaConfig
	httpClient *http.Client
}

// NewGiteaClient creates a new Gitea client
func NewGiteaClient(config GiteaConfig) *GiteaClient {
	return &GiteaClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Commit represents a git commit
type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Author  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

// CommitStatus represents a commit status
type CommitStatus struct {
	State       string `json:"state"`       // pending, success, error, failure
	TargetURL   string `json:"target_url"`  
	Description string `json:"description"`
	Context     string `json:"context"`
}

// GetPullRequestCommits gets commits for a pull request
func (c *GiteaClient) GetPullRequestCommits(owner, repo string, prNumber int) ([]Commit, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/pulls/%d/commits",
		c.config.BaseURL, owner, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.config.Token))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits: %s", resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// CreateCommitStatus creates a commit status
func (c *GiteaClient) CreateCommitStatus(owner, repo, sha string, status *CommitStatus) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/statuses/%s",
		c.config.BaseURL, owner, repo, sha)

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.config.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create status: %s", resp.Status)
	}

	return nil
}

// CreatePullRequestComment creates a comment on a pull request
func (c *GiteaClient) CreatePullRequestComment(owner, repo string, prNumber int, comment string) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/issues/%d/comments",
		c.config.BaseURL, owner, repo, prNumber)

	body := map[string]string{
		"body": comment,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.config.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create comment: %s", resp.Status)
	}

	return nil
}