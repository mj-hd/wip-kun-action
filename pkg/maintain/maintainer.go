package maintain

import (
	"context"
	"errors"
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
	switch event.Type {
	case github.EVENT_TYPE_OPENED:
		err = m.maintainOpened(ctx, status)
	case github.EVENT_TYPE_EDITED:
		err = m.maintainEdited(ctx, status)
	case github.EVENT_TYPE_LABELED:
		err = m.maintainLabeled(ctx, status)
	case github.EVENT_TYPE_UNLABELED:
		err = m.maintainUnlabeled(ctx, status)
	case github.EVENT_TYPE_SYNCHRONIZED:
		err = m.maintainSynchronized(ctx, status)
	default:
		err = fmt.Errorf("maintain: unsupported event type %d", event.Type)
	}
	return
}

func (m *Maintainer) maintainOpened(ctx context.Context, status check.WIPStatus) error {
	// all
	_ = ctx
	_ = status
	return errors.New("TODO: not implemented yet")
}

func (m *Maintainer) maintainEdited(ctx context.Context, status check.WIPStatus) error {

	// label
	_ = ctx
	_ = status
	return errors.New("TODO: not implemented yet")
}

func (m *Maintainer) maintainLabeled(ctx context.Context, status check.WIPStatus) error {

	// title
	_ = ctx
	_ = status
	return errors.New("TODO: not implemented yet")
}

func (m *Maintainer) maintainUnlabeled(ctx context.Context, status check.WIPStatus) error {

	// title
	_ = ctx
	_ = status
	return errors.New("TODO: not implemented yet")
}

func (m *Maintainer) maintainSynchronized(ctx context.Context, status check.WIPStatus) error {

	// title, label
	_ = ctx
	_ = status
	return errors.New("TODO: not implemented yet")
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
	return errors.New("TODO: not implemented yet")
}

func (m *Maintainer) removeTitle(ctx context.Context, pr github.PullRequest) error {
	return errors.New("TODO: not implemented yet")
}
