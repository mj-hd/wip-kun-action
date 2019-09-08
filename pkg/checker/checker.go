package checker

import (
	"context"
	"errors"
	"strings"

	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

type Checker struct {
	client github.Client
}

func New(client github.Client) *Checker {
	return &Checker{
		client: client,
	}
}

func (c *Checker) Check(ctx context.Context, sha string) error {
	prs, err := c.client.ListPullRequestsWithCommit(ctx, sha)
	if err != nil {
		panic(err)
	}
	pr := prs[0]
	if err := c.checkPR(pr); err != nil {
		return err
	}
	commits, err := c.client.ListCommits(ctx, pr.Number)
	if err != nil {
		panic(err)
	}
	for _, commit := range commits {
		if err := c.checkCommit(commit); err != nil {
			return err
		}
	}
	return nil
}

func (c *Checker) checkPR(pr github.PullRequest) error {
	title := strings.ToLower(pr.Title)
	if strings.HasPrefix(title, "wip") {
		return errors.New("PR title contains WIP")
	}
	return nil
}

func (c *Checker) checkCommit(commit github.Commit) error {
	message := strings.ToLower(commit.Message)
	if strings.HasPrefix(message, "fixup!") {
		return errors.New("fixup commit found")
	}
	if strings.HasPrefix(message, "wip") {
		return errors.New("WIP commit found")
	}
	return nil
}
