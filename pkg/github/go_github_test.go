package github

import (
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/stretchr/testify/require"
)

func TestToPullRequests(t *testing.T) {
	prTitle1 := "PR1"
	prTitle2 := "PR2"
	labelName := "WIP"
	prNumber1 := 1
	prNumber2 := 2
	opened := "open"
	closed := "close"
	target := []*github.PullRequest{
		{
			Number: &prNumber1,
			Title:  &prTitle1,
			Labels: []*github.Label{
				{
					Name: &labelName,
				},
			},
			State: &opened,
		},
		{
			Number: &prNumber2,
			Title:  &prTitle2,
			Labels: nil,
			State:  &closed,
		},
	}
	expect := []PullRequest{
		{
			Number: prNumber1,
			Title:  prTitle1,
			Labels: []Label{
				{Name: labelName},
			},
			Opened: true,
		},
		{
			Number: prNumber2,
			Title:  prTitle2,
			Labels: []Label{},
			Opened: false,
		},
	}
	require.Equal(t, expect, toPullRequests(target))
}

func TestToCommits(t *testing.T) {
	commitMessage1 := "message1"
	commitMessage2 := "message2"
	target := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: &commitMessage1,
			},
		},
		{
			Commit: &github.Commit{
				Message: &commitMessage2,
			},
		},
	}
	expect := []Commit{
		{
			Message: commitMessage1,
		},
		{
			Message: commitMessage2,
		},
	}
	require.Equal(t, expect, toCommits(target))
}

func TestNewEvent(t *testing.T) {
	from_title := "from_title"
	testcases := map[string]struct {
		data   []byte
		err    bool
		expect Event
	}{
		"unsupported event type": {
			data: []byte(`
{
  "action": "assigned",
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			err: true,
		},
		"opened": {
			data: []byte(`
{
  "action": "opened",
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_OPENED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
			},
		},
		"reopened": {
			data: []byte(`
{
  "action": "reopened",
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_OPENED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
			},
		},
		"edited": {
			data: []byte(`
{
  "action": "edited",
  "changes": {
    "title": {
      "from": "from_title"
    }
  },
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_EDITED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
				ChangedTitle: &from_title,
			},
		},
		"labeled": {
			data: []byte(`
{
  "action": "labeled",
  "label": {
    "name": "work-in-progress"
  },
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_LABELED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
				ChangedLabel: &Label{
					Name: "work-in-progress",
				},
			},
		},
		"unlabeled": {
			data: []byte(`
{
  "action": "unlabeled",
  "label": {
    "name": "work-in-progress"
  },
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_UNLABELED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
				ChangedLabel: &Label{
					Name: "work-in-progress",
				},
			},
		},
		"synchronize": {
			data: []byte(`
{
  "action": "synchronize",
  "number": 1,
  "pull_request": {
    "body": "body",
    "labels": [],
    "number": 1,
    "state": "open",
    "title": "title"
  }
}
			`),
			expect: Event{
				Type: EVENT_TYPE_SYNCHRONIZED,
				PR: PullRequest{
					Number: 1,
					Title:  "title",
					Labels: []Label{},
					Opened: true,
				},
			},
		},
	}
	for title, testcase := range testcases {
		testcase := testcase
		test := func(t *testing.T) {
			actual, err := NewEvent("pull_request", testcase.data)
			if testcase.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testcase.expect, actual)
		}
		t.Run(title, test)
	}
}
