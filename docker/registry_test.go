package docker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
)

func TestRemoteRegistries_Tags(t *testing.T) {
	t.Skip("queries dockerhub")
	cases := []string{
		"alpine",
		"datadog/agent",
		"ghcr.io/thepwagner-smurf/alpine",
	}

	for _, tc := range cases {
		tl := docker.NewRemoteRegistries("")
		t.Run(tc, func(t *testing.T) {
			tags, err := tl.Tags(context.Background(), tc)
			require.NoError(t, err)
			assert.NotEmpty(t, tags)
		})
	}
}

func TestRemoteTagRegistries_TrustTags(t *testing.T) {
	t.Skip("queries notary")
	cases := []struct {
		image string
		key   string
	}{
		{image: "alpine", key: "a2489bcac7a79aa67b19b96c4a3bf0c675ffdf00c6d2fabe1a5df1115e80adce"},
		{image: "debian", key: "575d013f89e3cbbb19e0fb06aa33566c22718318e0c9ffb1ab5cc4291e07bf84"},
		{image: "datadog/agent", key: "5e06443f1750bffcf43423454f3cd06ac13e370e81cced1e83296ab7d62b458b"},
		{image: "jess/vlc", key: "fb4278e3366754a0117b2feaadddfe6aca7b42ecfaacd5a92af3b2a1a1146695"},
	}

	for _, tc := range cases {
		reg := docker.NewRemoteRegistries(tc.key)
		tags, err := reg.Tags(context.Background(), tc.image)
		require.NoError(t, err)
		assert.NotEmpty(t, tags)
	}
}

func TestRemoteTagLister_TrustTags_Untrusted(t *testing.T) {
	rootKey := "a2489bcac7a79aa67b19b96c4a3bf0c675ffdf00c6d2fabe1a5df1115e80adce"
	reg := docker.NewRemoteRegistries(rootKey)

	_, err := reg.Tags(context.Background(), "debian")
	assert.EqualError(t, err, "trusted root key not found")
}

func TestRemoteRegistries_Pin(t *testing.T) {
	t.Skip("queries dockerhub")
	cases := []string{
		"alpine:3.11.0",
		"datadog/agent:7",
		"ghcr.io/thepwagner-smurf/alpine:3.11.0",
	}

	for _, tc := range cases {
		tl := docker.NewRemoteRegistries("")
		t.Run(tc, func(t *testing.T) {
			pinned, err := tl.Pin(context.Background(), tc)
			require.NoError(t, err)
			t.Log(pinned)
			assert.Contains(t, pinned, "@sha256:")
		})
	}
}
