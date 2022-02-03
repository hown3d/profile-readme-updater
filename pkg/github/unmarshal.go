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

	repo := event.GetRepo()
	err = c.getRepoInfo(ctx, repo)
	if err != nil {
		return err
	}
	// override pr since current information about the pr could be newer
	pr := e.GetPullRequest()
	err = c.getPullRequestInfo(ctx, getRepoName(repo), pr)
	if err != nil {
		return err
	}

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

func (c *Client) unmarshalIssueCommentEvent(ctx context.Context, event *github.Event) error {
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
	err := c.getRepoInfo(ctx, repo)
	if err != nil {
		return IssueWithRepository{}, err
	}
	// override issue since current information about the issue could be newer
	err = c.getIssueInfo(ctx, getRepoName(repo), issue)
	if err != nil {
		return IssueWithRepository{}, err
	}
	return IssueWithRepository{
		Issue: issue,
		Repo:  repo,
	}, nil
}

func (c *Client) getPullRequestInfo(ctx context.Context, repoName string, pr *github.PullRequest) error {
	owner, name := splitRepoName(repoName)
	prWithInfo, _, err := c.client.PullRequests.Get(ctx, owner, name, pr.GetNumber())
	if err != nil {
		return fmt.Errorf("getting pullrequest %v/%v-#%v: %w", owner, name, pr.GetNumber(), err)
	}
	*pr = *prWithInfo
	return nil
}

func (c *Client) getIssueInfo(ctx context.Context, repoName string, issue *github.Issue) error {
	owner, name := splitRepoName(repoName)
	issueWithInfo, _, err := c.client.Issues.Get(ctx, owner, name, issue.GetNumber())
	if err != nil {
		return fmt.Errorf("getting issue %v/%v-#%v: %w", owner, name, issue.GetNumber(), err)
	}
	*issue = *issueWithInfo
	return nil
}

func (c *Client) getRepoInfo(ctx context.Context, repo *github.Repository) error {
	owner, name := splitRepoName(getRepoName(repo))
	repoWithInfo, _, err := c.client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("getting repo %v/%v: %w", owner, name, err)
	}
	*repo = *repoWithInfo
	return nil
}

func getRepoName(repo *github.Repository) string {
	if repo.FullName != nil {
		return repo.GetFullName()
	}
	return repo.GetName()
}
