# Cloud Hosting Plan — Options & Cost Model

## Status: v1.1 — decisions locked (2026-08-15); landing page added (§15). Ready for Phase 0 breakdown.

Hosting options for running Photobox as an **open SaaS**, targeting **<$25/mo** at launch with an architecture that **scales without rework**, **no cloud-specific API lock-in**, and **no photo loss** as the hard durability requirement. AI features (Ollama) are **optional/disabled in cloud**.

---

## 1. Goals & Constraints

| Goal | Constraint |
|------|-----------|
| Tenancy | Open SaaS — strangers sign up; needs isolation, quotas, approval flow |
| Budget | <$25/mo at launch |
| Growth | Architecture must scale to 50–500 users / 2–20TB without rework |
| Portability | No cloud-specific APIs; K8s/Helm portable; S3-compatible storage only |
| Durability | Photos are irreplaceable — object storage + independent second copy (accepting up to 24h metadata RPO, see §6) |
| AI | Disabled by default in cloud; GPU cost is the biggest variable |

---

## 2. Current State (codebase findings, verified)

- ✅ **Registration + admin approval** exists (`POST /register`, `User.Approved`) — a working gate for open signups
- ✅ **Storage SDK is portable** — MinIO Go SDK (S3-compatible); works against R2/B2/S3/Oracle/MinIO
- ✅ **Valkey cache tier exists** (thumbnails, geocoder responses, filesystem-scan index) — Redis-compatible; note `CACHE_ENABLED` defaults to **false** (§13.6)
- ✅ **Helm chart + K3s deployment** exist (currently validated against home K3s, ns `media`)
- ❌ **Photos have no owner** — `domain.Photo` has no `UserID`; multi-tenancy needs schema work
- ❌ **No quota concept** anywhere in the API
- ❌ **Originals are filesystem-only** — served via `ctx.File`; need an S3 adapter + serving-path change (S3 stream or presigned URLs) before cloud hosting
- ❌ **Shares are an unscoped IDOR today** — `ListShares`/`RevokeShare` don't scope by user; any authenticated user can list/revoke every share link. `Share.CreatedBy` exists but is ignored by the service layer.
- ❌ **Share passwords are stored plaintext** (marked `// NOT for production` in `service/share.go`)
- ❌ **No billing/subscription plumbing** (needed if paid tiers fund storage at scale)

**Phase 0 prerequisite code work (blocks all hosting options):**
1. `owner_id` on `photos`/`albums` — and also scope **tags, photo analysis, settings, search** (all currently global) — row-level tenancy in one shared DB
2. Per-user storage quota (bytes used / limit), enforced at upload, with an orphan-object sweep (upload-to-S3 succeeds, DB insert fails → object leaks quota)
3. S3-compatible adapter for **originals** + serving-path change (key pattern `originals/{owner_id}/{photo_id}`); thumbnail keys re-key to `thumbnails/{owner_id}/{photo_id}/{size}.webp` (dual-lookup or re-generate)
4. **Share scoping fix + share-password hashing** (live IDOR + plaintext — security blockers for open signup)
5. **AuthZ integration tests as definition of done**: cross-tenant 403 asserted on every photos/albums/shares/search/tags endpoint; `owner_id` indexes on existing query paths
6. Signup hardening: email verification, approval flow polish, stricter auth rate limits

---

## 3. The Core Economic Tension (read first)

Open SaaS + growth + $25/mo only co-exist if **users don't store much, or storage is paid for**. Storage dominates cost past the free tiers — and **durability roughly doubles the storage line** (second copy, §6):

| Scale | Users × avg | Primary storage (B2 $6.95/TB) | + backup copy | Fits $25? |
|-------|------------|-------------------------------|---------------|-----------|
| Launch | 10 users × 20GB = 200GB | ~$1.40/mo | ~$1 (home mirror $0) | ✅ easily |
| Small | 50 users × 40GB = 2TB | ~$14/mo | ~€4.50–9 | 🟡 borderline (~$20–28 all-in) |
| Medium | 200 users × 50GB = 10TB | ~$70/mo | ~$25+ | ❌ |
| Growth | 500 users × 40GB = 20TB | ~$140/mo | ~$45+ | ❌ |

**Implication:** the $25 budget covers launch + early small scale. Beyond ~3TB total, the service must be funded by the **10GB free quota** (§10) and/or paid plans. The architecture below keeps the fixed base at ~$6–12/mo so nearly all budget headroom goes to storage.

---

## 4. Compute/Platform Options

### Option A — Hybrid: home K3s + cloud object storage
Keep compute on the existing home K3s; photos in cloud object storage.

- **Cost:** ~$0 compute; storage only
- **Pros:** cheapest; K3s/Helm skills reused; AI can run on home GPU free
- **Cons:** strangers' traffic flows through home uplink/NAT; home power/ISP is the SPOF; upload throughput to SaaS users is poor; CGNAT may complicate ingress (needs a tunnel)
- **Verdict:** ❌ for open SaaS serving. Fine as a dev/staging lane and as a **free third backup copy** (nightly rclone pull).

