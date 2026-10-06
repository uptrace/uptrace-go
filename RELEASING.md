# Releasing new versions

## Prerequisites

- You are on the `master` branch with a clean working tree
- All changes intended for the release are already merged to `master`
- Go 1.26+ is installed (required by OpenTelemetry Go v1.47+)

To upgrade dependencies (e.g. a new OpenTelemetry release), run
`make go_mod_tidy` on a separate branch, verify that `make test` passes and the
examples build, and merge that PR before starting the release.

## Step 1: Run the release script

```bash
TAG=v1.42.0 ./scripts/release.sh
```

This script:

1. Validates the tag format (must be valid semver, e.g. `v1.42.0`)
2. Checks the tag doesn't already exist
3. Ensures the working tree is clean
4. Runs `go mod tidy` on all packages
5. Updates version references in all example `go.mod` files
6. Updates the version string in `uptrace/version.go`
7. Creates a `release/v1.42.0` branch from master
8. Commits and pushes the branch to origin

## Step 2: Create and merge a pull request

Open a pull request from the `release/v1.42.0` branch to `master`. CI will run
tests and linting. Review and merge the PR.

## Step 3: Tag the release

After the PR is merged, pull master and run the tag script:

```bash
git checkout master
git pull origin master
TAG=v1.42.0 ./scripts/tag.sh
```

This script:

1. Verifies that `uptrace/version.go` contains the version
2. Creates and pushes the main tag (`v1.42.0`)
3. Creates and pushes tags for nested Go modules, excluding `example/` and
   `internal/` (currently there are none, so only the main tag is pushed)

## Step 4: Verify the GitHub release

Pushing a `v*` tag triggers the `.github/workflows/release.yml` workflow, which
automatically creates a GitHub Release. Verify it appears at
https://github.com/uptrace/uptrace-go/releases.

## Quick reference

```bash
# Full release flow
TAG=v1.42.0 ./scripts/release.sh
# ... create PR, get it reviewed, merge ...
git checkout master && git pull origin master
TAG=v1.42.0 ./scripts/tag.sh
```
