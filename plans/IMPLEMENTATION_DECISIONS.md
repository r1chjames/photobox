# Implementation Decisions & Requirements

## Status: Pending Rich's Input

This document captures the decisions and requirements needed before we can begin autonomous implementation.

---

## 1. K3s Cluster Access

### Current State
✅ **kubectl installed**: Available at `/opt/data/kubectl`  
✅ **kubeconfig found**: `/opt/data/home/.kube/config`  
✅ **API server reachable**: K3s v1.35.6  
❌ **Insufficient permissions**: Service account `hermes-sa` can only discover API resources, cannot manage pods/deployments

### What I Can Do
- Query API server version and metadata
- List available API resources (but not actual resources)
- Access health/readiness endpoints

### What I Cannot Do
- List/get/create/update/delete pods, deployments, services
- View logs or events
- Trigger deployments or rollbacks

### Required Action
**Option A**: Grant the `hermes-sa` service account additional permissions:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: hermes-deployment-access
  namespace: media  # or whichever namespace photobox runs in
subjects:
- kind: ServiceAccount
  name: hermes-sa
  namespace: personal
roleRef:
  kind: ClusterRole
  name: edit  # or a custom role with deployment permissions
  apiGroup: rbac.authorization.k8s.io
```

**Option B**: Provide a different kubeconfig with broader permissions (e.g., admin kubeconfig)

**Option C**: Limit Layer 4 validation to `helm lint` + `helm template` only (no actual deployment)

**Recommendation**: Option A with a scoped role limited to the `media` namespace and photobox resources only.

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

## 4. Approval Workflow

### Confirmed: Option B (Checkpoint Approval)

**Workflow**:
1. Rich selects issue(s) to implement
2. Hermes implements + validates locally (Layers 1-3)
3. Hermes opens PR with validation evidence
4. Rich reviews PR
5. Rich approves + merges
6. Hermes updates issue status

**Hermes Cannot**:
- Merge their own PRs
- Push directly to `main` or `develop` branches
- Deploy to K3s without Rich's approval

**Hermes Can**:
- Create feature branches
- Open PRs
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
1. **K3s access decision**: Option A (grant permissions), B (new kubeconfig), or C (helm-only validation)?
2. **ArgoCD credentials**: Server URL + API token (or confirm manual sync workflow)
3. **Test data approval**: Confirm Option C (hybrid approach) or suggest alternative

### Once Unblocked:
1. Rich selects first 2-3 issues to implement
2. Hermes creates feature branch for first issue
3. Hermes implements + validates
4. Hermes opens PR with ADR + Test Plan + Test Outcome
5. Rich reviews + merges
6. Repeat

---

## Questions for Rich

1. **K3s access**: Which option (A/B/C) for cluster access?
2. **ArgoCD**: Can you provide ArgoCD server URL + API token, or should we stick with manual sync?
3. **Test data**: Approve hybrid approach (synthetic + public domain + edge cases)?
4. **First issues**: Which 2-3 issues should I tackle first? (My recommendation: #94 → #95 → #96)
5. **Session length**: How long should each autonomous session run? (1h? 4h? Until done?)