### Option B — Single cloud VPS running k3s *(recommended at launch)*
One **netcup VPS** running k3s; self-hosted Postgres on the same box with nightly `pg_dump` shipped **off-provider** (§6); photos in B2; Cloudflare free CDN in front of the API.

netcup G12 lineup (prices incl. 19% VAT, verified Aug 2026):

| VM | vCPU | RAM | NVMe | Price | Fit |
|----|------|-----|------|-------|-----|
| VPS 500 G12 | 2 | 4 GB | 128 GB | €5.91/mo | Budget floor — tight for k3s+API+PG+valkey; pin PG `shared_buffers` low |
| **VPS 1000 G12** | **4** | **8 GB** | **256 GB** | **€10.37/mo** | ✅ **recommended** — comfortable headroom incl. ArgoCD |
| VPS 2000 G12 | 8 | 16 GB | 512 GB | €19.25/mo | Phase 2/3 headroom |
| RS 1000 G12 | 4 dedicated | 8 GB | 256 GB | €12.79/mo | Upgrade if VPS noisy neighbours appear (dedicated cores; 12-mo term saves ~16%) |

Included on all: IPv4, DDoS protection, Copy-On-Write snapshots (free). Traffic "included" with no disclosed cap ⚠️ (verify; netcup is historically generous). German DC only.

- **Cost:** ~$6–12/mo compute + storage
- **Pros:** K8s skills fully reused; Helm chart as-is; zero platform lock-in (box is replaceable — photos in object storage, DB dumps off-provider, IaC rebuild); predictable pricing; generous included traffic; free snapshots
- **Cons:** single node (acceptable at launch); you own Postgres backups; German DC only (latency for non-EU users — fine at launch, revisit if the user base goes transatlantic); no S3 offering (photos stay in B2 by design anyway); smaller English-language community than Hetzner
- **Ops note:** see **§12** for the frictionless K3s bootstrap (cloud-init reference, hardening, Tailscale kubectl access, ArgoCD GitOps option, k3s state backup)
- **Verdict:** ✅ best fit for budget + portability + skills.

### Option C — Oracle Always Free (A1 ARM VM or OKE free control plane)
2 OCPU Ampere ARM + 12GB RAM free; OKE control plane free (pay only workers, which can be Always-Free A1).

