module github.com/thepwagner/action-update-docker

go 1.15

require (
	github.com/docker/distribution v0.0.0-20200223014041-6b972e50feee
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/moby/buildkit v0.7.2
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/thepwagner/action-update v0.0.11
	golang.org/x/mod v0.3.0
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
