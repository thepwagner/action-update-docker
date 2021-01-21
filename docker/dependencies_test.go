package docker_test

import (
	"testing"

	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/updatertest"
)

func TestUpdater_Dependencies(t *testing.T) {
	cases := map[string][]updater.Dependency{
		"simple": {
			{Path: "alpine", Version: "3.11.0"},
		},
		"buildarg": {
			{Path: "redis", Version: "6.0.0-alpine"},
			{Path: "redis", Version: "6.0.0-alpine"},
			{Path: "alpine", Version: "3.11.0"},
		},
		"comments": {
			{Path: "redis", Version: "6.0.0-alpine"},
			{Path: "alpine", Version: "3.11.0"},
		},
		"pinned": {
			{Path: "alpine", Version: "sha256:7c92a2c6bbcb6b6beff92d0a940779769c2477b807c202954c537e2e0deb9bed"},
		},
	}
	updatertest.DependenciesFixtures(t, updaterFactory(), cases)
}
