module github.com/thepwagner/action-update-docker

go 1.15

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cloudflare/cfssl v1.4.1 // indirect
	github.com/docker/cli v0.0.0-20200915230204-cd8016b6bcc5
	github.com/docker/distribution v0.0.0-20200223014041-6b972e50feee
	github.com/docker/go v1.5.1-1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/lib/pq v1.8.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.4 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/moby/buildkit v0.7.2
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/thepwagner/action-update v0.0.15
	github.com/theupdateframework/notary v0.6.1 // indirect
	golang.org/x/mod v0.3.0
	gopkg.in/dancannon/gorethink.v3 v3.0.5 // indirect
	gopkg.in/fatih/pool.v2 v2.0.0 // indirect
	gopkg.in/gorethink/gorethink.v3 v3.0.5 // indirect

)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
