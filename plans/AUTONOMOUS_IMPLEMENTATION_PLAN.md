# Photobox — Autonomous Implementation Plan

## Overview

This document defines how Hermes (the AI agent) autonomously implements Photobox issues, validates each feature, and integrates with the existing deployment pipeline.

**Goal**: Rich picks an issue → Hermes implements it end-to-end → validates it works → opens a PR → Rich reviews and merges.

---

## 1. Infrastructure & Access

### What Hermes Has
| Resource | Access | Purpose |
|----------|--------|---------|
| GitHub API | `r1chjames/photobox` (full repo access) | Create branches, push code, open PRs, update issues |
| GitLab API | `r1chjames/kubernetes-helm-charts` | Update Helm chart when deployment changes needed |
| Terminal | Full Linux environment | Build, test, run docker-compose, run CI locally |
| Docker | Available (or installable) | Run local dev environment for validation |
| Obsidian notes | `/opt/data/obsidian-code/` | Reference project docs and decisions |

### What Hermes Needs From Rich
| Resource | Status | Notes |
|----------|--------|-------|
| K3s cluster access | ✅ Available | K3s credentials are in scope — Layer 4 validation can run against the real cluster |
| ArgoCD sync trigger | ❓ Unknown | How to trigger a sync after Helm chart changes |
| Test photo library | ❓ Unknown | Sample photos for integration testing |
| Approval workflow | ✅ Confirmed | See §5 below |

### Action Items for Rich
1. **ArgoCD sync**: Does pushing to `main` on `kubernetes-helm-charts` auto-sync via ArgoCD, or is a manual sync needed?
2. **Test data**: Should Hermes generate test photos, or is there an existing test library?

---

## 2. Validation Framework

Every issue goes through a **4-layer validation pyramid**. Not all layers apply to every issue — the plan specifies which layers are required.

```
┌─────────────────────────────────────────────────┐
│  Layer 4: K3s Deployment Validation             │  ← Helm/ArgoCD
│  (Production-like environment)                  │
├─────────────────────────────────────────────────┤
│  Layer 3: Docker Compose E2E Validation         │  ← Full stack locally
│  (API + Webapp + Postgres + Ollama)             │
├─────────────────────────────────────────────────┤
│  Layer 2: Integration Tests                     │  ← Real DB, real filesystem
│  (Go integration tests + Playwright E2E)        │
├─────────────────────────────────────────────────┤
│  Layer 1: Unit Tests + Lint                     │  ← Fast feedback
│  (Go test -race, Vitest, golangci-lint)         │
└─────────────────────────────────────────────────┘
```

### Layer 1: Unit Tests + Lint (Required for ALL issues)
**When**: Every commit, every PR
**How**:
```bash
# Go backend
cd api && golangci-lint run ./...
cd api && go test -v -race -count=1 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out  # Must be ≥70%

# React frontend
cd webapp && npm ci && npm run lint
cd webapp && npm test -- --coverage  # Must be ≥60%
cd webapp && npm run build  # Must compile without errors
```

**Pass criteria**:
- All existing tests still pass (no regressions)
- New code has tests covering the new functionality
- Lint passes with zero errors
- Build succeeds

### Layer 2: Integration Tests (Required for backend/API changes)
**When**: Any issue that touches API handlers, services, or repositories
**How**:
```bash
# Start test dependencies
docker compose -f api/test/docker-compose.test.yml up -d postgres valkey

# Run integration tests
cd api && go test -v -tags=integration ./test/integration/...

# Cleanup
docker compose -f api/test/docker-compose.test.yml down
```

**Pass criteria**:
- All integration tests pass
- API responds correctly to test requests
- Database migrations run cleanly

### Layer 3: Docker Compose E2E (Required for UI changes, new features)
**When**: Any issue that adds/modifies user-facing functionality
**How**:
```bash
# Build and start full stack
cd photobox
docker compose -f docker-compose.local.yml up --build -d

# Wait for health
docker compose -f docker-compose.local.yml ps  # All containers healthy

# Run E2E validation (see §3)
# ... curl commands, Playwright tests, manual verification ...

# Cleanup
docker compose -f docker-compose.local.yml down -v
```

**Pass criteria**:
- All containers start and pass health checks
- Feature works end-to-end (upload → process → view)
- No console errors in browser
- API returns expected responses

