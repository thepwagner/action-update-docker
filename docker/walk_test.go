package docker_test

import (
	"sync/atomic"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
)

const fixtureCount = 3

func TestWalkDockerfiles(t *testing.T) {
	var cnt int64
	err := docker.WalkDockerfiles("testdata/", func(path string, _ *parser.Result) error {
		atomic.AddInt64(&cnt, 1)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, int64(fixtureCount), cnt, "function not invoked N times")
}