- **Cost:** $0 compute
- **Pros:** free; generous free block storage (200GB)
- **Cons:** ⚠️ capacity frequently unavailable in popular regions ("out of host capacity"); Oracle **reclaims idle instances** (7-day <20% utilisation heuristic — a <10-user photo app **will** look idle, and you'd discover the stop via user complaints); ARM-only free shape (verify ffmpeg/imaging deps in CI); ⚠️ "10TB/mo free egress" claim is **unverified**
- **Verdict:** 🟡 use as a **staging/experiment lane**, not production launch. If capacity and utilisation prove stable over months, it can displace Option B later — the Helm chart makes that a config-level move.

### Option D — PaaS containers (fly.io / Railway / Render) + Neon free Postgres
Deploy API as a container; Neon scale-to-zero Postgres (0.5GB free); photos in R2/B2.

- **Cost:** ~$2–7/mo (+ fly.io shared IPv4 ~$2/mo; fly egress $0.02/GB)
- **Pros:** least ops; fast to launch; Neon free tier fits the trivial DB
- **Cons:** per-unit pricing scales worse than a VPS; free tiers sleep/pause (Render) — bad for a photo app; mild platform lock-in
- **Verdict:** 🟡 fine for a quick launch, but Option B is cheaper at equal capability and more portable.

### Comparison summary

| Option | Launch cost | Portability | Ops burden | Risk |
|--------|------------|-------------|------------|------|
| A — Home hybrid | ~$0 + storage | High | Low | Home uplink/power for strangers ❌ |
| **B — netcup VPS + k3s** | **~$6–12 + storage** | **High** | **Low-med** | Single node (mitigated by object storage + off-provider dumps) |
| C — Oracle free | $0 + storage | High | Med | Idle reclamation can stop production ⚠️ |
| D — PaaS | ~$2–7 + storage | Medium | Lowest | Free-tier sleep; unit pricing at scale |

---

## 5. Object Storage Comparison (photos)

All options below are S3-compatible (MinIO SDK works unchanged). Prices from Aug 2026 research; ⚠️ = unverified, confirm before committing.

| Provider | Storage | Egress | Free tier | Gotchas |
|----------|---------|--------|-----------|---------|
| **Backblaze B2** | **$6.95/TB/mo** *(verified)* | Free up to 3× stored/mo; unlimited free via Cloudflare | 10GB | 3× rule can be exceeded in heavy-download months unless fronted by Cloudflare |
| **Cloudflare R2** | $15/TB/mo | **$0** | 10GB + ops | Class A/B op fees ($4.50/M writes, $0.36/M reads); **no Object Lock** (§6) |
| **Wasabi** | $7.99/TB/mo ⚠️ | $0 | Trial only | 90-day minimum storage duration ⚠️; no free tier |
| **Oracle Object Storage** | ~$25.50/TB Standard (~$10/TB infrequent-access) | ⚠️ generous free egress (unverified) | 10GB | Pairs with Option C; Standard tier is 2.5× the price initially assumed |
| **Hetzner Object Storage** | ~€10/TB ⚠️ | 1TB included, ~€1/TB over | None | Pricing not fully public; verify |
| AWS S3 | $23/TB/mo | ~$0.09/GB ❌ | None | Egress kills a photo app; excluded |

**Why "R2's zero egress" is not the deciding factor it appears to be:** the API proxies every byte (originals via `ctx.File`, thumbnails via S3 GET → `ctx.Data`), so origin egress flows to the VPS — and netcup includes traffic with no disclosed cap (⚠️ verify the fair-use terms), making provider egress fees moot at launch either way. What actually controls cost is caching: **Cloudflare free CDN in front of the API, honouring the already-set `Cache-Control: immutable` thumbnail headers — then most thumbnail views never touch the origin or object storage at all.** That is launch architecture, not an optimization.

**Recommendation: B2 from day 1 + Cloudflare free CDN.** With CF caching, B2 is cheaper than R2 at essentially every size ($6.95 vs $15/TB, both 10GB-free), and B2 uniquely offers **Object Lock** (§6) — immutability that survives even a compromised operator credential, which matters most in the early single-operator phase. R2 remains the zero-egress fallback if the CF pairing ever breaks. Re-model R2 Class A/B op counts (multipart uploads ≈3+ ops, trash copy+delete = 2, 3 thumbnail PUTs/photo) before Phase 2 either way. **Caveat:** CF only caches thumbnail traffic if the auth design allows it — see §13.4, which is a Phase-0/1 design decision.

---

## 6. Durability & Backup Strategy (no photo loss)

Object storage provides ~11-nines durability against hardware loss. The residual risks are **account/provider failure**, **accidental/malicious deletion**, and **silent backup failure** — the design must survive all three:

1. **Immutability on the primary:** B2 bucket versioning + **Object Lock (governance mode)** + 90-day noncurrent-version lifecycle. Versioning defends against app bugs; Object Lock defends against credential theft and operator error. 30 days is too short — a wipe unnoticed over a quiet month must not be fatal.
2. **Independent second copy at a different provider — including the database:**
   - Target: **Hetzner Storage Box** (1TB ≈ €4.50/mo, SFTP/restic) **from day 1**. Home mirror is a free **third** copy, not the only off-provider one (the home box is the same SPOF §4 dismisses for serving).
   - Scope: originals + **nightly `pg_dump` + bucket/credential inventory + rclone config**. Thumbnails excluded (regenerable via job — verified the job exists; note regen is not on-miss, so a restore runs the job afterwards). **Without the DB dumps in the second copy, provider-account loss destroys your ability to restore accounts, shares, quotas, and the owner↔object mapping even though photo bytes survive.**
   - Method: the backup **pulls** with **read-only credentials** on the primary; `rclone copy` (never deletes at destination) or restic snapshots; 90-day deletion retention on the backup side. A cron holding write creds on the primary is how one bug or one stolen key wipes both copies in a night.
3. **Postgres:** nightly `pg_dump` off-provider (stated RPO: up to 24h of accounts/approvals/quota/album metadata — acceptable at launch; photos themselves are unaffected).
4. **Detection, not just drills:** scheduled `rclone check` (size+hash compare primary vs backup) with **failure alerting** (healthchecks.io free tier, or Uptime Kuma on the box). A backup silently failing for months is "no backup" discovered too late. Restore drill: full restore once at launch (object counts/sizes vs DB `FileHash` anchors integrity), then quarterly spot-restores.
5. **Deletion semantics change on S3:** filesystem trash is a rename; S3 trash is copy+delete (ordering matters — copy must succeed first), and "purge" only writes delete markers until the lifecycle expires them. That's a durability feature, but: purge frees user quota instantly while storage is still paid for the retention window (model it), and **GDPR erasure cannot complete until versions + backups expire** — document a maximum erasure latency (≤90 days) in the ToS/privacy policy.
6. **Restore path note:** after restoring originals + DB to a fresh deployment, run the thumbnail regeneration job to rebuild the thumbnail tier.

Estimated backup cost: **~€4.50/mo** (Storage Box 1TB) + **$0** home mirror.

---

## 7. Multi-Tenancy Design (shared everything, row-level isolation)

- **One shared Postgres** — `owner_id` on photos/albums **and owner scoping on tags, analysis, search, settings**; every query scoped by owner (admin cross-tenant for ops). DB stays trivial — no per-tenant DBs.
- **One shared bucket** — keys prefixed `originals/{owner_id}/...` and `thumbnails/{owner_id}/...`. No per-user buckets (doesn't scale; all object access flows through the API).
- **Quotas** — `storage_used_bytes` / `storage_limit_bytes` per user, transactional on upload/delete; orphan sweep reconciles leaked objects. Free quota: **10GB** (§10).
- **Signup** — existing register + `Approved` flag; add email verification. **Decided (§10): open signup from launch** (no approval gate — the flag stays as an admin tool); moderation is reactive-reporting (below).
- **Security Phase-0 blockers (live today):** share listing/revocation IDOR, plaintext share passwords, global tags/search leaking across would-be tenants. AuthZ integration tests are the definition of done.
- **Product model — decided (§10):** isolated tenant per signup now (collaboration via share links only), **designed so multi-user tenants can be added later without rework** — all queries scoped by a single ownership dimension (`owner_id`) that a future `tenant_id`/membership layer can replace. Today's home-instance semantics (global roles, album membership) stay on the home instance, not the SaaS.
- **Compliance & moderation — decided (§10):** an EU-hosted photo SaaS holding strangers' personal data needs ToS/privacy policy (and a DPA posture). Moderation is **reactive reporting only**: minimal report channel at launch (mailto link + ToS language), report button + takedown workflow in Phase 2. With Ollama disabled there is no automated NSFW detection (`PhotoAnalysis.IsNSFW` is the only hook) — this exposure is accepted and tracked in §11.
- **Encryption:** server-side encryption (default on B2/R2) + TLS in transit is the industry norm at this tier. Per-tenant/client-side E2E is rejected (breaks thumbnailing and AI) — YAGNI.

---

## 8. Recommended Phased Path

| Phase | Scale | Stack | Est. cost |
|-------|-------|-------|-----------|
| **0 — prereqs** | — | Code: ownership (all resources), quotas, S3 originals + serving path, share security fixes, authZ tests, signup hardening | $0 |
| **1 — launch** | <10 users, <250GB | **Option B**: netcup VPS 1000 G12 + k3s (bootstrap per §12; plain Helm or ArgoCD); B2 (versioning + Object Lock); CF free CDN with immutable thumbnail caching; self-hosted PG + nightly dumps **to the Storage Box**; Storage Box backup (rclone pull, read-only, 90-day retention); home mirror third copy; healthchecks.io alerting; Resend free tier for email; open signup + mailto report channel + ToS; **landing page on CF Pages (§15)**; Option C as staging lane only | **~$11–16/mo** (VPS 500 G12 €5.91 or VPS 1000 G12 €10.37 + backup €4.50; storage ≈ free tier) |
| **2 — early SaaS** | 10–50 users, 0.25–2TB | Same shape; B2 storage line grows; enlarge Storage Box; quotas enforced (10GB free); report button + takedown workflow; billing decision revisited (§10) | **~$25–35/mo** (e.g. 2TB: VPS 1000 G12 ~$11 + B2 $14 + Storage Box 2TB ~$10) — brushes the budget ceiling; paid tiers land by here |
| **3 — growth** | 50–500 users, 2–20TB | Storage is the bill (§3). Fund via paid tiers **before** this phase; then multi-node k3s and managed PG become options when revenue covers them | **$70–200+/mo, revenue-funded** |

Fixed cost lines to budget from day 1: domain ~$12/yr; email $0 (Resend free, 3k/mo); monitoring $0 (healthchecks.io); VM snapshots $0 (netcup COW snapshots included — still not a substitute for off-box backups).

**Recommendation: keep the home K3s as your personal instance; stand the SaaS up as a fresh deployment.** Avoids a painful data cutover, keeps personal photos off the SaaS box, and dodges the owner-backfill problem for existing data entirely.

---

## 9. Lock-in Guardrails

- S3-compatible storage only (MinIO SDK already) — no provider-native APIs
- k3s + Helm for all deployment (chart exists); plain `helm install` is the zero-overhead path, ArgoCD GitOps is a supported option (§12.4) that mirrors the home workflow
- **OpenTofu** (not Terraform) for cloud resources — BUSL-free, drop-in compatible, built-in state encryption (§14)
- Postgres portable (plain `pg_dump` restore path)
- Cloudflare used only for commodity features (CDN/DNS); avoid Workers/Pages/Images
- Backup copy at an **independent** provider, **pull-based with read-only primary credentials**, so exit from any provider is always possible
- ARM-capable builds verified in CI (keeps the Oracle free-tier escape hatch open)

---

## 10. Decisions (locked 2026-08-15)

| # | Question | Decision | Rationale / consequence |
|---|----------|----------|------------------------|
| 1 | Free quota per user at launch | **10GB** | Tight cost curve: 50 users ≈ 500GB ≈ $3.50/mo B2. Paid tiers add more (Phase 2/3) |
| 2 | Product model | **Isolated tenant per signup now, multi-user later** | Phase 0 ships `owner_id` scoping; schema must admit a future tenant/membership layer without rework (queries written against a single ownership dimension that a `tenant_id` mapping can later replace) |
| 3 | Billing | **Defer to Phase 2/3** | No billing code in Phase 0/1; revisit Stripe vs merchant-of-record when quotas pinch (MoR attractive for EU VAT) |
| 4 | Content moderation | **Reactive reporting only** | Open signup (no approval gate); minimal report channel at launch (mailto + ToS language), report button + takedown workflow in Phase 2. **Accepted risk:** legal/operational exposure until then (§11) |
| 5 | Compute option for launch | **Option B (netcup VPS)** — resolved in review | Oracle free tier's idle reclamation is disqualifying for production; C is the staging lane |
| 6 | Backup topology | **Confirmed: Storage Box 2nd copy, home mirror 3rd** | Off-prem second copy survives provider loss; home mirror is a free bonus, not a dependency |
| 7 | Infra repo | **New `photobox-infra` repo** (§14.4) | tofu/Ansible/cloud-init change on a different cadence from app deploys (ArgoCD watches helm-charts) |
| 8 | Landing page | **Static site (Astro/Hugo) on Cloudflare Pages, planned now / built in Phase 1** (§15) | $0 on CF free tier; apex = landing, `app.*` = the app; needed at launch for the open-signup funnel |

---

## 11. Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Silent backup failure | "No backup" discovered at disaster time | `rclone check` + healthchecks.io alerting (§6.4) |
| Live share IDORs / plaintext share passwords | Cross-tenant data breach on day 1 of open signup | Phase 0 blockers (§2, §7) |
| Provider account loss | Loss of restore ability even with bytes safe | Second copy includes DB dumps + config inventory (§6.2) |
| Credential compromise / bad sync | Both copies wiped | Object Lock + pull-based read-only backup + 90-day retention (§6) |
| Single VPS failure | Downtime (not data loss) | Photos in object storage; dumps off-provider; IaC rebuild |
| Storage cost outruns budget at growth | Service unfundable | Hard quotas from day 1; paid tiers before Phase 3 |
| Unverified prices (Oracle egress/storage, Hetzner obj storage, Wasabi min) | Wrong cost model | Verify before committing (§5 flags) |
| GDPR erasure vs retention | Compliance gap | Documented max erasure latency; ToS/privacy policy (§6.5, §7) |
| No content moderation with AI off | Legal/operational exposure | **Accepted risk (§10):** reactive reporting only — mailto channel + ToS at launch, report/takedown workflow Phase 2 |
| netcup traffic cap undisclosed ("included") | Surprise throttle/fee at photo-app traffic levels | Verify fair-use terms before launch; CF CDN absorbs most reads (§5) |
| German-DC-only (netcup) | Latency for non-EU users | Acceptable at launch; CF cache hides most of it; revisit if user base shifts |

---

## 12. Appendix — Running K3s on the VM (frictionless path)

How to go from a fresh netcup VM to a working cluster in ~10 minutes, in the order you'd do it.

### 12.1 Bootstrap: cloud-init + official install script

For a single node, **cloud-init running the official install script** is the lowest-friction, most reproducible path — the VM self-provisions on first boot, no SSH choreography needed. `k3sup` (still maintained) is the ad-hoc alternative from your laptop; `k3s-ansible` only earns its keep at 2–3 nodes; Cluster API is fleet tooling — overkill.

Reference cloud-init (paste into the netcup order panel / apply post-install):

```yaml
#cloud-config
packages: [curl, rsync, fail2ban, unattended-upgrades]

write_files:
  - path: /etc/rancher/k3s/config.yaml
    permissions: "0644"
    content: |
      write-kubeconfig-mode: "0644"
      tls-san: ["api.photobox.example.com"]
      disable:
        - traefik      # Cloudflare terminates TLS + caches
        - servicelb    # no cloud LB behind CF

runcmd:
  - swapoff -a && sed -i '/swap/d' /etc/fstab   # k3s refuses to start with swap on
  - curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server" sh -
  - systemctl enable --now fail2ban unattended-upgrades
```

Result: k3s up with the bundled sqlite datastore (right choice for one node — no etcd overhead), kubeconfig at `/etc/rancher/k3s/k3s.yaml`.

### 12.2 Hardening baseline (10 minutes, once)

- SSH: key-only auth, `PermitRootLogin no`
- UFW: deny incoming; allow SSH, and 6443 **only from your IP/Tailscale** (or rely on netcup's DDoS protection + CF for 443; never expose VXLAN 8472 — that's L2 access to the cluster network)
- `fail2ban` + `unattended-upgrades` (in the cloud-init above)
- Swap off (already in cloud-init); Ubuntu 22.04+ cgroup v2 works out of the box

### 12.3 kubectl access: Tailscale

Install Tailscale on the VM and your laptop (free tier covers this easily), then kubectl over the WireGuard tunnel — **6443 never touches the public internet**. IP allowlist on 6443 is the acceptable alternative if you have a static IP; fully public 6443 is not recommended.

### 12.4 GitOps: ArgoCD (optional but matches your home workflow)

Plain `helm install photobox run/chart/` is the zero-overhead path and entirely sufficient for one node. **ArgoCD is worth it if you want the same push-to-`kubernetes-helm-charts` → auto-sync workflow you run at home:**

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

- **Resource cost:** ~0.5–1GB RAM for the full install — comfortable on the VPS 1000 G12 (8GB), too tight on the VPS 500 (4GB); a `--core` (serverless) install halves that if ever needed
- **Pattern:** one `Application` pointing at the GitLab `r1chjames/kubernetes-helm-charts` repo, auto-sync + prune + self-heal; ArgoCD UI exposed only over Tailscale (no public ingress)
- Flux is the lighter alternative (~half the footprint) if ArgoCD's UI isn't valued — but you already know ArgoCD, so there's no re-learning argument for switching

### 12.5 Backing up the k3s control plane (sqlite — not etcd)

The `--etcd-snapshot-*` flags **do not apply** to the default sqlite datastore. Back up state with a nightly cron, shipped to B2 alongside the PG dumps:

```bash
# /etc/cron.daily/k3s-state-backup
tar -czf - /var/lib/rancher/k3s/server/db \
            /var/lib/rancher/k3s/server/token \
            /var/lib/rancher/k3s/server/tls \
  | rclone rcat b2:photobox-backup/k3s/$(date +%F).tar.gz
```

The **token** and **tls** dirs are as important as the DB — without them a restore can't re-join agents or verify identity. With GitOps (§12.4) the cluster's desired state is also in git, so this backup is mostly for secrets and speed of rebuild.

### 12.6 Gotchas worth knowing upfront

| Gotcha | Handling |
|--------|----------|
| Swap enabled → k3s won't start | `swapoff -a` + remove from fstab (cloud-init above) |
| Port 8472 (VXLAN) exposed | UFW default-deny; netcup DDoS protection is not a firewall |
| MTU weirdness behind tunnels | If pods lose connectivity over Tailscale/CF tunnel, check flannel MTU first |
| Old iptables (1.8.0–1.8.4) | Only on older distros; use a current Ubuntu LTS image |
| Single-node disk pressure | 256GB NVMe is ample; set `systemReserved`/`evictionHard` if you later fill it with image layers |

### 12.7 Provisioning flow (the whole thing)

1. Order VPS 1000 G12 with the cloud-init above → coffee → `kubectl get nodes` works over Tailscale
2. Point Cloudflare DNS at the VM (proxied), deploy CF origin certs (or CF Tunnel if you'd rather not open 443)
3. `helm install photobox run/chart/` (or bootstrap ArgoCD → Application → auto-sync)
4. Wire the two backup crons (PG dump → Storage Box; k3s state → B2) + healthchecks.io pings
5. Done: a replaceable box whose only irreplaceable data lives off-provider

---

## 13. Storage Performance (NAS vs S3, and how to make S3 fast)

### 13.1 Framing: what S3 does and doesn't win

S3 won't beat a LAN NAS on single-request latency — it wins on parallel throughput and consistency under load. Perceived speed comes from the caching layers in front of the origin, not the origin itself.

| | NAS (LAN) | S3 (B2/R2) from the VPS |
|---|---|---|
| Single GET latency (TTFB) | ~1–5ms idle, 50ms+ loaded | ~40–150ms in-region, very consistent |
| Parallel throughput | Collapses (one spindle queue) | Scales ~linearly with concurrency |
| Under load | Degrades hard | Flat |
| Bottleneck becomes | NAS disks | VPS↔bucket latency on **cache misses only** |

### 13.2 Why the current NAS is slow (confirm before blaming the medium)

Likely culprits for grid sluggishness on the home setup:
- **NFS/SMB metadata latency** — every thumbnail read is `open/stat/read/close` over the network; one grid page = 50+ round-trips
- **Spinning-disk random I/O** — hundreds of small-file reads contend on one spindle queue (head-of-line blocking)
- **K3s PVC overhead** — NFS provisioner or Longhorn replication adds another network hop per read
- **Diagnostic:** `iostat -x 5` on the NAS while scrolling a large album — `%util` near 100% with high `await` = spindle contention, not bandwidth

### 13.3 Getting good S3 performance (origin-side)

1. **Co-locate region** — bucket in the same geography as the VPS (B2 EU region: Amsterdam — pairs well with netcup Nuremberg; R2 supports an EU location hint). Biggest single lever on cache-miss latency.
2. **Concurrency** — MinIO SDK pools connections; thumbnail fetches must be parallel (webapp already does 6-concurrent), never serial waterfalls.
3. **Never LIST on the hot path** — LISTs are slow and metered; listing already comes from Postgres. Keep it that way.
4. **Right-size every fetch** — grid requests `s` thumbnails only; never proxy an original to render a card.
5. **Range GETs for video/streaming** — no whole-object reads to seek (assumed by #103's design).
6. **Multipart parallel-chunk uploads** for 50MB+ HEICs (aligns with #98).

### 13.4 Cache layers — and the auth caveat that decides whether they work

```
Browser → Cloudflare edge (goal: ~80% of thumbnail reads) → VPS nginx cache → Valkey → B2
```

Cold miss: ~50–150ms once, then warm for everyone. The app's `Cache-Control: immutable` thumbnail headers are already correct for this.

**⚠️ The caveat the rest of this plan depends on:** Cloudflare **does not cache responses to requests carrying an `Authorization` header by default** — and the thumbnail routes are authenticated today. For edge caching to work, pick one during the S3-originals work (Phase 0/1):
- **(a) Auth-free thumbnail route against unguessable UUID photo IDs** — simplest; treat the UUID as the capability (acceptable for thumbnails, not originals)
- **(b) Short-lived signed thumbnail URLs** — stronger; more moving parts
- **(c) CF Worker that validates auth and strips the header** — keeps server-side authZ; adds CF-specific code (mild guardrail tension, §9)

This choice is the difference between "CF absorbs ~80% of reads" and "CF absorbs 0%". Multi-tenant note: never cache originals at a shared CDN keyed only by URL — thumbnails only, and only with one of the above.

### 13.5 Performance targets (acceptance numbers for the hosting work)

| Path | Target |
|------|--------|
| Warm edge (CF hit) thumbnail | <20ms |
| VPS-cached (nginx/Valkey) thumbnail | <10ms |
| Cold S3 miss (in-region) | 50–150ms TTFB |
| Full grid first paint (50 thumbnails, warm) | <1s |

If production shows multi-second grid loads, the origin is not the problem — the cache chain is broken (check CF cache-status headers first).

### 13.6 Valkey (Redis-compatible) cache tier — role & sizing

The app already ships a Valkey cache (`CacheService`, go-redis client, graceful degradation) caching thumbnails, geocoder responses, and the filesystem-scan index — but `CACHE_ENABLED` **defaults to false**. Two consequences:

- **Today (NAS):** if the current deployment hasn't set `CACHE_ENABLED=true`, it's running with zero thumbnail caching — enabling it is the cheapest performance win available, no migration needed.
- **In the cloud architecture, Valkey gets more important, not less:** it absorbs S3 misses for *authenticated* traffic — exactly the traffic Cloudflare won't cache (§13.4). Until the thumbnail-auth design is resolved, Valkey **is** the primary thumbnail cache. Chain: `CF edge → nginx → Valkey → B2`.

**Sizing on the VPS 1000 G12 (8GB):** WebP thumbnails run ~20–50KB, so 1GB ≈ 20–50K thumbnails — comfortably covers a small SaaS hot set. Run as a pure cache:

```
maxmemory 1gb
maxmemory-policy allkeys-lru
# no AOF/RDB persistence — contents are re-fetchable from S3 or regenerable;
# a restart costs only a brief warm-up
```

**Not for Valkey:** originals (they'd evict everything), and defer caching API responses / job status / session lookups until profiling justifies it (YAGNI).

---

## 14. Infrastructure as Code (everything around the Helm chart)

The Helm chart already manages the full app stack (deployment, ingress, PVC, secrets, HPA, Ollama sidecar, Postgres/Valkey as chart dependencies). The IaC gap is everything around the cluster: VM, DNS/CDN, buckets, secrets, backup wiring. Layered approach — each layer uses the lightest tool that covers it:

### 14.1 What's automatable vs one-time manual

| Component | Automatable? | Tool |
|-----------|-------------|------|
| netcup VM | ❌ **one-time manual order** — no official provisioning API; community TF providers are stale (last commits 2022). No loss: DNS is on Cloudflare, so no netcup provider is needed at all | netcup panel |
| VM day-0 config (swap, k3s, hardening) | ✅ | cloud-init (§12.1) |
| VM day-2 config (rclone, backup crons, Tailscale, argocd install) | ✅ | Ansible (`ansible-pull` from the infra repo) |
| Cloudflare (DNS records, cache rules, origin certs) | ✅ official provider | OpenTofu |
| B2 buckets + lifecycle + Object Lock + app keys | ✅ official provider (`backblazeb2/backblazeb2`) | OpenTofu |
| Tailscale (auth keys, ACLs) | ✅ official provider | OpenTofu |
| healthchecks.io monitors | ✅ community provider (or 3-line curl script) | OpenTofu / script |
| Hetzner Storage Box | ❌ no provider (Robot-side product, no API) — created once manually, then configured via Ansible/rclone | panel + Ansible |
| K8s cluster state (app, ArgoCD apps, secrets CRs) | ✅ | ArgoCD app-of-apps (§12.4) |
| Secrets in git | ✅ | Sealed Secrets (K8s) + SOPS/age (Ansible) |
| OpenTofu state | ✅ | R2 S3 backend (officially documented; free tier; zero egress) |

### 14.2 Toolchain decisions

- **OpenTofu over Terraform** — BUSL-free, drop-in compatible (same providers/registry), built-in state encryption. Matches the no-lock-in constraint.
- **State backend: R2**, not B2 — Cloudflare officially documents the S3 backend on R2; state is tiny (free tier); no DynamoDB-style locking on either, acceptable for a single operator.
- **Secrets — split by side, respecting ArgoCD's own warning.** ArgoCD officially recommends *against* argocd-vault-plugin and friends: injected secrets land **plaintext in ArgoCD's Redis cache**. So:
  - **K8s side: Sealed Secrets** (Bitnami, active) — encrypted `SealedSecret` manifests live in git, the on-cluster controller decrypts them; ArgoCD syncs ciphertext and never sees values. No external vault to run (ESO is the alternative if you ever want one — rejected for now, YAGNI).
  - **Ansible/OS side: SOPS + age** — encrypted rclone configs, Tailscale auth keys, cron credentials in the infra repo.
- **Ansible via `ansible-pull`** from the infra repo — the VM configures itself from git; no control-node dependency, works for rebuilds.

### 14.3 Automation chain (empty account → running SaaS)

```
[1. manual]  Order netcup VPS + Storage Box (once, ~10 min in panels)
      ↓
[2. tofu]    tofu apply  →  CF DNS/cache rules, B2 buckets (+Object Lock),
                            Tailscale auth key, healthchecks checks, state in R2
      ↓
[3. cloud-init]  First boot: swap off, hardening base, k3s install (§12.1),
                 then `ansible-pull -U <infra-repo> site.yml`
      ↓
[4. ansible]  rclone remotes (B2, Storage Box), backup crons (PG dump, k3s
              state → §12.5), Tailscale up, helm install argocd
      ↓
[5. argocd]  app-of-apps takes over: photobox chart (+sealed secrets),
             everything else in-cluster from then on
      ↓
Rebuild test: kill the VM → repeat steps 1–5 → restore PG dump +
              rclone check originals → thumbnail regen job (§6.6)
```

### 14.4 Infra repo layout (slim)

```
infra/
├── tofu/                  # OpenTofu: cloudflare.tf, b2.tf, tailscale.tf,
│   └── backend.tf         #   healthchecks.tf; backend = R2 S3
├── ansible/
│   ├── site.yml           # roles: base, tailscale, k3s, rclone, backups, argocd
│   └── group_vars/all/vault.yml   # SOPS+age encrypted
├── k8s/
│   ├── app-of-apps.yaml   # ArgoCD root app → kubernetes-helm-charts repo
│   └── sealed-secrets/    # SealedSecret manifests (ciphertext in git)
├── cloud-init/user-data.yaml
└── .sops.yaml             # age recipients
```

**Infra repo — decided (§10):** new `photobox-infra` repo. It keeps app-deploys (ArgoCD watches helm-charts) and infra (tofu/Ansible, run by hand) on different cadences.

### 14.5 What this buys

- **The VM is cattle:** full rebuild = order VM + `tofu apply` + boot + ~15 min of automation. Combined with §6 (all state off-box), the single-node risk is reduced to downtime only.
- **Every drift-prone thing is in git:** DNS, cache rules, bucket policies, backup crons, cluster state.
- **No lock-in added:** OpenTofu, Ansible, Helm, SOPS, Sealed Secrets are all FOSS; the only provider-specific resource blocks are thin (CF records, B2 buckets).

---

## 15. Product Landing Page

An open-signup SaaS needs a landing page at launch — it *is* the signup funnel. **Decided (§10.8): static site on Cloudflare Pages, planned now, built in Phase 1.**

### 15.1 Stack & hosting

- **Static site** — Astro (component-friendly, ships zero JS by default) or Hugo (simpler); final pick at build time, both are equally portable. Plain HTML output means zero runtime to maintain.
- **Cloudflare Pages free tier** — $0, git-integrated deploys (push → build → global CDN), custom domains + TLS included. Fits the budget exactly.
- **Lock-in check:** negligible — the output is static files; redeploying to GitHub Pages/Netlify/S3+CDN is a one-line change if CF Pages ever disappoints. Only the build/deploy integration is CF-specific.
- **Domain layout:** apex (`photobox.example`) = landing; `app.photobox.example` = the webapp; `www` redirects to apex. DNS + Pages project managed in OpenTofu (cloudflare provider covers both) — lands in the `photobox-infra` repo (§14.4).
- **Repo:** `landing/` folder in the main photobox repo — CF Pages builds from a subdirectory; keeps product + marketing versioned together.

### 15.2 Positioning & content outline

Positioned between the two incumbents the codebase already benchmarks against:

> **Google Photos convenience, without Google.** Immich proves people want to own their photos — but not everyone wants to run a server. Photobox is the hosted version of that promise: your library, EU-hosted, nobody training on your memories.

Sections:
1. **Hero** — the one-liner above + "Create free account" CTA (10GB free, §10.1)
2. **Feature grid** — AI tagging & search, map view, albums & sharing, Memories/"On This Day", video support (as they ship — keep the page honest about what's live)
3. **Import story** — Google Takeout + Immich migration tools (#106/#107): "bring your library with you" — the strongest switching argument
4. **Pricing teaser** — 10GB free; paid tiers "coming" (Phase 2/3, §10.3)
5. **Trust strip** — EU hosting, no ads/tracking, open source (link to repo)
6. **Footer** — ToS, privacy policy, **report-abuse mailto** (carries the compliance links required by the moderation decision, §7/§10.4)

### 15.3 Timing

Phase 1 launch item — build with @designer when Phase 1 starts (real screenshots from the running app, not mockups). Until then the apex can serve a minimal "coming soon" page from the same CF Pages project.

---

*Pricing data: Aug 2026 research pass against provider pricing pages (B2 price and netcup G12 lineup independently verified). Items flagged ⚠️ were not verifiable and must be confirmed before use in a commitment.*
