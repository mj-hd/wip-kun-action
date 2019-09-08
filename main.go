package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mjhd-devlion/wip-kun/pkg/checker"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

func main() {
	ctx := context.Background()
	conf, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client := github.NewGoGithub(ctx, conf)
	checker := checker.New(client)
	if err := checker.Check(ctx, conf.GithubSHA); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
