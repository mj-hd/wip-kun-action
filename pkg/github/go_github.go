package github

import (
	"context"

	"github.com/google/go-github/v28/github"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"golang.org/x/oauth2"
)

type GoGithubClient struct {
	client *github.Client
	repo   string
	owner  string
}

func NewGoGithub(ctx context.Context, conf *config.Config) Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GithubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GoGithubClient{
		client: client,
		repo:   conf.GithubRepo,
		owner:  conf.GithubOwner,
	}
}

func (g *GoGithubClient) ListPullRequestsWithCommit(ctx context.Context, sha string) ([]PullRequest, error) {
	prs, _, err := g.client.PullRequests.ListPullRequestsWithCommit(ctx, g.owner, g.repo, sha, nil)
	if err != nil {
		return nil, err
	}
	return toPullRequests(prs), nil
}

func (g *GoGithubClient) ListCommits(ctx context.Context, pullRequestNumber int) ([]Commit, error) {
	commits, _, err := g.client.PullRequests.ListCommits(ctx, g.owner, g.repo, pullRequestNumber, nil)
	if err != nil {
		return nil, err
	}
	return toCommits(commits), nil
}

func (g *GoGithubClient) filterLabelsByName(labels []*github.Label, name string) []*github.Label {
	results := make([]*github.Label, 0, len(labels))
	for _, label := range labels {
		if label.GetName() == name {
			continue
		}
		results = append(results, label)
	}
	return results
}

func (g *GoGithubClient) AddLabel(ctx context.Context, pullRequestNumber int, label Label) error {
	pr, _, err := g.client.PullRequests.Get(ctx, g.owner, g.repo, pullRequestNumber)
	if err != nil {
		return err
	}
	labels := make([]*github.Label, 0, len(pr.Labels))
	for _, label := range pr.Labels {
		if label.GetName() == label.Name {
			continue
		}
		labels = append(labels, label)
	}
	pr.Labels = g.filterLabelsByName(pr.Labels, label.Name)
	pr.Labels = append(pr.Labels, toGithubLabel(label))
	_, _, err = g.client.PullRequests.Edit(ctx, g.owner, g.repo, pullRequestNumber, pr)
	return err
}

func (g *GoGithubClient) RemoveLabel(ctx context.Context, pullRequestNumber int, label Label) error {
	pr, _, err := g.client.PullRequests.Get(ctx, g.owner, g.repo, pullRequestNumber)
	if err != nil {
		return err
	}
	pr.Labels = g.filterLabelsByName(pr.Labels, label.Name)
	_, _, err = g.client.PullRequests.Edit(ctx, g.owner, g.repo, pullRequestNumber, pr)
	return err
}

func toPullRequests(prs []*github.PullRequest) []PullRequest {
	result := make([]PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = toPullRequest(pr)
	}
	return result
}

func toPullRequest(pr *github.PullRequest) PullRequest {
	return PullRequest{
		Number: pr.GetNumber(),
		Title:  pr.GetTitle(),
		Labels: toLabels(pr.Labels),
	}
}

func toLabels(labels []*github.Label) []Label {
	result := make([]Label, len(labels))
	for i, label := range labels {
		result[i] = toLabel(label)
	}
	return result
}

func toLabel(label *github.Label) Label {
	return Label{Name: label.GetName()}
}

func toCommits(commits []*github.RepositoryCommit) []Commit {
	result := make([]Commit, len(commits))
	for i, commit := range commits {
		result[i] = toCommit(commit)
	}
	return result
}

func toCommit(commit *github.RepositoryCommit) Commit {
	return Commit{
		SHA:     commit.Commit.GetSHA(),
		Message: commit.Commit.GetMessage(),
	}
}
