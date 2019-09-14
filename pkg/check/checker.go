package check

import (
	"context"
	"strings"

	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

type Checker struct {
	client github.Client
	config *config.Config
}

func NewChecker(client github.Client, config *config.Config) *Checker {
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

func (w WIPStatus) WIP() bool {
	return w.hasWIPTitle || w.hasWIPLabel || w.hasWIPCommits
}

func (c *Checker) Check(ctx context.Context, event github.Event) (status WIPStatus, err error) {

	status.hasWIPTitle = c.checkPR(event.PR)
	status.hasWIPLabel = c.checkLabels(event.PR)
	commits, err := c.client.ListCommits(ctx, event.PR.Number)
	if err != nil {
		return
	}
	status.hasWIPCommits = c.checkCommits(commits)
	return
}

func (c *Checker) checkPR(pr github.PullRequest) bool {
	title := strings.ToLower(pr.Title)
	if !strings.HasPrefix(title, c.config.WIPTitle) {
		return false
	}
	return true
}

func (c *Checker) checkCommits(commits []github.Commit) bool {
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
