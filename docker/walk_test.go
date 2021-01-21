package docker_test

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
)

const fixtureCount = 4

func TestWalkDockerfiles(t *testing.T) {
	var cnt int64
	err := docker.WalkDockerfiles("testdata/", nil, func(path string, _ *parser.Result) error {
		atomic.AddInt64(&cnt, 1)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, int64(fixtureCount), cnt, "function not invoked N times")
}

func TestWalkDockerfiles_PathFilter(t *testing.T) {
	filter := func(s string) bool {
		return !strings.HasPrefix(s, "simple")
	}

	var cnt int64
	err := docker.WalkDockerfiles("testdata/", filter, func(path string, _ *parser.Result) error {
		fmt.Println(path)
		atomic.AddInt64(&cnt, 1)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, int64(1), cnt, "function not invoked N times")
}
