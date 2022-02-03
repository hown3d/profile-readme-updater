package github

import (
	"context"
	"encoding/json"
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
}

func NewGithubClient() *github.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return github.NewClient(httpClient)
}

func NewClient(githubClient *github.Client) (Client, error) {
	username, err := getUsername(context.Background(), githubClient)
	if err != nil {
		return Client{}, nil
	}
	return Client{
		client: githubClient,
		user:   username,
	}, nil
}

type Events struct {
	PullRequests map[int64]PullRequestWithRepository
	Issues       map[int64]IssueWithRepository
}

type PullRequestWithRepository struct {
	PullRequest *github.PullRequest
	Repo        *github.Repository
}

type IssueWithRepository struct {
	Issue *github.Issue
	Repo  *github.Repository
}

func (c Client) GetContributions(ctx context.Context, earliest time.Time) (Events, error) {
	opts := &github.ListOptions{
		PerPage: 10,
	}
	collected := Events{
		PullRequests: map[int64]PullRequestWithRepository{},
		Issues:       map[int64]IssueWithRepository{},
	}

pagingLoop:
	for {
		events, resp, err := c.client.Activity.ListEventsPerformedByUser(ctx, c.user, false, opts)
		if err != nil {
			return collected, err
		}
		for _, event := range events {
			if event.CreatedAt.Before(earliest) {
				break pagingLoop
			}
			switch *event.Type {
			case "IssuesEvent":
				e := new(github.IssuesEvent)
				err := json.Unmarshal(*event.RawPayload, e)
				if err != nil {
					return collected, fmt.Errorf("unmarshaling payload: %w", err)
				}

				issue := e.GetIssue()
				_, exists := collected.Issues[issue.GetID()]
				// key is already populated, so this issue was already present in a previous event
				if exists {
					break
				}

				repo := event.GetRepo()
				// override state, since it could be newer
				currentIssueState, err := c.getIssueState(ctx, repo.GetName(), issue.GetNumber())
				if err != nil {
					return collected, err
				}
				*issue.State = currentIssueState

				issueWithRepo := IssueWithRepository{
					Issue: issue,
					Repo:  repo,
				}
				collected.Issues[issue.GetID()] = issueWithRepo
			case "PullRequestEvent":
				e := new(github.PullRequestEvent)
				err := json.Unmarshal(*event.RawPayload, e)
				if err != nil {
					return collected, fmt.Errorf("unmarshaling payload: %w", err)
				}

				pr := e.GetPullRequest()
				_, exists := collected.Issues[pr.GetID()]
				// key is already populated, so this pr was already in an event
				if exists {
					break
				}

				repo := event.GetRepo()
				// override state, since it could be newer
				currentPRState, err := c.getPullRequestState(ctx, repo.GetName(), pr.GetNumber())
				if err != nil {
					return collected, err
				}
				*pr.State = currentPRState

				prWithRepo := PullRequestWithRepository{
					PullRequest: pr,
					Repo:        repo,
				}
				collected.PullRequests[pr.GetID()] = prWithRepo
			}
		}

		if resp.NextPage == 0 {
			break pagingLoop
		}

		opts.Page = resp.NextPage
	}
	return collected, nil
}

func (c Client) getPullRequestState(ctx context.Context, repoName string, number int) (string, error) {
	owner, name := splitRepoName(repoName)
	pr, _, err := c.client.PullRequests.Get(ctx, owner, name, number)
	if err != nil {
		return "", fmt.Errorf("getting pullrequest %v/%v-#%v: %w", owner, name, number, err)
	}
	return pr.GetState(), nil
}

func (c Client) getIssueState(ctx context.Context, repoName string, number int) (string, error) {
	owner, name := splitRepoName(repoName)
	issue, _, err := c.client.Issues.Get(ctx, owner, name, number)
	if err != nil {
		return "", fmt.Errorf("getting issue %v/%v-#%v: %w", owner, name, number, err)
	}
	return issue.GetState(), nil
}

func getUsername(ctx context.Context, client *github.Client) (string, error) {
	// if user is empty, use the authenticated user
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("getting user: %w", err)
	}
	return *user.Login, nil
}

func splitRepoName(repoName string) (owner, name string) {
	owner = strings.Split(repoName, "/")[0]
	name = strings.Split(repoName, "/")[1]
	return
}
