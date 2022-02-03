package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hown3d/profile-readme-updater/pkg/github"
	"github.com/hown3d/profile-readme-updater/pkg/template"
)

var (
	templateFile   *string = flag.String("template-file", "", "path to the template file")
	outputFilePath *string = flag.String("out", "", "path to the output")
)

func main() {
	flag.Parse()

	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	githubClient, err := github.NewGithubClient()
	if err != nil {
		log.Fatal(err)
	}
	client, err := github.NewClient(githubClient)
	if err != nil {
		log.Fatal(fmt.Errorf("creating client: %v", err))
	}

	now := time.Now()
	oneMonthEarlier := now.AddDate(0, -1, 0)
	err = client.GetContributions(context.Background(), oneMonthEarlier)
	if err != nil {
		log.Fatal(fmt.Errorf("get contributions: %v", err))
	}

	err = template.Render(outputFile, *templateFile, client.CollectedEvents())
	if err != nil {
		log.Fatal(err)
	}
	// w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	// fmt.Fprintf(w, "Number of issue events: %v\n", len(collectedEvents.Issues))
	// for _, issueWithRepo := range collectedEvents.Issues {
	// 	issue := issueWithRepo.Issue
	// 	fmt.Fprintf(w, "Title: %v\tRepo: %v\tStatus: %v\tDate: %v\tURL: %v\n", issue.GetTitle(), issueWithRepo.Repo.GetName(), issue.GetState(), issue.GetCreatedAt(), issue.GetHTMLURL())
	// }

	// fmt.Fprintf(w, "Number of pr events: %v\n", len(collectedEvents.PullRequests))
	// for _, prWithRepo := range collectedEvents.PullRequests {
	// 	pr := prWithRepo.PullRequest
	// 	fmt.Println(pr.GetMerged())
	// }
	// err = w.Flush()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
