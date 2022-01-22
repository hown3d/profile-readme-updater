package github

import "context"

type mockGithubClient struct{}

var _ githubClient = (*mockGithubClient)(nil)

func (c mockGithubClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	panic("not implemented")
}
