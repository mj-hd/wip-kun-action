package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mjhd-devlion/wip-kun/pkg/check"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
	"github.com/mjhd-devlion/wip-kun/pkg/maintain"
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
	checker := check.NewChecker(client, conf)
	maintainer := maintain.NewMaintainer(client, conf)
	status, err := checker.Check(ctx, event)
	if err != nil {
		panic(err)
	}
	err = maintainer.Maintain(ctx, event, status)
	if err != nil {
		panic(err)
	}
	if status.WIP() {
		os.Exit(1)
	}
}
