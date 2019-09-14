package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type GoGithubClient struct {
	client *github.Client
	repo   string
	owner  string
}

func NewGoGithub(ctx context.Context, token, owner, repo string) Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GoGithubClient{
		client: client,
		repo:   repo,
		owner:  owner,
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

func (g *GoGithubClient) AddLabel(ctx context.Context, pullRequestNumber int, label Label) error {
	_, _, err := g.client.Issues.AddLabelsToIssue(ctx, g.owner, g.repo, pullRequestNumber, []string{label.Name})
	return err
}

func (g *GoGithubClient) RemoveLabel(ctx context.Context, pullRequestNumber int, label Label) error {
	_, err := g.client.Issues.RemoveLabelForIssue(ctx, g.owner, g.repo, pullRequestNumber, label.Name)
	return err
}

func NewEvent(typ string, data []byte) (Event, error) {
	event, err := github.ParseWebHook(typ, data)
	if err != nil {
		return Event{}, err
	}
	switch e := event.(type) {
	case *github.PullRequestEvent:
		return toPullRequestEvent(e)
	}
	return Event{}, errors.New("github: unsupported event type")
}

func toPullRequestEvent(e *github.PullRequestEvent) (Event, error) {
	typ, err := toEventType(e.GetAction())
	return Event{
		Type: typ,
		PR:   toPullRequest(e.GetPullRequest()),
	}, err
}

func toEventType(action string) (EventType, error) {
	switch action {
	case "opened", "reopened":
		return EVENT_TYPE_OPENED, nil
	case "edited":
		return EVENT_TYPE_EDITED, nil
	case "labeled":
		return EVENT_TYPE_LABELED, nil
	case "unlabeled":
		return EVENT_TYPE_UNLABELED, nil
	case "synchronized":
		return EVENT_TYPE_SYNCHRONIZED, nil
	}
	return 0, errors.New("github: unsupported pull request event")
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