### Layer 4: K3s Deployment Validation (Required for infra/deployment changes)
**When**: Any issue that touches Helm chart, Dockerfile, or deployment config
**How**:
```bash
# Option A: Local validation only (if no K3s access)
cd kubernetes-helm-charts/charts/photobox
helm lint .
helm template . --debug  # Must render without errors

# Option B: Full deployment (if K3s access available)
# Push to kubernetes-helm-charts → ArgoCD syncs → verify pods
kubectl -n media get pods -l app=photobox
kubectl -n media logs -l app=photobox --tail=50
curl https://photos.r1chjames.co.uk/api/health
```

**Pass criteria**:
- Helm chart lints cleanly
- Template renders without errors
- (If deployed) Pods start, pass health checks, serve traffic

---

## 3. Issue-Specific Validation Recipes

Each issue type has a specific validation recipe. Below are the patterns for the 24 issues created.

### Short-Term Issues

| Issue | Layer 1 | Layer 2 | Layer 3 | Layer 4 | Validation Method |
|-------|---------|---------|---------|---------|-------------------|
| **#90 UI bugs** | ✅ | — | ✅ | — | Visual: load Settings, Duplicates, PhotoDetail without crash. Dark mode check. |
| **#91 Helm chart** | — | — | — | ✅ | `helm lint` + `helm template` + deploy to K3s |
| **#92 Virtual scroll** | ✅ | — | ✅ | — | Load 10K+ photos, verify 60fps scroll, <500 DOM nodes |
| **#93 Thumbnail prefetch** | ✅ | — | ✅ | — | Network tab: prefetch requests visible, no pop-in |
| **#94 Health checks** | ✅ | ✅ | ✅ | ✅ | `curl /api/health` → 200, `curl /api/ready` → 200/503 |
| **#95 Env validation** | ✅ | — | ✅ | — | Start with missing env → clear error. Start with all env → success. |
| **#96 Test pipeline** | ✅ | ✅ | — | — | CI runs green on PR. Coverage report posted. |
| **#97 Mobile responsive** | ✅ | — | ✅ | — | Chrome DevTools mobile viewport. Touch gestures work. |

### Medium-Term Issues

| Issue | Layer 1 | Layer 2 | Layer 3 | Layer 4 | Validation Method |
|-------|---------|---------|---------|---------|-------------------|
| **#98 Upload** | ✅ | ✅ | ✅ | — | Drag-drop files → progress bars → photos appear in grid |
| **#99 Bulk ops** | ✅ | ✅ | ✅ | — | Select 10 photos → batch add to album → verify album |
| **#100 Smart albums** | ✅ | ✅ | ✅ | — | Create smart album with tag rule → upload matching photo → appears |
| **#101 Search** | ✅ | ✅ | ✅ | — | Combine 3+ filters → correct results in <500ms |
| **#102 Photo edit** | ✅ | ✅ | ✅ | — | Rotate → verify orientation. Crop → verify thumbnail reflects crop. |
| **#103 Video** | ✅ | ✅ | ✅ | — | Upload MP4 → plays in browser with seek bar |
| **#104 EXIF viewer** | ✅ | ✅ | ✅ | — | View photo with EXIF → metadata panel shows camera/settings/GPS |
| **#105 Trash cleanup** | ✅ | ✅ | ✅ | — | Set retention to 1 day → wait → verify auto-deleted |

### Long-Term Issues

| Issue | Layer 1 | Layer 2 | Layer 3 | Layer 4 | Validation Method |
|-------|---------|---------|---------|---------|-------------------|
| **#106 Takeout import** | ✅ | ✅ | ✅ | — | Import test Takeout → verify albums, metadata, dates |
| **#107 Immich migration** | ✅ | ✅ | ✅ | — | Migrate test Immich DB → verify photos/albums preserved |
| **#108 PWA offline** | ✅ | — | ✅ | — | Disable network → cached photos still viewable |
| **#109 Live Photos** | ✅ | ✅ | ✅ | — | Upload HEIC+MOV pair → "Live" badge → press-and-hold plays video |
| **#110 Memories** | ✅ | ✅ | ✅ | — | Photos from 2023 → "On This Day" memory appears |
| **#111 Map view** | ✅ | — | ✅ | — | 100+ GPS photos → clusters form → zoom breaks apart |
| **#112 API docs** | ✅ | — | ✅ | — | `/api/docs` renders Swagger UI. API key auth works. |
| **#113 Plugin arch** | ✅ | ✅ | ✅ | — | Load example plugin → hook fires → metadata added |

---

## 4. Implementation Workflow

### Step-by-Step Process

