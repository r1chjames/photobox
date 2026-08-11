# Implementation Decisions & Requirements

## Status: Ready to proceed

Access is confirmed for K3s and the branch flow is enforced by repo rulesets. Remaining decisions are listed in §6.

---

## 1. K3s Cluster Access

### Current State
✅ **kubectl installed**: Available at `/opt/data/kubectl`  
✅ **kubeconfig found**: `/opt/data/home/.kube/config`  
✅ **API server reachable**: K3s v1.35.6  
✅ **Credentials in scope**: K3s tokens are available to Hermes — cluster access is no longer a blocker

### What I Can Do
- Full Layer 4 validation: `helm lint`, `helm template`, and `kubectl` against the live cluster
- Verify deployments, logs, events, and health endpoints (`/api/health`, `/api/ready`)
- Trigger deployments/rollbacks as needed (with Rich's approval for production-affecting changes)

---

## 2. ArgoCD Sync Trigger

### Current State
ArgoCD sync is manual. Need to determine if we can trigger it via API.

### ArgoCD API for Sync

ArgoCD provides a REST API to trigger application syncs:

```bash
# Trigger sync for photobox application
curl -X POST https://<argocd-server>/api/v1/applications/photobox/sync \
  -H "Authorization: Bearer <argocd-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "prune": true,
    "dryRun": false,
    "strategy": {
      "hook": {
        "force": true
      }
    }
  }'
```

### Required Information
To use this, I need:
1. **ArgoCD server URL**: e.g., `https://argocd.r1chjames.co.uk` or internal URL
2. **ArgoCD API token**: A service account token with `sync` permissions for the `photobox` application

### How to Create ArgoCD Token
```bash
# Login to ArgoCD CLI
argocd login <argocd-server>

# Create a service account token
argocd account generate-token --account hermes-ci --id hermes-photobox-sync

# Or create via UI:
# Settings → Accounts → Create Application Account → Generate Token
```

### Alternative: Git-Based Sync
If ArgoCD is configured with automated sync on git push (even if currently disabled), we could:
1. Push changes to `kubernetes-helm-charts` repo
2. Rich manually clicks "Sync" in ArgoCD UI
3. Or enable auto-sync for photobox app specifically

**Recommendation**: Provide ArgoCD server URL + API token so I can trigger syncs programmatically after Helm chart changes.

---

## 3. Test Data Strategy

### Current State
No existing test photo library.

### Proposed Strategy

#### Option A: Generate Synthetic Test Photos (Recommended)
Create a script that generates test photos with realistic EXIF data:

```bash
# Script: generate-test-photos.sh
# Generates 100 test photos with varied:
# - File formats (JPEG, HEIC, PNG, WebP)
# - EXIF metadata (camera models, dates, GPS coordinates)
# - File sizes (1MB - 20MB)
# - Resolutions (1080p to 24MP)
# - Orientations (portrait, landscape, square)
```

**Pros**:
- Reproducible
- No privacy concerns
- Can test edge cases (corrupt EXIF, huge files, etc.)
- Can generate on-demand

**Cons**:
- Not "real" photos
- May not catch real-world edge cases

#### Option B: Use Public Domain Photo Sets
Download curated test photo collections:
- [Sample Photos](https://www.samplephotos.com/) - Various cameras, EXIF-rich
- [Unsplash](https://unsplash.com/) - High-quality, free to use
- [Pexels](https://www.pexels.com/) - Similar to Unsplash

**Pros**:
- Real photos with real EXIF
- Variety of cameras and conditions

**Cons**:
- Download time
- Storage space
- May need to curate for specific test cases

#### Option C: Hybrid Approach (Best)
1. **Base set**: 50 synthetic photos (fast, reproducible)
2. **Real set**: 20 public domain photos (EXIF variety)
3. **Edge cases**: 10 manually crafted problem photos (corrupt EXIF, huge files, videos)

### Test Data Requirements
For comprehensive testing, we need:

| Category | Count | Purpose |
|----------|-------|---------|
| JPEG photos | 50 | Basic upload, thumbnail generation |
| HEIC photos | 20 | iPhone format, conversion testing |
| RAW photos | 10 | Large files, DNG/CR2/NEF |
| Videos | 10 | MP4, MOV, various codecs |
| Photos with GPS | 30 | Map view testing |
| Photos without EXIF | 10 | Edge case handling |
| Duplicate photos | 5 | De-duplication testing |
| Large photos (>10MB) | 5 | Performance testing |
| Corrupt files | 3 | Error handling |

**Recommendation**: Option C (Hybrid). I'll create a `test-data/` directory with a generation script + a small curated set of real photos.

---

## 4. Approval Workflow & Branch Flow

### Branch Flow (enforced by repo rulesets)

```
feature/*  ──PR──▶  develop  ──(Rich)──▶  main
 (agent)           (agent merges)      (Rich, batched stable releases)
```

- **feature → develop (agent)**: Hermes opens PRs targeting `develop` and merges them after validation. The `Protect develop` ruleset requires pull requests — direct pushes to `develop` are blocked.
- **develop → main (Rich)**: Only Rich can update `main`. The `Protect main` ruleset grants Rich the sole bypass; no agent can push to or merge into `main`.
- `develop → main` merges are batched for stable releases.

**Workflow**:
1. Rich selects issue(s) to implement
2. Hermes implements + validates locally (Layers 1-3)
3. Hermes opens PR against `develop` with validation evidence
4. Hermes merges PR into `develop` (Rich may review first if a checkpoint is wanted)
5. Rich merges `develop` → `main` in release batches
6. Hermes updates issue status

**Hermes Cannot** (enforced by rulesets):
- Push directly to `main` or merge PRs into `main`
- Push directly to `develop` (PR required)
- Deploy to K3s without Rich's approval

**Hermes Can**:
- Create feature branches
- Open and merge PRs into `develop`
- Run local validation (docker-compose, tests)
- Update issue comments with progress

---

## 5. Report Format

### Required: ADR + Test Plan + Test Outcome

For each implemented issue, Hermes will provide a report in the PR description with three sections:

### Section 1: Architecture Decision Record (ADR)

```markdown
## ADR: [Issue #] [Title]

### Status
[Proposed / Accepted / Deprecated / Superseded]

### Context
What is the issue? What forces are at play? What is the current state?

### Decision
What is the change that we're proposing and/or requiring?

### Consequences
What becomes easier or more difficult to do because of this change?

### Implementation Details
- Files changed
- New dependencies
- Configuration changes
- Migration steps (if any)
```

### Section 2: Test Plan

```markdown
## Test Plan

### Validation Layers
- [ ] Layer 1: Unit tests + lint
- [ ] Layer 2: Integration tests (if applicable)
- [ ] Layer 3: Docker Compose E2E (if applicable)
- [ ] Layer 4: K3s deployment (if applicable)

### Test Cases

#### Test Case 1: [Name]
**Given**: [Initial state]  
**When**: [Action]  
**Then**: [Expected result]  
**Validation**: [How to verify]

#### Test Case 2: [Name]
...

### Manual Verification Steps
1. [Step 1]
2. [Step 2]
3. [Step 3]
```

### Section 3: Test Outcome

```markdown
## Test Outcome

### Layer 1: Unit Tests + Lint
```bash
$ go test -v -race ./...
[output]

$ golangci-lint run ./...
[output]

$ npm test -- --coverage
[output]
```

**Result**: ✅ PASS / ❌ FAIL

### Layer 2: Integration Tests (if applicable)
```bash
$ docker compose -f api/test/docker-compose.test.yml up -d
$ go test -v -tags=integration ./test/integration/...
[output]
```

**Result**: ✅ PASS / ❌ FAIL / ⏭️ SKIPPED (not applicable)

### Layer 3: Docker Compose E2E (if applicable)
```bash
$ docker compose -f docker-compose.local.yml up --build -d
$ docker compose -f docker-compose.local.yml ps
[output showing all containers healthy]

# API validation
$ curl -s http://localhost:8080/api/health | jq
[output]

# UI validation
[Screenshot or description of manual verification]
```

**Result**: ✅ PASS / ❌ FAIL / ⏭️ SKIPPED (not applicable)

### Layer 4: K3s Deployment (if applicable)
```bash
$ helm lint charts/photobox
[output]

$ helm template charts/photobox --debug
[output snippet]

# If deployed:
$ kubectl -n media get pods -l app=photobox
[output]
```

**Result**: ✅ PASS / ❌ FAIL / ⏭️ SKIPPED (not applicable)

### Summary
- **Total test cases**: X
- **Passed**: Y
- **Failed**: Z
- **Skipped**: W

**Overall**: ✅ READY TO MERGE / ❌ NEEDS WORK
```

---

## 6. Next Steps

### Blocked Until Rich Provides:
1. **ArgoCD credentials**: Server URL + API token (or confirm manual sync workflow)
2. **Test data approval**: Confirm Option C (hybrid approach) or suggest alternative

### Once Unblocked:
1. Rich selects first 2-3 issues to implement
2. Hermes creates feature branch for first issue
3. Hermes implements + validates
4. Hermes opens PR with ADR + Test Plan + Test Outcome
5. Rich reviews + merges
6. Repeat

---

## Questions for Rich

1. **ArgoCD**: Can you provide ArgoCD server URL + API token, or should we stick with manual sync?
2. **Test data**: Approve hybrid approach (synthetic + public domain + edge cases)?
3. **First issues**: Which 2-3 issues should I tackle first? (My recommendation: #94 → #95 → #96)
4. **Session length**: How long should each autonomous session run? (1h? 4h? Until done?)
