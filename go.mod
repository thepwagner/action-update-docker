module github.com/thepwagner/action-update-docker

go 1.16

require (
	github.com/cloudflare/cfssl v1.6.0 // indirect
	github.com/docker/cli v20.10.8+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/moby/buildkit v0.9.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thepwagner/action-update v0.0.42
	github.com/theupdateframework/notary v0.7.0 // indirect
	golang.org/x/mod v0.4.2

)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
