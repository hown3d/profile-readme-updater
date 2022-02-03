package github

import "github.com/google/go-github/v42/github"

type Events struct {
	PullRequests map[int64]PullRequestWithRepository
	Issues       map[int64]IssueWithRepository
}

type PullRequestWithRepository struct {
	PullRequest *github.PullRequest
	Repo        *github.Repository
}

type IssueWithRepository struct {
	Issue   *github.Issue
	Comment *github.IssueComment
	Repo    *github.Repository
}
