package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hown3d/profile-readme-updater/pkg/github"
)

var (
	user *string = flag.String("user", "", "github username to use")
)

func main() {
	flag.Parse()
	client := github.NewClient(*user, github.NewGithubGraphQLClient())
	infos, err := client.GetInfos(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("get infos: %v", err))
	}
	for _, info := range infos {
		fmt.Println(info)
	}
}
