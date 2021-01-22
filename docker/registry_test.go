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

func TestRemoteRegistries_Tags_Trust(t *testing.T) {
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
	cases := []struct {
		image    string
		key      string
		expected string
	}{
		{
			image:    "alpine:3.11.0",
			key:      "a2489bcac7a79aa67b19b96c4a3bf0c675ffdf00c6d2fabe1a5df1115e80adce",
			expected: "sha256:7c92a2c6bbcb6b6beff92d0a940779769c2477b807c202954c537e2e0deb9bed",
		},
		{
			image:    "alpine:3.11.0",
			expected: "sha256:7c92a2c6bbcb6b6beff92d0a940779769c2477b807c202954c537e2e0deb9bed",
		},
		{
			image:    "ghcr.io/thepwagner-smurf/alpine:3.11.0",
			expected: "sha256:d371657a4f661a854ff050898003f4cb6c7f36d968a943c1d5cde0952bd93c80",
		},
	}

	for _, tc := range cases {
		t.Run(tc.image, func(t *testing.T) {
			reg := docker.NewRemoteRegistries(tc.key)
			pinned, err := reg.Pin(context.Background(), tc.image)
			require.NoError(t, err)
			t.Log(pinned)
			assert.Equal(t, tc.expected, pinned)
		})
	}
}

func TestRemoteRegistries_Unpin(t *testing.T) {
	t.Skip("queries dockerhub")
	cases := []struct {
		image    string
		key      string
		hash     string
		expected string
	}{
		{
			image:    "alpine",
			hash:     "sha256:7c92a2c6bbcb6b6beff92d0a940779769c2477b807c202954c537e2e0deb9bed",
			expected: "3.11.0",
		},
		{
			image:    "alpine",
			key:      "a2489bcac7a79aa67b19b96c4a3bf0c675ffdf00c6d2fabe1a5df1115e80adce",
			hash:     "sha256:7c92a2c6bbcb6b6beff92d0a940779769c2477b807c202954c537e2e0deb9bed",
			expected: "3.11.0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.image, func(t *testing.T) {
			reg := docker.NewRemoteRegistries(tc.key)
			unpinned, err := reg.Unpin(context.Background(), tc.image, tc.hash)
			require.NoError(t, err)
			t.Log(unpinned)
			assert.Equal(t, tc.expected, unpinned)
		})
	}
}
