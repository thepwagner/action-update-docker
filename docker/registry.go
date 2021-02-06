package docker

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	manifest "github.com/docker/cli/cli/manifest/types"
	"github.com/docker/cli/cli/registry/client"
	"github.com/docker/cli/cli/trust"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/sirupsen/logrus"
	notary "github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/tuf/data"
)

type TagLister interface {
	// Tags returns potential version tags given a updater.Dependency path
	Tags(ctx context.Context, path string) ([]string, error)
}

type ImagePinner interface {
	// Pin normalizes Docker image name to sha256 pinned image.
	Pin(ctx context.Context, image string) (string, error)
	Unpin(ctx context.Context, image, hash string) (string, error)
}

type RemoteRegistries struct {
	rt       http.RoundTripper
	trustKey string
}

func NewRemoteRegistries(trustKey string) *RemoteRegistries {
	return &RemoteRegistries{
		rt:       http.DefaultTransport,
		trustKey: trustKey,
	}
}

const userAgent = "action-update-docker/1.0"

func (r *RemoteRegistries) Tags(ctx context.Context, image string) ([]string, error) {
	// Normalize image name:
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image name: %w", err)
	}
	logrus.WithField("image", normalized.String()).Debug("listing image tags")

	cli, err := r.newDockerCLI()
	if err != nil {
		return nil, err
	}
	resolver := authResolver(cli)

	if !cli.ContentTrustEnabled() {
		logrus.Debug("content trust disabled, listing registry")
		tags, err := client.NewRegistryClient(resolver, userAgent, false).GetTags(ctx, normalized)
		if err != nil {
			return nil, fmt.Errorf("listing tags: %w", err)
		}
		return tags, nil
	}

	logrus.Debug("content trust enabled, listing notary")
	targets, err := r.notaryListTargets(ctx, normalized.Name(), resolver, cli)
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0, len(targets))
	for _, targetWithRole := range targets {
		tags = append(tags, targetWithRole.Name)
	}
	return tags, nil
}

func authResolver(cli *command.DockerCli) func(ctx context.Context, index *registry.IndexInfo) types.AuthConfig {
	resolver := func(ctx context.Context, index *registry.IndexInfo) types.AuthConfig {
		return command.ResolveAuthConfig(ctx, cli, index)
	}
	return resolver
}

func (r *RemoteRegistries) newDockerCLI() (*command.DockerCli, error) {
	cli, err := command.NewDockerCli(command.WithContentTrust(r.trustKey != ""))
	if err != nil {
		return nil, err
	}
	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return nil, fmt.Errorf("initializing cli: %w", err)
	}
	return cli, nil
}

func (r *RemoteRegistries) verifyRootTrust(notaryRepo notary.Repository) error {
	roles, err := notaryRepo.ListRoles()
	if err != nil {
		return fmt.Errorf("listing roles: %w", err)
	}
	for _, role := range roles {
		if role.Name == data.CanonicalRootRole {
			for _, keyID := range role.KeyIDs {
				if keyID == r.trustKey {
					for _, sig := range role.Signatures {
						if sig.KeyID == keyID && sig.IsValid {
							return nil
						}
					}
				}
			}
		}
	}
	return fmt.Errorf("trusted root key not found")
}

func (r *RemoteRegistries) Pin(ctx context.Context, image string) (string, error) {
	// Normalize image name:
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", fmt.Errorf("invalid image name: %w", err)
	}
	logrus.WithField("image", normalized.String()).Debug("listing image tags")

	cli, err := r.newDockerCLI()
	if err != nil {
		return "", err
	}
	resolver := authResolver(cli)

	if !cli.ContentTrustEnabled() {
		registryClient := client.NewRegistryClient(resolver, userAgent, false)
		mf, err := r.getManifest(ctx, registryClient, normalized)
		if err != nil {
			return "", fmt.Errorf("getting manifest: %w", err)
		}
		return mf.Descriptor.Digest.String(), nil
	}

	logrus.Debug("content trust enabled, listing notary")
	targets, err := r.notaryListTargets(ctx, image, resolver, cli)
	if err != nil {
		return "", err
	}

	name := normalized.(reference.Tagged).Tag()
	for _, targetWithRole := range targets {
		if targetWithRole.Name != name {
			continue
		}
		sha256Hash, ok := targetWithRole.Hashes["sha256"]
		if !ok {
			return "", fmt.Errorf("hash not found")
		}
		return fmt.Sprintf("sha256:%x", sha256Hash), nil
	}
	return "", fmt.Errorf("image not found in content trust")
}

