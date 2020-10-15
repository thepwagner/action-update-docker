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
	}

	for _, tc := range cases {
		tl := docker.NewRemoteTagLister()
		t.Run(tc, func(t *testing.T) {
			tags, err := tl.Tags(context.Background(), tc)
			require.NoError(t, err)
			assert.NotEmpty(t, tags)
		})
	}
}
