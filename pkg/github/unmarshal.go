package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v42/github"
)

func (c *Client) unmarshalPullRequestEvent(ctx context.Context, event *github.Event) error {
	e := new(github.PullRequestEvent)
	err := json.Unmarshal(event.GetRawPayload(), e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}

	pr := e.GetPullRequest()
	repo := event.GetRepo()
	// override pr since current information about the pr could be newer
	currentPR, err := c.getPullRequest(ctx, repo.GetName(), pr.GetNumber())
	if err != nil {
		return err
	}
	*pr = *currentPR

	prWithRepo := PullRequestWithRepository{
		PullRequest: pr,
		Repo:        repo,
	}
	c.collectPullRequest(pr.GetID(), prWithRepo)
	return nil
}

func (c *Client) unmarshalIssuesEvent(ctx context.Context, event *github.Event) error {
	e := new(github.IssuesEvent)
	err := json.Unmarshal(*event.RawPayload, e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}
	issue := e.GetIssue()
	issueWithRepo, err := c.createIssueWithRepo(ctx, event, issue)
	if err != nil {
		return fmt.Errorf("creating issue with repo struct: %w", err)
	}
	c.collectIssue(issue.GetID(), issueWithRepo)
	return nil
}

func (c *Client) unmarshalIssuesCommentEvent(ctx context.Context, event *github.Event) error {

	e := new(github.IssueCommentEvent)
	err := json.Unmarshal(*event.RawPayload, e)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}
	issue := e.GetIssue()
	issueWithRepo, err := c.createIssueWithRepo(ctx, event, issue)
	issueWithRepo.Comment = e.Comment
	c.collectIssue(e.GetComment().GetID(), issueWithRepo)
	return nil
}

func (c *Client) createIssueWithRepo(ctx context.Context, event *github.Event, issue *github.Issue) (IssueWithRepository, error) {
	repo := event.GetRepo()
	// override pr since current information about the issue could be newer
	currentIssue, err := c.getIssue(ctx, repo.GetName(), issue.GetNumber())
	if err != nil {
		return IssueWithRepository{}, err
	}
	*issue = *currentIssue
	return IssueWithRepository{
		Issue: issue,
		Repo:  repo,
	}, nil
}

func (c *Client) getPullRequest(ctx context.Context, repoName string, number int) (*github.PullRequest, error) {
	owner, name := splitRepoName(repoName)
	pr, _, err := c.client.PullRequests.Get(ctx, owner, name, number)
	if err != nil {
		return nil, fmt.Errorf("getting pullrequest %v/%v-#%v: %w", owner, name, number, err)
	}
	return pr, nil
}

func (c *Client) getIssue(ctx context.Context, repoName string, number int) (*github.Issue, error) {
	owner, name := splitRepoName(repoName)
	issue, _, err := c.client.Issues.Get(ctx, owner, name, number)
	if err != nil {
		return nil, fmt.Errorf("getting issue %v/%v-#%v: %w", owner, name, number, err)
	}
	return issue, nil
}
