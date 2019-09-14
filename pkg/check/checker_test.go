package check

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
	"github.com/mjhd-devlion/wip-kun/pkg/github/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckerCheck(t *testing.T) {
	pr := github.PullRequest{
		Number: 1,
		Title:  "Title",
	}
	commit := github.Commit{
		Message: "Message",
	}
	testcases := map[string]struct {
		setup     func(client *mock.MockClient)
		sha       string
		expectErr bool
	}{
		"ready for review": {
			setup: func(client *mock.MockClient) {
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
				client.EXPECT().ListCommits(gomock.Any(), pr.Number).Return([]github.Commit{commit}, nil)
			},
			sha: "sha",
		},
		"PR title has WIP as its prefix": {
			setup: func(client *mock.MockClient) {
				pr := pr
				pr.Title = "WIP: Title"
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
			},
			sha:       "sha",
			expectErr: true,
		},
		"PR contains fixup commits": {
			setup: func(client *mock.MockClient) {
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
				commit := commit
				commit.Message = "fixup! feat: Message"
				client.EXPECT().ListCommits(gomock.Any(), pr.Number).Return([]github.Commit{commit}, nil)
			},
			sha:       "sha",
			expectErr: true,
		},
		"PR contains WIP commits": {
			setup: func(client *mock.MockClient) {
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
				commit := commit
				commit.Message = "WIP"
				client.EXPECT().ListCommits(gomock.Any(), pr.Number).Return([]github.Commit{commit}, nil)
			},
			sha:       "sha",
			expectErr: true,
		},
	}
	for title, testcase := range testcases {
		testcase := testcase
		test := func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mock.NewMockClient(ctrl)
			testcase.setup(client)
			checker, err := New(ctx, client, testcase.sha)
			require.NoError(t, err)
			err = checker.Check(ctx)
			if testcase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
		t.Run(title, test)
	}
}

func TestCheckerEnsureLabel(t *testing.T) {
	pr := github.PullRequest{
		Number: 1,
		Title:  "Title",
	}
	wipLabel := "work-in-progress"
	testcases := map[string]struct {
		setup    func(client *mock.MockClient)
		sha      string
		wip      bool
		wipLabel string
	}{
		"add wip": {
			setup: func(client *mock.MockClient) {
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
				client.EXPECT().AddLabel(gomock.Any(), pr.Number, github.Label{Name: wipLabel})
			},
			sha:      "sha",
			wip:      true,
			wipLabel: wipLabel,
		},
		"remove wip": {
			setup: func(client *mock.MockClient) {
				client.EXPECT().ListPullRequestsWithCommit(gomock.Any(), "sha").Return([]github.PullRequest{pr}, nil)
				client.EXPECT().RemoveLabel(gomock.Any(), pr.Number, github.Label{Name: wipLabel})
			},
			sha:      "sha",
			wip:      false,
			wipLabel: wipLabel,
		},
	}
	for title, testcase := range testcases {
		testcase := testcase
		test := func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mock.NewMockClient(ctrl)
			testcase.setup(client)
			checker, err := New(ctx, client, testcase.sha)
			require.NoError(t, err)
			err = checker.EnsureLabel(ctx, testcase.wip, testcase.wipLabel)
			require.NoError(t, err)
		}
		t.Run(title, test)
	}
}
