package checker

import (
	"context"
	"errors"
	"strings"

	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

type Checker struct {
	client github.Client
	pr     github.PullRequest
}

func New(ctx context.Context, client github.Client, sha string) (*Checker, error) {
	prs, err := client.ListPullRequestsWithCommit(ctx, sha)
	if err != nil {
		return nil, err
	}
	return &Checker{
		client: client,
		pr:     prs[0],
	}, nil
}

func (c *Checker) Check(ctx context.Context, diff string) error {
	if err := c.checkPR(); err != nil {
		return err
	}
	commits, err := c.client.ListCommits(ctx, c.pr.Number)
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

func (c *Checker) EnsureLabel(ctx context.Context, wip bool, wipLabel string) error {
	if len(wipLabel) == 0 {
		return nil
	}
	if wip {
		return c.client.AddLabel(ctx, c.pr.Number, github.Label{Name: wipLabel})
	}
	return c.client.RemoveLabel(ctx, c.pr.Number, github.Label{Name: wipLabel})
}

func (c *Checker) checkPR() error {
	title := strings.ToLower(c.pr.Title)
	if strings.HasPrefix(title, "wip") {
		return errors.New("checker: PR title contains WIP")
	}
	return nil
}

func (c *Checker) checkCommit(commit github.Commit) error {
	message := strings.ToLower(commit.Message)
	if strings.HasPrefix(message, "fixup!") {
		return errors.New("checker: fixup commit found")
	}
	if strings.HasPrefix(message, "wip") {
		return errors.New("checker: WIP commit found")
	}
	return nil
}
