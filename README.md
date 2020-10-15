# action-update-docker

## This is not endorsed by or associated with GitHub, Dependabot, etc.

This Action checks for available updates to semver-compatible Docker images in Dockerfiles.

* Supports `ARG` interpolation
* All the features common to [action-update](https://github.com/thepwagner/action-update) actions

## Simplest setup

```
- uses: actions/checkout@v2
  # If you use Actions "push" for CI too, a Personal Access Token is required for update PRs to trigger
  with:
    token: ${{ secrets.MY_GITHUB_PAT }}
- uses: actions/setup-go@v2
  with:
    go-version: '1.15.0'
- uses: thepwagner/action-update-docker@main
  with:
    # If you use Actions "pull_request" for CI too, a Personal Access Token is required for update PRs to trigger
    token: ${{ secrets.MY_GITHUB_PAT }}
```

## Private dependencies

If your image requires authentication, you can configure before invoking the action:

```yaml
- uses: actions/checkout@v2
- uses: actions/setup-go@v2
  with:
    go-version: '1.15.0'
- uses: docker/login-action@v1
  with:
    registry: ghcr.io
    username: ${{ github.repository_owner }}
    password: ${{ secrets.MY_GITHUB_PAT }}
- uses: thepwagner/action-update-docker@main
  with:
    token: ${{ secrets.MY_GITHUB_PAT }}
```
