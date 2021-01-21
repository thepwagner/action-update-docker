package docker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update-docker/docker"
	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/updatertest"
)

var alpine3110 = updater.Dependency{Path: "alpine", Version: "3.11.0"}

func TestUpdater_CheckDockerhub(t *testing.T) {
	t.Skip("early integration test, check https://hub.docker.com/_/alpine?tab=tags")
	u := updatertest.CheckInFixture(t, "simple", updaterFactory(), alpine3110, nil)
	assert.NotNil(t, u)
	assert.Equal(t, "3.12.0", u.Next)
}

func TestUpdater_CheckAuth(t *testing.T) {
	t.Skip("early integration test")
	var privateImage = updater.Dependency{Path: "ghcr.io/thepwagner/alpine", Version: "3.11.0"}
	u := updatertest.CheckInFixture(t, "simple", updaterFactory(), privateImage, nil)
	assert.NotNil(t, u)
	assert.Equal(t, "3.12.0", u.Next)
}

func TestUpdater_CheckAuthPinned(t *testing.T) {
	t.Skip("early integration test")
	var privateImage = updater.Dependency{Path: "ghcr.io/thepwagner/alpine", Version: "sha256:d371657a4f661a854ff050898003f4cb6c7f36d968a943c1d5cde0952bd93c80"}
	u := updatertest.CheckInFixture(t, "pinned", updaterFactory(), privateImage, nil)
	assert.NotNil(t, u)
	assert.Equal(t, "3.12.0", u.Next)
}

func TestUpdater_Check(t *testing.T) {
	cases := map[string]struct {
		dep      updater.Dependency
		tags     []string
		expected string
	}{
		"proposes update": {
			dep:      alpine3110,
			tags:     []string{"3.10.0", "3.11.0", "3.11.1"},
			expected: "3.11.1",
		},
		"no update available": {
			dep:      alpine3110,
			tags:     []string{"3.10.0", "3.11.0"},
			expected: "",
		},
		"prefer longest": {
			dep:      alpine3110,
			tags:     []string{"3.12", "3.12.0", "3.11.0", "3.11.1"},
			expected: "3.12.0",
		},
		"prefer longest 2": {
			dep:      alpine3110,
			tags:     []string{"3.12.0", "3.12", "3.11.0", "3.11.1"},
			expected: "3.12.0",
		},
		"maintains suffix": {
			dep:      updater.Dependency{Path: "redis", Version: "6.0.0-alpine"},
			tags:     []string{"6.0.0-alpine", "6.0.8-alpine", "6.0.9-alpine3.12"},
			expected: "6.0.8-alpine",
		},
		"maintains suffix 2": {
			dep:      updater.Dependency{Path: "redis", Version: "6.0.0-alpine3.12"},
			tags:     []string{"6.0.0-alpine", "6.0.9-alpine", "6.0.8-alpine3.12"},
			expected: "6.0.8-alpine3.12",
		},
		"maintains suffix 3": {
			dep:      updater.Dependency{Path: "redis", Version: "6.0.0"},
			tags:     []string{"6.0.8", "6.0.9-alpine", "6.0.9-alpine3.12"},
			expected: "6.0.8",
		},
	}

	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			u := newUpdaterMockTags(tc.tags...)
			update, err := u.Check(context.Background(), tc.dep, nil)
			require.NoError(t, err)

			if tc.expected != "" {
				if assert.NotNil(t, update) {
					assert.Equal(t, tc.expected, update.Next)
				}
			} else {
				assert.Nil(t, update)
			}
		})
	}
}

func newUpdaterMockTags(tags ...string) *docker.Updater {
	tl := &mockTagLister{}
	tl.On("Tags", mock.Anything, mock.Anything).Return(tags, nil)
	return docker.NewUpdater("", docker.WithTagsLister(tl))
}
