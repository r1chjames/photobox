## CI / CD

The project builds a single Docker image containing both the Go API and React webapp via GitHub Actions.

### Image tags

| Trigger | Tags |
|---|---|
| Push to `main` | `latest`, `<sha>` |
| Push to `develop` | `develop`, `<sha>` |
| Push tag `v1.2.3` | `1.2.3`, `1.2`, `1`, `<sha>` |
| PR to `main` | Build only (no push) |

Images are published to `ghcr.io/<owner>/photobox`.

### Release flow

1. Feature branches merge into `develop` → CI builds `:develop` (rolling integration image).
2. When ready to release, merge `develop` into `main` → CI builds `:latest`.
3. Tag the merge commit on `main` with a semver tag (`v1.2.3`) → CI builds `:1.2.3`, `:1.2`, `:1`, all pointing to the same image.
