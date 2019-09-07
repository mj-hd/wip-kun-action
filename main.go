package main

import (
	"context"

	"github.com/mjhd-devlion/wip-kun/pkg/checker"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

func main() {
	ctx := context.Background()
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	client := github.NewOctokit(ctx, conf)
	checker := checker.New(client)
	if err := checker.Check(ctx, conf.GithubSHA); err != nil {
		panic(err)
	}
}
