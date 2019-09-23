package check

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
	"github.com/mjhd-devlion/wip-kun/pkg/github/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckerCheck(t *testing.T) {
	config := &config.Config{
		WIPLabel:   "work-in-progress",
		WIPTitle:   "WIP: ",
		WIPCommits: "fixup!",
	}
	testcases := map[string]struct {
		event  github.Event
		setup  func(mock *mock.MockClient)
		expect WIPStatus
	}{
		"has WIP title": {
			event: github.Event{
				PR: github.PullRequest{
					Number: 1,
					Title:  config.WIPTitle + "hello",
				},
			},
			setup: func(m *mock.MockClient) {
				m.EXPECT().ListCommits(gomock.Any(), 1).Return(nil, nil)
			},
			expect: WIPStatus{
				HasWIPTitle: true,
			},
		},
		"has fixup commits": {
			event: github.Event{
				PR: github.PullRequest{
					Number: 1,
				},
			},
			setup: func(m *mock.MockClient) {
				m.EXPECT().ListCommits(gomock.Any(), 1).Return([]github.Commit{
					{
						Message: config.WIPCommits,
					},
				}, nil)
			},
			expect: WIPStatus{
				HasWIPCommits: true,
			},
		},
		"has WIP label": {
			event: github.Event{
				PR: github.PullRequest{
					Number: 1,
					Labels: []github.Label{
						{
							Name: config.WIPLabel,
						},
					},
				},
			},
			setup: func(m *mock.MockClient) {
				m.EXPECT().ListCommits(gomock.Any(), 1).Return([]github.Commit{}, nil)
			},
			expect: WIPStatus{
				HasWIPLabel: true,
			},
		},
	}
	for title, testcase := range testcases {
		testcase := testcase
		test := func(t *testing.T) {
			ctx := context.Background()
			c := gomock.NewController(t)
			defer c.Finish()
			m := mock.NewMockClient(c)
			testcase.setup(m)
			checker := NewChecker(m, config)
			status, err := checker.Check(ctx, testcase.event)
			require.NoError(t, err)
			require.Equal(t, testcase.expect, status)
		}
		t.Run(title, test)
	}
}

func TestWIPStatus(t *testing.T) {
	status := WIPStatus{}
	require.Equal(t, status.WIP(), false)
	status.HasWIPTitle = true
	require.Equal(t, status.WIP(), true)
}
