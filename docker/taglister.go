package docker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/cli/cli/command"
	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/sirupsen/logrus"
)

type TagLister interface {
	Tags(ctx context.Context, path string) ([]string, error)
}

type RemoteTagLister struct {
	rt http.RoundTripper
}

func NewRemoteTagLister() *RemoteTagLister {
	return &RemoteTagLister{
		rt: http.DefaultTransport,
	}
}

func (r *RemoteTagLister) Tags(ctx context.Context, image string) ([]string, error) {
	// Normalize image name:
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image name: %w", err)
	}
	domain := reference.Domain(normalized)
	if domain == "docker.io" {
		domain = "index.docker.io"
	}
	imagePath := reference.Path(normalized)
	logrus.WithFields(logrus.Fields{
		"domain": domain,
		"path":   imagePath,
	}).Debug("listing image tags")

	repo, err := r.newRepository(domain, imagePath)
	if err != nil {
		return nil, err
	}

	tags, err := repo.Tags(ctx).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	return tags, nil
}

func (r *RemoteTagLister) newRepository(domain string, imagePath string) (distribution.Repository, error) {
	// Build authorizing RoundTripper:
	baseURL := fmt.Sprintf("https://%s", domain)
	cm, err := r.challengeManager(baseURL)
	if err != nil {
		return nil, err
	}

	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}


	// TODO: extract ~/.docker/config.json .auths
	// Like https://github.com/docker/cli/blob/8107a381c181e3dec6757d9ffca801863fb1fc6f/cli/registry/client/client.go#L165-L171 ?
	tokenHandler := auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
		Transport: r.rt,
		Scopes: []auth.Scope{
			auth.RepositoryScope{
				Repository: imagePath,
				Actions:    []string{"pull"},
			},
		},
	})
	authedTransport := transport.NewTransport(r.rt, auth.NewAuthorizer(cm, tokenHandler))

	repoName, _ := reference.WithName(imagePath)
	repo, err := client.NewRepository(repoName, baseURL, authedTransport)
	if err != nil {
		return nil, fmt.Errorf("creating repo: %w", err)
	}
	return repo, nil
}

func (r *RemoteTagLister) challengeManager(baseURL string) (challenge.Manager, error) {
	cm := challenge.NewSimpleManager()

	// Send a ping, store challenge:
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/", baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("creating challenge request: %w", err)
	}
	res, err := r.rt.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("sending challenge request: %w", err)
	}
	if err := cm.AddResponse(res); err != nil {
		return nil, fmt.Errorf("parsing challenge request: %w", err)
	}

	return cm, nil
}
