package docker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
)

func TestRemoteTagLister_Tags(t *testing.T) {
	t.Skip("queries dockerhub")
	cases := []string{
		"alpine",
		"datadog/agent",
		"ghcr.io/thepwagner-smurf/alpine",
	}

	for _, tc := range cases {
		tl := docker.NewRemoteRegistries()
		t.Run(tc, func(t *testing.T) {
			tags, err := tl.Tags(context.Background(), tc)
			require.NoError(t, err)
			assert.NotEmpty(t, tags)
		})
	}
}

func TestRemoteTagLister_Pin(t *testing.T) {
	t.Skip("queries dockerhub")
	cases := []string{
		"alpine:3.11.0",
		"datadog/agent:7",
		"ghcr.io/thepwagner-smurf/alpine:3.11.0",
	}

	for _, tc := range cases {
		tl := docker.NewRemoteRegistries()
		t.Run(tc, func(t *testing.T) {
			pinned, err := tl.Pin(context.Background(), tc)
			require.NoError(t, err)
			t.Log(pinned)
			assert.Contains(t, pinned, "@sha256:")
		})
	}
}
