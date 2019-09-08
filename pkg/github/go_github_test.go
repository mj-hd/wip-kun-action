package github

import (
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/stretchr/testify/require"
)

func TestToPullRequests(t *testing.T) {
	prTitle1 := "PR1"
	prTitle2 := "PR2"
	prNumber1 := 1
	prNumber2 := 2
	target := []*github.PullRequest{
		{
			Number: &prNumber1,
			Title:  &prTitle1,
			Labels: nil,
		},
	}
	expect := []PullRequest{
		{
			Number: prNumber2,
			Title:  prTitle2,
			Labels: nil,
		},
	}
	require.Equal(t, toPullRequests(target), expect)
}

func TestToCommits(t *testing.T) {
	commitSHA1 := "sha1"
	commitSHA2 := "sha2"
	commitMessage1 := "message1"
	commitMessage2 := "message2"
	target := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				SHA:     &commitSHA1,
				Message: &commitMessage1,
			},
		},
		{
			Commit: &github.Commit{
				SHA:     &commitSHA2,
				Message: &commitMessage2,
			},
		},
	}
	expect := []Commit{
		{
			SHA:     commitSHA1,
			Message: commitMessage1,
		},
		{
			SHA:     commitSHA2,
			Message: commitMessage2,
		},
	}
	require.Equal(t, toCommits(target), expect)
}
