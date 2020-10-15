package docker_test

import (
	"github.com/thepwagner/action-update-docker/docker"
	"github.com/thepwagner/action-update/updater"
)

//go:generate mockery --outpkg docker_test --output . --testonly --name TagLister --structname mockTagLister --filename mocktaglister_test.go

type testFactory struct{}

func (u *testFactory) NewUpdater(root string) updater.Updater {
	return docker.NewUpdater(root)
}

func updaterFactory() updater.Factory {
	return &testFactory{}
}
