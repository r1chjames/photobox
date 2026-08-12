# Photobox Helm Chart

Deploys the Photobox application (Go API + React webapp served from a single
nginx container) to Kubernetes.

## Requirements

- Kubernetes 1.24+
- Helm 3.x

## Installing the chart

```bash
helm dependency update
helm upgrade --install photobox . -n media --create-namespace
```

By default the chart deploys:

- A single `photobox` Deployment (combined API + webapp image, port 80)
- A Service (ClusterIP) on port 80
- A Bitnami PostgreSQL subchart
- Ollama is **disabled by default** — enable with `ollama.enabled=true`

## Required secrets (must be overridden in production)

| Value | Requirement |
|-------|-------------|
| `secrets.token` | Paseto symmetric key, **exactly 32 bytes** |
| `secrets.adminPassword` | ≥ 8 chars, not a known-weak value |
| `secrets.dbPassword` | Must match `database.postgresql.auth.password` |

The defaults are insecure placeholders that only exist so `helm template`
renders. Set them via `--set` or a values file in every real deployment.

## Configuration

### Image & replicas

| Value | Default | Description |
|-------|---------|-------------|
| `image.repository` | `ghcr.io/r1chjames/photobox` | Combined image |
| `image.tag` | `latest` | Image tag |
| `image.pullPolicy` | `IfNotPresent` | Pull policy |
| `replicaCount` | `1` | Number of replicas |

### Application environment (`env`)

Non-secret configuration injected as a ConfigMap. Keys: `DB_PORT`, `DB_USER`,
`DB_NAME`, `PHOTO_DIR`, `API_BASE_PATH`, `AI_ENABLED`, `OLLAMA_MODEL`,
`OLLAMA_HOST`, `DEFAULT_ADMIN_USERNAME`, `CORS_ALLOWED_ORIGINS`,
`TOKEN_DURATION`, `CACHE_ENABLED`, `CACHE_HOST`.

### Secrets (`secrets`)

`dbPassword`, `token`, `adminPassword`, `s3Endpoint`, `s3AccessKey`,
`s3SecretKey`, `s3Bucket`. Set `existingSecret` to use a pre-created Secret
instead of the generated one.

### Service

`service.type` (`ClusterIP`), `service.port` (`80`), `service.loadBalancerIP`.

### Persistence

| Value | Default | Description |
|-------|---------|-------------|
| `persistence.enabled` | `false` | Create a PVC for photos |
| `persistence.storageClass` | `""` | StorageClass (default class when empty) |
| `persistence.size` | `10Gi` | PVC size |
| `persistence.hostPath` | `""` | Use a hostPath instead of a PVC |

When both `enabled=false` and `hostPath=""`, an `emptyDir` is used (photos do
not persist).

### Ollama (`ollama`)

`ollama.enabled` (`false`), image, service, resources, and persistence. When
enabled, set `env.AI_ENABLED="true"` and `env.OLLAMA_HOST` is automatically
defaulted to the in-cluster Ollama service.

### Database (`database`)

`database.enabled` (`true`) deploys the Bitnami PostgreSQL subchart. Set
`database.host` when using an external database with `database.enabled=false`.

### Ingress

`ingress.enabled` (`false`), `ingress.className`, `ingress.annotations`
(Traefik TLS annotations are provided as comments), `ingress.hosts`,
`ingress.tls`. The API's WebSocket endpoint (`/api/ws`) works through Traefik
out of the box; for other controllers add websocket/upgrade annotations.

### Autoscaling

`autoscaling.enabled` (`false`), `minReplicas`, `maxReplicas`,
`targetCPUUtilizationPercentage`.

### Security context

- `podSecurityContext.fsGroup` (default `1000`) — group ownership of the photos volume.
- `securityContext` (default drops all capabilities, disables privilege
  escalation) — full non-root / `readOnlyRootFilesystem` hardening requires a
  non-root image; see `README` in the repository root.

## Health & probes

The chart configures the following probes against the nginx port (80), which
proxies `/api/` to the in-container API:

| Probe | Path | Purpose |
|-------|------|---------|
| `livenessProbe` | `/api/health` | Process is alive |
| `readinessProbe` | `/api/ready` | DB + storage reachable |
| `startupProbe` | `/api/healthz` | Slow initialization (migrations) |
