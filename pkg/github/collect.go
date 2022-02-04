package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v42/github"
)

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
	c.Infos.PullRequests.add(pr.GetID(), prWithRepo)
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

	c.Infos.Issues.add(issue.GetID(), issueWithRepo)
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
	c.Infos.Issues.add(e.GetComment().GetID(), issueWithRepo)
	return nil
}
