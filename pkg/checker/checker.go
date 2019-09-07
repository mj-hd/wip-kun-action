package checker

import (
	"context"

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
		return err
	}
	pr := prs[0]
	if checkPR(pr) {
		return errors.New("still WIP!")
	}
	commits, err := c.client.ListCommits(ctx, pr.Number)
	if err != nil {
		return err
	}
	for _, commit := range commits {
		if checkCommit(commit) {
		  return errors.New("failed to check commit")
		}
	}
	return nil
}

func (c *Checker) checkPR(pr github.PullRequest) bool {
	if (strings.HasPrefix(pr.Title, "WIP")) {
		return true
	}
	return false
}

func (c *Checker) checkCommit(commit github.Commit) bool {
	if (strings.HasPrefix(commit.Message, "fixup!")) {
		return true
	}
	if (strings.HasPrefix(commit.Message, "WIP")) {\
		return true
	}
	if (strings.HasPrefix(commit.Message, "wip")) {\
		return true
	}
	return false
}