func (r *RemoteRegistries) notaryListTargets(ctx context.Context, image string, resolver func(ctx context.Context, index *registry.IndexInfo) types.AuthConfig, cli *command.DockerCli) ([]*notary.TargetWithRole, error) {
	imgRefAndAuth, err := trust.GetImageReferencesAndAuth(ctx, nil, resolver, image)
	if err != nil {
		return nil, err
	}
	notaryRepo, err := trust.GetNotaryRepository(cli.In(), cli.Out(), userAgent, imgRefAndAuth.RepoInfo(), imgRefAndAuth.AuthConfig(), trust.ActionsPullOnly...)
	if err != nil {
		return nil, err
	}

	if err := r.verifyRootTrust(notaryRepo); err != nil {
		return nil, err
	}

	return notaryRepo.ListTargets(data.CanonicalTargetsRole)
}

func (r *RemoteRegistries) Unpin(ctx context.Context, image, hash string) (string, error) {
	normalized, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s@%s", image, hash))
	if err != nil {
		return "", fmt.Errorf("invalid image name: %w", err)
	}

	cli, err := r.newDockerCLI()
	if err != nil {
		return "", err
	}
	resolver := authResolver(cli)

	if !cli.ContentTrustEnabled() {
		registryClient := client.NewRegistryClient(resolver, userAgent, false)
		return r.unpinFromRegistry(ctx, registryClient, normalized, hash)
	}

	logrus.Debug("content trust enabled, listing notary")
	targets, err := r.notaryListTargets(ctx, image, resolver, cli)
	if err != nil {
		return "", err
	}

	for _, targetWithRole := range targets {
		sha256Hash, ok := targetWithRole.Hashes["sha256"]
		if !ok {
			continue
		}
		if hex.EncodeToString(sha256Hash) == hash[7:] {
			return targetWithRole.Name, nil
		}
	}

	return "", fmt.Errorf("TODO: content trust")
}

func (r *RemoteRegistries) unpinFromRegistry(ctx context.Context, registryClient client.RegistryClient, normalized reference.Named, hash string) (string, error) {
	tags, err := registryClient.GetTags(ctx, normalized)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", fmt.Errorf("tag not found")
	}

	// Filter semver tags, work backwards (assuming the pinned sha is a near-latest version)
	semverTags := make([]string, 0)
	for _, tag := range tags {
		if semverIsh(tag) == "" {
			continue
		}
		semverTags = append(semverTags, tag)
	}
	semverTags = semverSort(semverTags)

	logrus.WithFields(logrus.Fields{
		"image": normalized.String(),
		"hash":  hash,
		"tags":  len(semverTags),
	}).Info("listing tags to identify SHA")

	for _, tag := range semverTags {
		normalizedTag, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s:%s", normalized.Name(), tag))
		if err != nil {
			continue
		}
		mf, err := r.getManifest(ctx, registryClient, normalizedTag)
		if err != nil {
			continue
		}
		digest := mf.Descriptor.Digest.String()
		logrus.WithFields(logrus.Fields{
			"tag":    tag,
			"digest": digest,
		}).Debug("fetched image details")

		if digest == hash {
			logrus.WithFields(logrus.Fields{
				"digest": digest,
				"tag":    tag,
			}).Info("resolved pinned image to tag")
			return tag, nil
		}
	}
	return "", fmt.Errorf("manifest not found")
}

func (r *RemoteRegistries) getManifest(ctx context.Context, registryClient client.RegistryClient, normalized reference.Named) (*manifest.ImageManifest, error) {
	// Assume this image is available for one platform:
	mf, err := registryClient.GetManifest(ctx, normalized)
	if err == nil {
		return &mf, nil
	}
	if !strings.Contains(err.Error(), "is a manifest list") {
		return nil, fmt.Errorf("getting manifest: %w", err)
	}

	// Multi-platform images have a list of manifests, select the "right" one:
	manifestList, err := registryClient.GetManifestList(ctx, normalized)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest list: %w", err)
	}
	for _, mf := range manifestList {
		pl := mf.Descriptor.Platform
		if pl.Architecture != "amd64" && pl.OS != "linux" {
			continue
		}
		return &mf, nil
	}
	return nil, fmt.Errorf("could not resolve %q", normalized.String())
}
