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
	client *github.Client
	user   string
	infos  *Infos
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
		infos: &Infos{
			PullRequests: map[int64]PullRequestWithRepository{},
			Issues:       map[int64]IssueWithRepository{},
			Languages:    Languages{},
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
			repo := event.GetRepo()
			err := c.getRepoInfo(ctx, repo)
			if err != nil {
				return fmt.Errorf("getting repo info: %w", err)
			}
			switch *event.Type {
			case "IssueCommentEvent":
				err = c.collectIssueCommentEvent(ctx, event, repo)
			case "IssuesEvent":
				err = c.collectIssuesEvent(ctx, event, repo)
			case "PullRequestEvent":
				err = c.collectPullRequestEvent(ctx, event, repo)
			// non of the wanted events, go again
			default:
				continue
			}
			if err != nil {
				return fmt.Errorf("unmarshal event: %w", err)
			}
			c.infos.Languages.IncreaseCount(repo.GetLanguage())
		}

		if resp.NextPage == 0 {
			break pagingLoop
		}
		opts.Page = resp.NextPage
	}
	return nil
}

func (c *Client) GetInfos() *Infos {
	return c.infos
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
