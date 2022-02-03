package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v42/github"
)

// collectIssue adds an issue to the map if the id of the issue isn't present already
func (c *Client) collectIssue(key int64, val IssueWithRepository) {
	_, exists := c.infos.Issues[key]
	if exists {
		return
	}
	c.infos.Issues[key] = val
}

// collectIssue adds a pr to the map if the id of the pr isn't present already
func (c *Client) collectPullRequest(key int64, val PullRequestWithRepository) {
	_, exists := c.infos.PullRequests[key]
	if exists {
		return
	}
	c.infos.PullRequests[key] = val
}

func (c *Client) collectPullRequestEvent(ctx context.Context, event *github.Event, repo *github.Repository) error {
	e := new(github.PullRequestEvent)
	err := json.Unmarshal(event.GetRawPayload(), e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}

	// override pr since current information about the pr could be newer
	pr := e.GetPullRequest()
	prWithRepo, err := c.newPullRequestWithRepo(ctx, repo, pr)
	if err != nil {
		return fmt.Errorf("creating pr with repo struct: %w", err)
	}
	c.collectPullRequest(pr.GetID(), prWithRepo)
	return nil
}

func (c *Client) collectIssuesEvent(ctx context.Context, event *github.Event, repo *github.Repository) error {
	e := new(github.IssuesEvent)
	err := json.Unmarshal(*event.RawPayload, e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}
	issue := e.GetIssue()
	issueWithRepo, err := c.newIssueWithRepo(ctx, repo, issue)
	if err != nil {
		return fmt.Errorf("creating issue with repo struct: %w", err)
	}
	c.collectIssue(issue.GetID(), issueWithRepo)
	return nil
}

func (c *Client) collectIssueCommentEvent(ctx context.Context, event *github.Event, repo *github.Repository) error {
	e := new(github.IssueCommentEvent)
	err := json.Unmarshal(*event.RawPayload, e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}
	issue := e.GetIssue()
	issueWithRepo, err := c.newIssueWithRepo(ctx, repo, issue)
	issueWithRepo.Comment = e.Comment
	c.collectIssue(e.GetComment().GetID(), issueWithRepo)
	return nil
}
