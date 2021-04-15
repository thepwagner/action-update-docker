module github.com/thepwagner/action-update-docker

go 1.15

require (
	github.com/cloudflare/cfssl v1.5.0 // indirect
	github.com/docker/cli v20.10.0-beta1.0.20201029214301-1d20b15adc38+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/lib/pq v1.10.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/moby/buildkit v0.8.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thepwagner/action-update v0.0.40
	github.com/theupdateframework/notary v0.7.0 // indirect
	golang.org/x/mod v0.4.2

)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
)
