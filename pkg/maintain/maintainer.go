package maintain

import (
	"context"
	"fmt"

	"github.com/mjhd-devlion/wip-kun/pkg/check"
	"github.com/mjhd-devlion/wip-kun/pkg/config"
	"github.com/mjhd-devlion/wip-kun/pkg/github"
)

type Maintainer struct {
	client github.Client
	config *config.Config
}

func NewMaintainer(client github.Client, config *config.Config) *Maintainer {
	return &Maintainer{
		client: client,
		config: config,
	}
}

func (m *Maintainer) Maintain(ctx context.Context, event github.Event, status check.WIPStatus) (err error) {
	if !event.PR.Opened {
		return nil
	}
	var f func(context.Context, github.Event, check.WIPStatus) error
	switch event.Type {
	case github.EVENT_TYPE_OPENED:
		f = m.maintainOpened
	case github.EVENT_TYPE_EDITED:
		f = m.maintainEdited
	case github.EVENT_TYPE_LABELED:
		f = m.maintainLabeled
	case github.EVENT_TYPE_UNLABELED:
		f = m.maintainUnlabeled
	case github.EVENT_TYPE_SYNCHRONIZED:
		f = m.maintainSynchronized
	default:
		return fmt.Errorf("maintain: unsupported event type %d", event.Type)
	}
	return f(ctx, event, status)
}

func (m *Maintainer) maintainOpened(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if status.WIP() {
		return m.wip(ctx, e, status)
	}
	return m.unwip(ctx, e, status)
}

func (m *Maintainer) maintainEdited(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if e.ChangedTitle == nil {
		return nil
	}
	if status.HasWIPCommits || status.HasWIPTitle {
		return m.wip(ctx, e, status)
	}
	return m.unwip(ctx, e, status)
}

func (m *Maintainer) maintainLabeled(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if e.ChangedLabel == nil || e.ChangedLabel.Name != m.config.WIPLabel {
		return nil
	}
	return m.wip(ctx, e, status)
}

func (m *Maintainer) maintainUnlabeled(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if e.ChangedLabel == nil || e.ChangedLabel.Name != m.config.WIPLabel {
		return nil
	}
	if status.HasWIPCommits {
		return m.wip(ctx, e, status)
	}
	return m.unwip(ctx, e, status)
}

func (m *Maintainer) maintainSynchronized(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if status.HasWIPCommits {
		return m.wip(ctx, e, status)
	}
	return m.unwip(ctx, e, status)
}

func (m *Maintainer) unwip(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if status.HasWIPLabel {
		if err := m.removeLabel(ctx, e.PR); err != nil {
			return err
		}
	}
	if status.HasWIPTitle {
		return m.removeTitle(ctx, e.PR)
	}
	return nil
}

func (m *Maintainer) wip(ctx context.Context, e github.Event, status check.WIPStatus) error {
	if !status.HasWIPLabel {
		if err := m.addLabel(ctx, e.PR); err != nil {
			return err
		}
	}
	if !status.HasWIPTitle {
		return m.addTitle(ctx, e.PR)
	}
	return nil
}

func (m *Maintainer) addLabel(ctx context.Context, pr github.PullRequest) error {
	return m.client.AddLabel(ctx, pr.Number, github.Label{
		Name: m.config.WIPLabel,
	})
}

func (m *Maintainer) removeLabel(ctx context.Context, pr github.PullRequest) error {
	return m.client.RemoveLabel(ctx, pr.Number, github.Label{
		Name: m.config.WIPLabel,
	})
}

func (m *Maintainer) addTitle(ctx context.Context, pr github.PullRequest) error {
	title := m.config.WIPTitle + pr.Title
	return m.client.UpdatePullRequestTitle(ctx, pr.Number, title)
}

func (m *Maintainer) removeTitle(ctx context.Context, pr github.PullRequest) error {
	title := pr.Title[len(m.config.WIPTitle):]
	return m.client.UpdatePullRequestTitle(ctx, pr.Number, title)
}
