package main

import (
	"context"
	"io/ioutil"
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
	client := github.NewGoGithub(ctx, conf.GithubToken, conf.GithubOwner, conf.GithubRepo)
	file, err := os.Open(conf.GithubEventPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	event, err := github.NewEvent(conf.GithubEventName, bytes)
	if err != nil {
		panic(err)
	}
	checker := checker.New(ctx, client)
	wip, err := checker.Check(ctx, event, conf.GithubRef)
	if err != nil {
		panic(err)
	}
}
