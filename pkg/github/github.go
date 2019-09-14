package github

import "context"

//go:generate mkdir -p mock
//go:generate mockgen -source=./github.go -package=mock -destination=./mock/mock.go

type Client interface {
	ListCommits(ctx context.Context, prNumber int) ([]Commit, error)
	AddLabel(ctx context.Context, prNumber int, label Label) error
	RemoveLabel(ctx context.Context, prNumber int, label Label) error
	UpdatePullRequestTitle(ctx context.Context, prNumber int, title string) error
}

type Label struct {
	Name string
}

type PullRequest struct {
	Number int
	Title  string
	Labels []Label
	Opened bool
}

type Commit struct {
	Message string
}

type EventType int

const (
	EVENT_TYPE_OPENED EventType = iota
	EVENT_TYPE_SYNCHRONIZED
	EVENT_TYPE_LABELED
	EVENT_TYPE_UNLABELED
	EVENT_TYPE_EDITED
)

// TODO simplify
type Event struct {
	Type  EventType
	PR    PullRequest
	Label *Label
	Title *string
}
