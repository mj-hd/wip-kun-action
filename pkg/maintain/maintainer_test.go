package maintain

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mjhd-devlion/wip-kun/pkg/check"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
	"github.com/mjhd-devlion/wip-kun/pkg/github/mock"
	"github.com/stretchr/testify/require"
)

func TestMaintainerMaintain(t *testing.T) {
	title := "hello"
	config := &config.Config{
		WIPTitle:   "WIP: ",
		WIPLabel:   "work-in-progress",
		WIPCommits: "fixup!",
	}
	testcases := map[string]struct {
		setup  func(m *mock.MockClient)
		event  github.Event
		status check.WIPStatus
		expect bool
	}{
		"opened with WIP commits": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, config.WIPTitle+title)
				m.EXPECT().AddLabel(gomock.Any(), 1, github.Label{
					Name: config.WIPLabel,
				})
			},
			event: github.Event{
				Type: github.EVENT_TYPE_OPENED,
				PR: github.PullRequest{
					Number: 1,
					Title:  title,
					Opened: true,
				},
			},
			status: check.WIPStatus{
				HasWIPCommits: true,
			},
			expect: true,
		},
		"opened with WIP title": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().AddLabel(gomock.Any(), 1, github.Label{
					Name: config.WIPLabel,
				})
			},
			event: github.Event{
				Type: github.EVENT_TYPE_OPENED,
				PR: github.PullRequest{
					Number: 1,
					Title:  config.WIPTitle + title,
					Opened: true,
				},
			},
			status: check.WIPStatus{
				HasWIPTitle: true,
			},
			expect: true,
		},
		"opened with WIP label": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, config.WIPTitle+title)
			},
			event: github.Event{
				Type: github.EVENT_TYPE_OPENED,
				PR: github.PullRequest{
					Number: 1,
					Title:  title,
					Labels: []github.Label{
						{Name: config.WIPLabel},
					},
					Opened: true,
				},
			},
			status: check.WIPStatus{
				HasWIPLabel: true,
			},
			expect: true,
		},
		"edited unrelated field": {
			setup: func(m *mock.MockClient) {
			},
			event: github.Event{
				Type: github.EVENT_TYPE_EDITED,
				PR: github.PullRequest{
					Opened: true,
				},
			},
			status: check.WIPStatus{},
		},
		"edited title with WIP": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().AddLabel(gomock.Any(), 1, github.Label{
					Name: config.WIPLabel,
				})
			},
			event: github.Event{
				Type: github.EVENT_TYPE_EDITED,
				PR: github.PullRequest{
					Number: 1,
					Title:  config.WIPTitle + title,
					Opened: true,
				},
				ChangedTitle: &title,
			},
			status: check.WIPStatus{
				HasWIPTitle: true,
			},
			expect: true,
		},
		"edited title with not WIP": {
			setup: func(m *mock.MockClient) {
			},
			event: github.Event{
				Type: github.EVENT_TYPE_EDITED,
				PR: github.PullRequest{
					Number: 1,
					Title:  "world",
					Opened: true,
				},
				ChangedTitle: &title,
			},
			status: check.WIPStatus{},
		},
		"unrelated labeled": {
			setup: func(m *mock.MockClient) {
			},
			event: github.Event{
				Type: github.EVENT_TYPE_LABELED,
				PR: github.PullRequest{
					Opened: true,
				},
				ChangedLabel: &github.Label{
					Name: "unrelated-label",
				},
			},
			status: check.WIPStatus{},
		},
		"labeled WIP": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, config.WIPTitle+title)
			},
			event: github.Event{
				Type: github.EVENT_TYPE_LABELED,
				PR: github.PullRequest{
					Number: 1,
					Title:  title,
					Labels: []github.Label{
						{Name: config.WIPLabel},
					},
					Opened: true,
				},
				ChangedLabel: &github.Label{
					Name: config.WIPLabel,
				},
			},
			status: check.WIPStatus{
				HasWIPLabel: true,
			},
			expect: true,
		},
		"unrelated unlabeled": {
			setup: func(m *mock.MockClient) {
			},
			event: github.Event{
				Type: github.EVENT_TYPE_UNLABELED,
				PR: github.PullRequest{
					Opened: true,
				},
				ChangedLabel: &github.Label{
					Name: "unrelated-label",
				},
			},
			status: check.WIPStatus{},
		},
		"unlabeled WIP": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, title)
			},
			event: github.Event{
				Type: github.EVENT_TYPE_UNLABELED,
				PR: github.PullRequest{
					Number: 1,
					Title:  config.WIPTitle + title,
					Opened: true,
				},
				ChangedLabel: &github.Label{
					Name: config.WIPLabel,
				},
			},
			status: check.WIPStatus{
				HasWIPTitle: true,
			},
		},
		"synchronized with WIP commits": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, config.WIPTitle+title)
				m.EXPECT().AddLabel(gomock.Any(), 1, github.Label{
					Name: config.WIPLabel,
				})
			},
			event: github.Event{
				Type: github.EVENT_TYPE_SYNCHRONIZED,
				PR: github.PullRequest{
					Number: 1,
					Title:  title,
					Opened: true,
				},
			},
			status: check.WIPStatus{
				HasWIPCommits: true,
			},
			expect: true,
		},
		"synchronized without WIP commits": {
			setup: func(m *mock.MockClient) {
				m.EXPECT().UpdatePullRequestTitle(gomock.Any(), 1, title)
				m.EXPECT().RemoveLabel(gomock.Any(), 1, github.Label{
					Name: config.WIPLabel,
				})
			},
			event: github.Event{
				Type: github.EVENT_TYPE_SYNCHRONIZED,
				PR: github.PullRequest{
					Number: 1,
					Title:  config.WIPTitle + title,
					Labels: []github.Label{
						{Name: config.WIPLabel},
					},
					Opened: true,
				},
			},
			status: check.WIPStatus{
				HasWIPTitle: true,
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
			maintainer := NewMaintainer(m, config)
			wip, err := maintainer.Maintain(ctx, testcase.event, testcase.status)
			require.NoError(t, err)
			require.Equal(t, testcase.expect, wip)
		}
		t.Run(title, test)
	}
}
