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

### Versioning

The app version is embedded into the Go binary at build time via `-ldflags` and surfaced in the **About** dialog (user menu → About) and the `/api/health` endpoint.

**Default:** `0.0.0-dev` (local builds without tags)

**CI automatically derives the version on every build** using `git describe --tags --always --dirty`:
- Push to `develop` → version reflects commits since last tag (e.g. `v1.0.0-5-gabc123`)
- Push tag `v1.2.3` → version is the clean tag (e.g. `v1.2.3`)

**Bumping the version:**

```bash
# Create an annotated semver tag (recommended)
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

This produces versions like:
- `v1.0.0` — clean tagged release
- `v1.0.0-3-gabc123` — 3 commits after `v1.0.0`
- `v1.0.0-3-gabc123-dirty` — uncommitted changes on top

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** — breaking changes (incompatible API or data migrations)
- **MINOR** — new features, backwards compatible
- **PATCH** — bug fixes, backwards compatible
