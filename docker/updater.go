package docker

import (
	"context"

	"github.com/thepwagner/action-update/updater"
)

type Updater struct {
	root string

	tags TagLister
}

var _ updater.Updater = (*Updater)(nil)

func NewUpdater(root string, opts ...UpdaterOpt) *Updater {
	u := &Updater{
		root: root,
		tags: NewRemoteTagLister(),
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

type UpdaterOpt func(*Updater)

func WithTagsLister(tags TagLister) UpdaterOpt {
	return func(u *Updater) {
		u.tags = tags
	}
}

func (u *Updater) ApplyUpdate(ctx context.Context, update updater.Update) error {
	panic("implement me")
}
