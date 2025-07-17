package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GiteaClient represents a client for Gitea API
type GiteaClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewGiteaClient creates a new Gitea API client
func NewGiteaClient(baseURL, token string) *GiteaClient {
	return &GiteaClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		client:  &http.Client{},
	}
}

// CommitStatus represents a commit status
type CommitStatus struct {
	State       string `json:"state"`       // pending, success, error, failure
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

// CommitInfo represents commit information from Gitea API
type CommitInfo struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

// SetCommitStatus sets the status of a commit
func (g *GiteaClient) SetCommitStatus(owner, repo, sha string, status CommitStatus) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/statuses/%s", g.baseURL, owner, repo, sha)
	
	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "token "+g.token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// GetPullRequestCommits gets all commits in a pull request
func (g *GiteaClient) GetPullRequestCommits(owner, repo string, number int) ([]CommitInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/pulls/%d/commits", g.baseURL, owner, repo, number)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "token "+g.token)
	
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}
	
	var commits []CommitInfo
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}
	
	return commits, nil
}