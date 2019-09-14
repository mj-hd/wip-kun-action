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
	var f func(context.Context, github.PullRequest, check.WIPStatus) error
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
	return f(ctx, event.PR, status)
}

// TODO simplify
func (m *Maintainer) maintainOpened(ctx context.Context, pr github.PullRequest, status check.WIPStatus) error {
	if status.WIP() {
		if !status.HasWIPTitle {
			if err := m.addTitle(ctx, pr); err != nil {
				return err
			}
		}
		if !status.HasWIPLabel {
			if err := m.addLabel(ctx, pr); err != nil {
				return err
			}
		}
		return nil
	}
	if status.HasWIPTitle {
		if err := m.removeTitle(ctx, pr); err != nil {
			return err
		}
	}
	if status.HasWIPLabel {
		if err := m.removeLabel(ctx, pr); err != nil {
			return err
		}
	}
	return nil
}

func (m *Maintainer) maintainEdited(ctx context.Context, pr github.PullRequest, status check.WIPStatus) error {
	if status.WIP() {
		if !status.HasWIPLabel {
			if err := m.addLabel(ctx, pr); err != nil {
				return err
			}
		}
		return nil
	}
	if status.HasWIPLabel {
		if err := m.removeLabel(ctx, pr); err != nil {
			return err
		}
	}
	return nil
}

func (m *Maintainer) maintainLabeled(ctx context.Context, pr github.PullRequest, status check.WIPStatus) error {
	if status.WIP() {
		if !status.HasWIPTitle {
			if err := m.addTitle(ctx, pr); err != nil {
				return err
			}
		}
		return nil
	}
	if status.HasWIPTitle {
		if err := m.removeTitle(ctx, pr); err != nil {
			return err
		}
	}
	return nil
}

func (m *Maintainer) maintainUnlabeled(ctx context.Context, pr github.PullRequest, status check.WIPStatus) error {
	if status.WIP() {
		if !status.HasWIPTitle {
			if err := m.addTitle(ctx, pr); err != nil {
				return err
			}
		}
		return nil
	}
	if status.HasWIPTitle {
		if err := m.removeTitle(ctx, pr); err != nil {
			return err
		}
	}
	return nil
}

func (m *Maintainer) maintainSynchronized(ctx context.Context, pr github.PullRequest, status check.WIPStatus) error {
	if status.WIP() {
		if !status.HasWIPTitle {
			if err := m.addTitle(ctx, pr); err != nil {
				return err
			}
		}
		if !status.HasWIPLabel {
			if err := m.addLabel(ctx, pr); err != nil {
				return err
			}
		}
		return nil
	}
	if status.HasWIPTitle {
		if err := m.removeTitle(ctx, pr); err != nil {
			return err
		}
	}
	if status.HasWIPLabel {
		if err := m.removeLabel(ctx, pr); err != nil {
			return err
		}
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
