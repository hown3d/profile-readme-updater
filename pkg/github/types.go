package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
)

type Infos struct {
	PullRequests map[int64]PullRequestWithRepository
	Issues       map[int64]IssueWithRepository
	Languages    Languages
}

type PullRequestWithRepository struct {
	PullRequest *github.PullRequest
	Repo        *github.Repository
}

func (c *Client) newPullRequestWithRepo(ctx context.Context, repo *github.Repository, pr *github.PullRequest) (PullRequestWithRepository, error) {
	err := c.getPullRequestInfo(ctx, getRepoName(repo), pr)
	if err != nil {
		return PullRequestWithRepository{}, fmt.Errorf("getting pr info: %w", err)
	}
	return PullRequestWithRepository{
		PullRequest: pr,
		Repo:        repo,
	}, nil
}

type IssueWithRepository struct {
	Issue   *github.Issue
	Comment *github.IssueComment
	Repo    *github.Repository
}

func (c *Client) newIssueWithRepo(ctx context.Context, repo *github.Repository, issue *github.Issue) (IssueWithRepository, error) {
	// override issue since current information about the issue could be newer
	err := c.getIssueInfo(ctx, getRepoName(repo), issue)
	if err != nil {
		return IssueWithRepository{}, fmt.Errorf("getting issue info: %w", err)
	}
	return IssueWithRepository{
		Issue: issue,
		Repo:  repo,
	}, nil
}

type Languages map[string]int

func (l Languages) IncreaseCount(language string) {
	l[language]++
}

func (l Languages) totalRepos() (total int) {
	for _, count := range l {
		total += count
	}
	return total
}

func (l Languages) Percentage(language string) float64 {
	return float64(l[language]) / float64(l.totalRepos()) * 100
}
