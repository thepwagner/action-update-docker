package docker_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/updatertest"
)

var (
	alpine3120 = updater.Update{Path: "alpine", Previous: "3.11.0", Next: "3.12.0"}
	redis608   = updater.Update{Path: "redis", Previous: "6.0.0-alpine", Next: "6.0.8-alpine"}
)

func TestUpdater_ApplyUpdate_Simple(t *testing.T) {
	dockerfile := applyUpdateToFixture(t, "simple", alpine3120)
	assert.Contains(t, dockerfile, "alpine:3.12.0")
	assert.NotContains(t, dockerfile, "alpine:3.11.0")
}

func TestUpdater_ApplyUpdate_Simple_Pinned(t *testing.T) {
	dockerfile := applyUpdateToFixture(t, "simple", alpine3120, withShaPinning("sha256:pinned")...)
	assert.Contains(t, dockerfile, "alpine@sha256:pinned")
	assert.NotContains(t, dockerfile, "alpine:3.12.0")
	assert.NotContains(t, dockerfile, "alpine:3.11.0")
}

func TestUpdater_ApplyUpdate_BuildArg(t *testing.T) {
	dockerfile := applyUpdateToFixture(t, "buildarg", alpine3120)
	assert.Contains(t, dockerfile, "FROM alpine:$ALPINE_VERSION")
	assert.Contains(t, dockerfile, "ARG ALPINE_VERSION=3.12.0")
	assert.NotContains(t, dockerfile, "ARG ALPINE_VERSION=3.11.0")
}

func TestUpdater_ApplyUpdate_BuildArgInterpolate(t *testing.T) {
	dockerfile := applyUpdateToFixture(t, "buildarg", redis608)
	assert.Contains(t, dockerfile, "FROM redis:${REDIS_VERSION}-alpine")
	assert.Contains(t, dockerfile, "FROM redis:$REDIS_VERSION-alpine")
	assert.Contains(t, dockerfile, "ARG REDIS_VERSION=6.0.8")
	assert.NotContains(t, dockerfile, "ARG REDIS_VERSION=6.0.0")
}

func TestUpdater_ApplyUpdate_Comments(t *testing.T) {
	dockerfile := applyUpdateToFixture(t, "comments", redis608)
	assert.Contains(t, dockerfile, "ARG REDIS_VERSION=6.0.8 # redis")
	assert.Contains(t, dockerfile, "# check out this whitespace\n\n\n# intentional trailing spaces  \n")
}

func applyUpdateToFixture(t *testing.T, fixture string, update updater.Update, opts ...docker.UpdaterOpt) string {
	tempDir := updatertest.ApplyUpdateToFixture(t, fixture, updaterFactory(opts...), update)
	b, err := ioutil.ReadFile(filepath.Join(tempDir, "Dockerfile"))
	require.NoError(t, err)
	dockerfile := string(b)
	t.Log(dockerfile)
	return dockerfile
}

func withShaPinning(pinned string) []docker.UpdaterOpt {
	mockPinner := &mockImagePinner{}
	mockPinner.On("Pin", mock.Anything, mock.Anything).Return(pinned, nil)
	return []docker.UpdaterOpt{
		docker.WithShaPinning(true),
		docker.WithImagePinner(mockPinner),
	}
}
