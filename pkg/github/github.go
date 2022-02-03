package github

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client          *github.Client
	user            string
	collectedEvents *Events
}

func NewGithubClient() (*github.Client, error) {
	accessToken := os.Getenv("GITHUB_TOKEN")
	if accessToken == "" {
		return nil, errors.New("Github Token env variable is empty")
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return github.NewClient(httpClient), nil
}

func NewClient(githubClient *github.Client) (*Client, error) {
	username, err := getUsername(context.Background(), githubClient)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: githubClient,
		user:   username,
		collectedEvents: &Events{
			PullRequests: map[int64]PullRequestWithRepository{},
			Issues:       map[int64]IssueWithRepository{},
		},
	}, nil
}

func (c *Client) GetContributions(ctx context.Context, earliest time.Time) error {
	opts := &github.ListOptions{
		PerPage: 10,
	}

pagingLoop:
	for {
		events, resp, err := c.client.Activity.ListEventsPerformedByUser(ctx, c.user, false, opts)
		if err != nil {
			return fmt.Errorf("listing events of user: %w", err)
		}
		for _, event := range events {
			if event.CreatedAt.Before(earliest) {
				break pagingLoop
			}
			var err error
			switch *event.Type {
			case "IssueCommentEvent":
				err = c.unmarshalIssueCommentEvent(ctx, event)
			case "IssuesEvent":
				err = c.unmarshalIssuesEvent(ctx, event)
			case "PullRequestEvent":
				err = c.unmarshalPullRequestEvent(ctx, event)
			}
			if err != nil {
				return fmt.Errorf("unmarshal event: %w", err)
			}
		}

		if resp.NextPage == 0 {
			break pagingLoop
		}

		opts.Page = resp.NextPage
	}
	return nil
}

func (c *Client) CollectedEvents() *Events {
	return c.collectedEvents
}

func getUsername(ctx context.Context, client *github.Client) (string, error) {
	// if user is empty, use the authenticated user
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("getting user: %w", err)
	}
	return user.GetLogin(), nil
}

func splitRepoName(repoName string) (owner, name string) {
	owner = strings.Split(repoName, "/")[0]
	name = strings.Split(repoName, "/")[1]
	return
}
