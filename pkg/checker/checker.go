package checker

import (
	"context"
	"errors"
	"strings"

	go_github "github.com/google/go-github/v28/github"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

type Checker struct {
	client github.Client
	config *config.Config
}

func New(ctx context.Context, client github.Client, config *config.Config) *Checker {
	return &Checker{
		client: client,
		config: config,
	}
}

type WIPStatus struct {
	hasWIPTitle   bool
	hasWIPCommits bool
	hasWIPLabel   bool
}

func (c *Checker) Check(ctx context.Context, event github.Event, ref string) (status WIPStatus, err error) {

	var pr 
	status.hasWIPTitle = c.checkPR(pr)
	status.hasWIPLabel = c.checkLabels(pr)
	commits, err := c.client.ListCommits(ctx, prNumber)
	if err != nil {
		return
	}
	status.hasWIPCommits = c.checkCommits(commits)
	return
}

func (c *Checker) checkPR(pr github.PullRequest) bool {
	title := strings.ToLower(pr.Title)
	if !strings.HasPrefix(pr.Title, c.config.WIPTitle) {
		return false
	}
	return true
}

func (c *Checker) checkCommits(commit []github.Commit) bool {
	for _, commit := range commits {
		if !c.checkCommit(commit) {
			continue
		}
		return true
	}
	return false
}

func (c *Checker) checkCommit(commit github.Commit) bool {
	message := strings.ToLower(commit.Message)
	for _, prefix := range c.config.WIPCommits() {
		if !strings.HasPrefix(message, prefix) {
			continue
		}
		return true
	}
	return false
}

func (c *Checker) checkLabels(pr github.PullRequest) bool {
	for _, label := range pr.Labels {
		if label.Name != c.config.WIPLabel {
			continue
		}
		return true
	}
	return false
}