```
1. Rich selects issue(s) to implement
   ↓
2. Hermes creates feature branch
   git checkout -b feat/issue-NNN-description
   ↓
3. Hermes implements the feature
   - Write code (Go + React)
   - Write tests (unit + integration)
   - Update docs/comments
   ↓
4. Hermes runs Layer 1 validation
   - go test -race ./...
   - npm test -- --coverage
   - golangci-lint run
   - npm run build
   ↓
5. Hermes runs Layer 2 validation (if applicable)
   - docker compose up test deps
   - go test -tags=integration ./test/integration/...
   ↓
6. Hermes runs Layer 3 validation (if applicable)
   - docker compose -f docker-compose.local.yml up --build
   - Manual/automated E2E checks
   - Screenshot or video of working feature
   ↓
7. Hermes runs Layer 4 validation (if applicable)
   - helm lint + helm template
   - Push to kubernetes-helm-charts if needed
   ↓
8. Hermes opens PR against `develop`
   - Links to issue
   - Describes changes
   - Includes validation evidence (test output, screenshots)
   ↓
9. Hermes merges PR into `develop`
   - Rich may review / request changes first if a checkpoint is wanted
   - Direct pushes to `develop` are blocked by the `Protect develop` ruleset
   ↓
10. Hermes updates issue status
    - Closes issue or moves to "Done"
    - Updates Obsidian notes if needed
```

### Branch Flow (enforced by repo rulesets)

```
feature/*  ──PR──▶  develop  ──(Rich)──▶  main
 (agent)           (agent merges)      (Rich, batched stable releases)
```

- **feature → develop**: Agent-owned. All changes land via pull requests; the `Protect develop` ruleset blocks direct pushes to `develop`.
- **develop → main**: Rich-owned. The `Protect main` ruleset grants Rich the sole bypass — only Rich can update `main`.
- **No agent ever pushes to `main`** — enforced at the GitHub ruleset level, not by convention.
- `develop → main` merges happen in batches when a stable release is wanted.

### Branch Naming Convention
```
feat/issue-90-fix-ui-bugs
feat/issue-92-virtual-scroll
feat/issue-98-drag-drop-upload
fix/issue-87-settings-job-state
```

### Commit Message Format
```
feat(issue-92): add virtual scrolling to photo grid

- Implement @tanstack/react-virtual for row-based virtualization
- Add IntersectionObserver for prefetch zone
- Maintain keyboard navigation (j/k, arrows, page up/down)
- Timeline scrubber integrates with scroll position

Closes #92
```

---

## 5. Approval Workflow

### Option A: Full Autonomy (Rich trusts Hermes completely)
- Hermes implements, validates, opens PR
- Rich reviews and merges at convenience
- Hermes can merge their own PRs if Rich grants permission

### Option B: Checkpoint Approval (Confirmed)
- Hermes implements and validates locally
- Hermes opens PR to `develop` with validation evidence
- Hermes merges their own PR into `develop` after validation passes
- Rich reviews at the release boundary and merges `develop` → `main` in batches

### Option C: Issue-by-Issue Approval
- Rich selects which issues to implement each session
- Hermes implements selected issues only
- Rich reviews each PR before merge

**Recommendation**: **Option B is confirmed.** Hermes owns `feature → develop` (via PRs); Rich owns `develop → main` and reviews at release time. See the Branch Flow in §4.

---

## 6. Validation Evidence

For each PR, Hermes provides **concrete evidence** that the feature works:

### Backend Changes
```bash
# Test output
$ go test -v -race ./...
PASS
ok  	gitlab.com/r1chjames/photobox/api/internal/core/service	0.234s

# Coverage report
$ go tool cover -func=coverage.out
total:	(statements)	78.3%
```

### Frontend Changes
```bash
# Test output
$ npm test -- --coverage
PASS  src/Components/PhotoGrid/PhotoGrid.test.tsx
  ✓ renders photo grid (23 ms)
  ✓ handles virtual scroll (45 ms)

Coverage:
File              | % Stmts | % Branch | % Funcs | % Lines
PhotoGrid.tsx     |   92.3  |   85.7   |  100.0  |   91.8
```

### E2E Changes
```bash
# Docker compose validation
$ docker compose -f docker-compose.local.yml ps
NAME                STATUS
photobox-api        Up (healthy)
photobox-webapp     Up (healthy)
photobox-postgres   Up (healthy)

# API validation
$ curl -s http://localhost:8080/api/health | jq
{
  "status": "ok",
  "version": "1.2.3"
}

# UI validation (screenshot or video)
# [Attach screenshot showing the feature working]
```

### Deployment Changes
```bash
# Helm validation
$ helm lint charts/photobox
==> Linting charts/photobox
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed

$ helm template charts/photobox --debug
# [Verify rendered YAML is correct]
```

---

## 7. Rollback & Safety

