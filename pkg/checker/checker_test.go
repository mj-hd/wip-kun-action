package checker

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mjhd-devlion/wip-kun/pkg/github/mock"
	"github.com/stretchr/testify/require"
)

func TestChecker(t *testing.T) {
	testcases := map[string]struct {
		setup     func(client *mock.MockClient)
		sha       string
		expectErr bool
	}{
		"ready for review": {
			setup: func(client *mock.MockClient) {

			},
			sha: "sha",
		},
		"PR title has WIP as its prefix": {
			setup: func(client *mock.MockClient) {
			},
			sha:       "sha",
			expectErr: true,
		},
		"PR contains fixup commits": {
			setup: func(client *mock.MockClient) {
			},
			sha:       "sha",
			expectErr: true,
		},
		"PR contains WIP commits": {
			setup: func(client *mock.MockClient) {
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
			checker := New(client)
			require.Equal(t, checker.Check(ctx, testcase.sha) != nil, testcase.expectErr)
		}
		t.Run(title, test)
	}
}
