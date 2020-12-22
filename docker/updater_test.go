package docker_test

import (
	"github.com/thepwagner/action-update-docker/docker"
	"github.com/thepwagner/action-update/updater"
)

//go:generate mockery --outpkg docker_test --output . --testonly --name TagLister --structname mockTagLister --filename mocktaglister_test.go
//go:generate mockery --outpkg docker_test --output . --testonly --name ImagePinner --structname mockImagePinner --filename mockimagepinner_test.go

type testFactory struct {
	opts []docker.UpdaterOpt
}

func (u *testFactory) NewUpdater(root string) updater.Updater {
	return docker.NewUpdater(root, u.opts...)
}

func updaterFactory(opts ...docker.UpdaterOpt) updater.Factory {
	return &testFactory{opts: opts}
}
