package github

import (
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Run("issues", testIssueWithRepository)
	t.Run("pullrequests", testPullRequestWithRepository)
}

type testCase[T Item] struct {
	name     string
	start    store[T]
	key      int64
	val      T
	expected store[T]
}

func testIssueWithRepository(t *testing.T) {
	t.Parallel()
	tests := []testCase[IssueWithRepository]{
		{
			name:  "empty map",
			start: store[IssueWithRepository]{},
			key:   123,
			val:   IssueWithRepository{},
			expected: store[IssueWithRepository]{
				123: IssueWithRepository{},
			},
		},
		{
			name: "existing key",
			start: store[IssueWithRepository]{
				123: IssueWithRepository{Issue: &github.Issue{}},
			},
			key: 123,
			val: IssueWithRepository{},
			expected: store[IssueWithRepository]{
				123: IssueWithRepository{Issue: &github.Issue{}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, runTestCase(tc))
	}
}

func testPullRequestWithRepository(t *testing.T) {
	t.Parallel()
	tests := []testCase[PullRequestWithRepository]{
		{
			name:  "empty map",
			start: store[PullRequestWithRepository]{},
			key:   123,
			val:   PullRequestWithRepository{},
			expected: store[PullRequestWithRepository]{
				123: PullRequestWithRepository{},
			},
		},
		{
			name: "existing key",
			start: store[PullRequestWithRepository]{
				123: PullRequestWithRepository{PullRequest: &github.PullRequest{}},
			},
			key: 123,
			val: PullRequestWithRepository{},
			expected: store[PullRequestWithRepository]{
				123: PullRequestWithRepository{PullRequest: &github.PullRequest{}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, runTestCase(tc))
	}
}

func runTestCase[T Item](tc testCase[T]) func(t *testing.T) {
	return func(t *testing.T) {
		tc.start.add(tc.key, tc.val)
		assert.Equal(t, tc.start, tc.expected)
	}
}
