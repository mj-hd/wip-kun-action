package github

import "context"

//go:generate mkdir -p mock
//go:generate mockgen -source=./github.go -package=mock -destination=./mock/mock.go

type Client interface {
	ListPullRequestsWithCommit(ctx context.Context, sha string) ([]PullRequest, error)
	ListCommits(ctx context.Context, pullRequestNumber int) ([]Commit, error)
}

type Label struct {
	Name string
}

type PullRequest struct {
	Number int
	Title  string
	Labels []Label
}

type Commit struct {
	SHA     string
	Message string
}
