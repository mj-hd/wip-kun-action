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
		panic(err)
	}
	client := github.NewGoGithub(ctx, conf)
	checker, err := checker.New(ctx, client, conf.GithubSHA)
	if err != nil {
		panic(err)
	}
	wip := false
	if err := checker.Check(ctx); err != nil {
		wip = true
		fmt.Println(err)
	}
	if err := checker.EnsureLabel(ctx, wip, conf.WIPLabel); err != nil {
		panic(err)
	}
	if wip {
		os.Exit(1)
	}
}