### What Could Go Wrong
| Risk | Mitigation |
|------|-----------|
| Hermes breaks existing functionality | Layer 1 tests catch regressions. PR must pass all existing tests. |
| Hermes introduces security vulnerability | golangci-lint catches common issues. Rich reviews code before merge. |
| Hermes deploys broken Helm chart | `helm lint` + `helm template` catch rendering errors. ArgoCD sync is manual. |
| Hermes wastes time on wrong approach | Rich approves issue selection. Hermes validates locally before PR. |

### Rollback Plan
- **Code rollback**: Git revert. Every PR is a discrete unit.
- **Deployment rollback**: ArgoCD `app rollback photobox <prev-id>`
- **Data rollback**: PostgreSQL backups (CNPG). Photos in S3 are immutable.

---

## 8. Progress Tracking

### Issue Status Updates
Hermes updates issue status as work progresses:

| Status | Meaning |
|--------|---------|
| `Open` | Not started |
| `In Progress` | Hermes is working on it (branch exists) |
| `PR Open` | PR submitted, awaiting review |
| `In Review` | Rich is reviewing |
| `Changes Requested` | Rich requested changes, Hermes is addressing |
| `Approved` | Rich approved, ready to merge |
| `Merged` | PR merged, issue closed |
| `Deployed` | Merged and deployed to K3s |

### Session Reports
At the end of each autonomous session, Hermes reports:
```
## Session Report — 2026-07-31

### Completed
- ✅ #90: Fix UI bugs — PR #123 merged
- ✅ #94: Health check endpoints — PR #124 merged

### In Progress
- 🔄 #92: Virtual scrolling — PR #125 open, awaiting review

### Blocked
- ⏸️ #91: Helm chart — needs K3s access for validation

### Next Session Plan
- Address review feedback on #125
- Start #93: Thumbnail prefetch
- Start #95: Env validation
```

---

## 9. Getting Started

### First Session Checklist
1. ✅ Rich confirms access (GitHub, GitLab, K3s if needed)
2. ✅ Rich selects first 2-3 issues to implement
3. ✅ Hermes sets up local dev environment (docker-compose)
4. ✅ Hermes implements first issue
5. ✅ Hermes validates (Layers 1-3)
6. ✅ Hermes opens PR against `develop`
7. ✅ Hermes merges PR into `develop`
8. ✅ Rich reviews at release time and merges `develop` → `main`
9. ✅ Hermes moves to next issue

### Recommended First Issues (Quick Wins)
| Issue | Why Start Here |
|-------|---------------|
| **#94 Health checks** | Small, well-defined, validates deployment pipeline |
| **#95 Env validation** | Small, improves reliability, no UI changes |
| **#90 UI bugs** | Fixes known issues, improves daily usability |
| **#96 Test pipeline** | Enables all future validation |

---

## 10. Open Questions for Rich

1. **ArgoCD sync**: Does pushing to `kubernetes-helm-charts` auto-sync, or manual?
2. **Test data**: Should Hermes generate test photos, or use existing library?
3. **First issues**: Which 2-3 issues should Hermes tackle first?
4. **Session length**: How long should each autonomous session run? (1 hour? 4 hours? Until done?)
5. **Communication**: Should Hermes report progress mid-session, or only at the end?

---

## Appendix: Validation Commands Quick Reference

### Backend
```bash
# Unit tests
cd api && go test -v -race -count=1 ./...

# Integration tests
cd api && go test -v -tags=integration ./test/integration/...

# Lint
cd api && golangci-lint run ./...

# Coverage
cd api && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# Build
cd api && go build -o /tmp/photobox-api ./main.go
```

### Frontend
```bash
# Install
cd webapp && npm ci

# Unit tests
cd webapp && npm test -- --coverage

# Lint
cd webapp && npm run lint

# Build
cd webapp && npm run build

# Storybook (visual testing)
cd webapp && npm run storybook
```

### Full Stack
```bash
# Start
docker compose -f docker-compose.local.yml up --build -d

# Health check
docker compose -f docker-compose.local.yml ps

# API health
curl -s http://localhost:8080/api/health | jq

# Webapp
curl -s http://localhost:80 | head -20

# Logs
docker compose -f docker-compose.local.yml logs -f api

# Stop
docker compose -f docker-compose.local.yml down -v
```

### Helm
```bash
# Lint
cd kubernetes-helm-charts/charts/photobox && helm lint .

# Template
cd kubernetes-helm-charts/charts/photobox && helm template . --debug

# Deploy (if K3s access)
cd kubernetes-helm-charts && git push origin main
# Wait for ArgoCD sync, then:
kubectl -n media get pods -l app=photobox
```
