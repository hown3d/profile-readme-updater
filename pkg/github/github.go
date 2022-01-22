package github

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type githubClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

type Client struct {
	client githubClient
	user   string
}

type Info struct {
	ContributionInfo ContributionInfo
	Repository       Repo
}

type ContributionInfo struct {
	Issues       []RepoItem
	PullRequests []RepoItem
}

type RepoItem struct {
	ID    string
	Title string
	URL   *url.URL
	State string
}

type Repo struct {
	ID          string
	Stars       int
	Name        string
	Description string
	URL         *url.URL
	ImageURL    *url.URL
}

func NewGithubGraphQLClient() githubClient {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubv4.NewClient(httpClient)
}

func NewClient(user string, client githubClient) Client {
	return Client{
		client: client,
		user:   user,
	}
}

func (c Client) GetInfos(ctx context.Context) ([]Info, error) {
	infos := []Info{}
	repos, err := c.getAllContributedRepos(ctx)
	if err != nil {
		return []Info{}, err
	}
	for _, repo := range repos {
		contribInfo, err := c.getIssuesAndPullRequests(ctx, repo.Name)
		if err != nil {
			return []Info{}, err
		}
		infos = append(infos, Info{
			ContributionInfo: contribInfo,
			Repository:       repo,
		})
	}
	return infos, nil
}

/*
{
	search(query: "author:$author repo:$repo", first: 10, type: ISSUE) {
		edges {
			node {
				... on PullRequest {
					id
					title
					url
					state
				}
				... on Issue {
					id
					title
					url
					state

				}
			}
		}
	}
}
*/
type getIssuesAndPullRequestsQuery struct {
	Search struct {
		Edges []struct {
			Node struct {
				TypeName    githubv4.RepositoryContributionType `graphql:"__typename"`
				PullRequest struct {
					ID    githubv4.String
					Title githubv4.String
					URL   githubv4.URI
					State githubv4.String
				} `graphql:"... on PullRequest"`
				Issue struct {
					ID    githubv4.String
					Title githubv4.String
					URL   githubv4.URI
					State githubv4.String
				} `graphql:"... on Issue"`
			}
		}
	} `graphql:"search(query: $query, first: 10, type: ISSUE"`
}

func (c Client) getIssuesAndPullRequests(ctx context.Context, repoName string) (ContributionInfo, error) {
	query := &getIssuesAndPullRequestsQuery{}
	variables := map[string]interface{}{
		"query": githubv4.String(fmt.Sprintf("author: %s repo: %s", c.user, repoName)),
	}
	err := c.client.Query(ctx, query, variables)
	if err != nil {
		return ContributionInfo{}, err
	}
	info := ContributionInfo{
		Issues:       []RepoItem{},
		PullRequests: []RepoItem{},
	}
	for _, edge := range query.Search.Edges {
		switch edge.Node.TypeName {
		case githubv4.RepositoryContributionTypeIssue:
			info.Issues = append(info.Issues, RepoItem{
				ID:    string(edge.Node.Issue.ID),
				Title: string(edge.Node.Issue.Title),
				URL:   edge.Node.Issue.URL.URL,
				State: string(edge.Node.Issue.State),
			})
		case githubv4.RepositoryContributionTypePullRequest:
			info.PullRequests = append(info.PullRequests, RepoItem{
				ID:    string(edge.Node.PullRequest.ID),
				Title: string(edge.Node.PullRequest.Title),
				URL:   edge.Node.PullRequest.URL.URL,
				State: string(edge.Node.PullRequest.State),
			})
		default:
			log.Printf("Error: Unknown type name: %s", edge.Node.TypeName)
		}
	}
	return info, nil
}

/*{
    user(login: "${login}") {
      repositoriesContributedTo(includeUserRepositories: false, first: 10, privacy: PUBLIC) {
        edges {
          node {
            id
            nameWithOwner
            shortDescriptionHTML(limit: 120)
            stargazers {
              totalCount
            }
            url
            openGraphImageUrl
          }
        }
      }
    }
  }
*/
type getAllContributedReposQuery struct {
	User struct {
		RepositoriesContributedTo struct {
			Edges []struct {
				Node struct {
					ID                   githubv4.String
					NameWithOwner        githubv4.String
					ShortDescriptionHTML githubv4.HTML `graphql:"shortDescriptionHTML(limit: 120)"`
					Stargazers           struct {
						TotalCount githubv4.Int
					}
					URL               githubv4.URI
					OpenGraphImageURL githubv4.URI
				}
			}
		} `graphql:"repositoriesContributedTo(first: 10, privacy: PUBLIC)"`
	} `graphql:"user(login: $login)"`
}

func (c Client) getAllContributedRepos(ctx context.Context) ([]Repo, error) {
	query := &getAllContributedReposQuery{}

	variables := map[string]interface{}{
		"login": githubv4.String(c.user),
	}

	err := c.client.Query(ctx, query, variables)
	if err != nil {
		return []Repo{}, err
	}

	repos := []Repo{}
	for _, edge := range query.User.RepositoriesContributedTo.Edges {
		repos = append(repos, Repo{
			ID:          string(edge.Node.ID),
			Stars:       int(edge.Node.Stargazers.TotalCount),
			Name:        string(edge.Node.NameWithOwner),
			Description: string(edge.Node.ShortDescriptionHTML),
			URL:         edge.Node.URL.URL,
			ImageURL:    edge.Node.OpenGraphImageURL.URL,
		})
	}
	return repos, nil
}

func (c Client) getLanguage() {
	panic("not implemented")
}
