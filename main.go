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

	infos := client.Infos
	err = template.Render(outputFile, *templateFile, infos)
	if err != nil {
		log.Fatal(err)
	}

}
