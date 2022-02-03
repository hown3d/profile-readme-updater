package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
)

func (c *Client) getPullRequestInfo(ctx context.Context, repoName string, pr *github.PullRequest) error {
	owner, name, err := splitRepoName(repoName)
	if err != nil {
		return err
	}
	prWithInfo, _, err := c.client.PullRequests.Get(ctx, owner, name, pr.GetNumber())
	if err != nil {
		return fmt.Errorf("getting pullrequest %v/%v-#%v: %w", owner, name, pr.GetNumber(), err)
	}
	*pr = *prWithInfo
	return nil
}

func (c *Client) getIssueInfo(ctx context.Context, repoName string, issue *github.Issue) error {
	owner, name, err := splitRepoName(repoName)
	if err != nil {
		return err
	}
	issueWithInfo, _, err := c.client.Issues.Get(ctx, owner, name, issue.GetNumber())
	if err != nil {
		return fmt.Errorf("getting issue %v/%v-#%v: %w", owner, name, issue.GetNumber(), err)
	}
	*issue = *issueWithInfo
	return nil
}

func (c *Client) getRepoInfo(ctx context.Context, repo *github.Repository) error {
	repoName := getRepoName(repo)
	owner, name, err := splitRepoName(repoName)
	if err != nil {
		return err
	}
	repoWithInfo, _, err := c.client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("getting repo %v: %w", repo.GetID(), err)
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
