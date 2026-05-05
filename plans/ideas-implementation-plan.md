# Ideas Implementation Plan

## 1. Multi-Tenancy

### Overview
Support multiple isolated users/organisations sharing a single Photobox instance, each with their own photos, albums, and settings.

### Option A: Row-Level Security (RLS) + Tenant ID Column (Recommended)
Add a `tenant_id` column to every table and enforce filtering at the application layer. PostgreSQL RLS policies can optionally enforce this at the DB level.

**Pros:**
- Single database, single schema — simplest ops
- Can add RLS policies later for defence-in-depth
- Minimal infrastructure change

**Cons:**
- All queries must include `tenant_id` filter (easy to miss)
- No true resource isolation (one tenant's heavy query affects others)

**Implementation Steps:**
1. Add `tenant_id uuid` to `users`, `albums`, `photos`, `settings`, `jobs`, `shared_links`
2. Update `User` model with `TenantID` field; assign new users to a default tenant
3. Add middleware that extracts tenant from JWT claim (`tenant_id`) and sets it on Gin context
4. Update **all** repository methods to filter by `tenant_id`:
   - `WHERE tenant_id = ?` in every SELECT/UPDATE/DELETE
5. Update `AuthService` to include `tenant_id` in JWT claims
6. Migration: backfill existing data with a default tenant UUID
7. Add `TENANT_ISOLATION` env var (default `false` for backward compatibility)
8. Update `appconfig.go` to read default tenant ID

---

### Option B: Schema-Per-Tenant
Each tenant gets their own PostgreSQL schema (`tenant_a`, `tenant_b`, etc.) with identical tables.

**Pros:**
- Strong data isolation
- Can migrate individual tenants independently
- Easy to export/delete a tenant

**Cons:**
- Complex connection management (schema search_path switching)
- GORM AutoMigrate must run per schema
- Higher operational overhead (N schemas to monitor/backup)

**Implementation Steps:**
1. Create a `tenants` table to track schema names and statuses
2. Add `TenantResolver` middleware that sets `search_path` on the DB connection per request
3. Update `InitDbConnection` to accept a schema name; create a connection pool per tenant
4. Run `AutoMigrate` for each tenant schema on creation
5. Update all queries to use the tenant-specific schema
6. Admin endpoints for provisioning/deprovisioning tenants

---

### Option C: Database-Per-Tenant
Each tenant gets a completely separate PostgreSQL database.

**Pros:**
- Maximum isolation (can run on separate hardware/region)
- Easiest to reason about security

**Cons:**
- Highest operational overhead
- Connection pooling becomes complex
- Hard to run cross-tenant analytics

**Implementation Steps:**
1. Create a `tenant_registry` database with a `tenants` table
2. Store connection strings per tenant
3. Build a connection factory that resolves tenant → DB connection
4. Run migrations per database
5. Update all service constructors to accept a tenant-aware DB env

---

### Recommended: Option A (Row-Level Security)
**Rationale:** Photobox is currently a single-tenant personal app. Option A adds multi-tenancy with the least risk and ops overhead. If tenant isolation needs increase later, migrate to Option B.

**Files to Reference:**
| File | Change |
|------|--------|
| `api/internal/core/domain/*.go` | Add `TenantID` to all domain models |
| `api/internal/adapter/storage/database/repository/*.go` | Add `tenant_id` filter to all queries |
| `api/internal/core/service/auth.go` | Add `tenant_id` to JWT claims |
| `api/internal/adapter/handler/http/router.go` | Add tenant resolution middleware |
| `api/internal/appconfig/appconfig.go` | Add `DEFAULT_TENANT_ID` env var |

---

---

## 2. AI Integration for Photo Analysis

### Overview
Use AI/ML to automatically tag, caption, classify, and search photos by content.

### Option A: Local Model (Ollama / llama.cpp) (Recommended for Privacy)
Run a local vision-language model (e.g., `llava`, `bakllava`, `moondream`) via Ollama on the same server or a nearby GPU machine.

**Pros:**
- Zero external data egress — photos never leave your infrastructure
- No per-image API costs
- Works offline

**Cons:**
- Requires GPU for reasonable speed (CPU is ~5-10s per image)
- Model quality lower than cloud APIs
- Operational burden of hosting the model

**Implementation Steps:**
1. Add `OLLAMA_HOST`, `OLLAMA_MODEL` env vars to `appconfig.go`
2. Create `api/internal/core/service/ai.go` with an `AIService` interface:
   ```go
   type AIService interface {
       AnalyzeImage(imagePath string) (*ImageAnalysis, error)
   }
   type ImageAnalysis struct {
       Caption   string   `json:"caption"`
       Tags      []string `json:"tags"`
       Objects   []string `json:"objects"`
       IsNSFW    bool     `json:"is_nsfw"`
       IsPortrait bool    `json:"is_portrait"`
   }
   ```
3. Implement `OllamaClient` that:
   - Reads the image file as base64
   - POSTs to `OLLAMA_HOST/api/generate` with a structured prompt
   - Parses JSON response
4. Add `AI_ENABLED` env var; when enabled, call `aiSvc.AnalyzeImage()` inside `SavePhoto`/`SavePhotos` after writing to DB
5. Store AI-derived tags in the `photo_tags` junction table with a source flag (`source = 'ai'`)
6. Add a background job (`scheduler.go`) to re-analyze existing photos in batches
7. Update search to include AI tags in full-text search index

**Prompt template example:**
```
Analyze this image. Return ONLY valid JSON in this format:
{"caption": "...", "tags": ["..."], "objects": ["..."], "is_nsfw": false, "is_portrait": false}
```

---

### Option B: Cloud API (OpenAI GPT-4o / Google Gemini / AWS Rekognition)
Send images to a managed vision API for analysis.

**Pros:**
- Best accuracy and speed
- No GPU infrastructure needed
- Multiple vendors to choose from

**Cons:**
- Cost per image (can be significant for large libraries)
- Privacy concerns — photos leave your infrastructure
- Rate limits and latency

**Implementation Steps:**
1. Add `AI_PROVIDER` (`openai`, `google`, `aws`) and API key env vars
2. Create `api/internal/core/service/ai.go` with provider-specific implementations:
   - `OpenAIClient` → uses GPT-4o vision API
   - `GeminiClient` → uses Google Gemini Pro Vision
   - `RekognitionClient` → uses AWS Rekognition `DetectLabels`
3. Implement a provider factory:
   ```go
   func NewAIService(config appconfig.AppConfig) port.AIService { ... }
   ```
4. Same integration points as Option A (SavePhoto hook, background job, search index)
5. Add cost tracking: log tokens/API calls per image

---

### Option C: Hybrid (Cloud + Local)
Use cloud API for initial high-quality analysis, then fine-tune a local model on the generated labels for future offline use.

**Pros:**
- Best of both worlds
- Can transition to fully local over time

**Cons:**
- Most complex implementation
- Requires ML ops expertise

**Implementation Steps:**
1. Implement Option B first for bootstrapping
2. Store all cloud API responses in a new `ai_training_data` table
3. Periodically export training data and fine-tune a local model
4. Switch to Option A once local model quality is acceptable

---

### Recommended: Option A (Local with Ollama)
**Rationale:** Photobox is a self-hosted personal photo app. Privacy is paramount — users won't want their photos sent to third parties. Start with Ollama + `moondream` (fast, small, decent quality). If quality is insufficient, upgrade GPU or switch to Option B.

**Files to Reference:**
| File | Change |
|------|--------|
| `api/internal/appconfig/appconfig.go` | Add `AI_ENABLED`, `OLLAMA_HOST`, `OLLAMA_MODEL` |
| `api/internal/core/port/ai.go` (new) | `AIService` interface |
| `api/internal/core/service/ai.go` (new) | Ollama/GPT-4o implementations |
| `api/internal/core/service/photo.go` | Call AI analysis after `SavePhoto` |
| `api/internal/adapter/storage/database/repository/photos.go` | Store AI tags with source flag |
| `api/internal/components/scheduler.go` | Background re-analysis job |

---

---

## 3. Object/Cloud Storage for Photos and Thumbnails

### Overview
Replace local filesystem storage with S3-compatible object storage (AWS S3, MinIO, Backblaze B2, Cloudflare R2) for photos, thumbnails, and backups.

### Option A: S3-Compatible Primary Storage (Recommended)
Store all photos and thumbnails in an S3 bucket. The local filesystem becomes a cache only.

**Pros:**
- Infinite scalability
- Built-in redundancy and durability
- Can serve photos via CDN
- Easy backups (S3 versioning, cross-region replication)

**Cons:**
- Latency for first fetch (mitigated by local cache)
- Cost for storage and egress
- Requires internet connectivity

**Implementation Steps:**
1. Add S3 config to `appconfig.go`:
   ```go
   S3Enabled      bool
   S3Endpoint     string   // e.g., s3.amazonaws.com or minio:9000
   S3Bucket       string
   S3Region       string
   S3AccessKey    string
   S3SecretKey    string
   S3UseSSL       bool
   S3Prefix       string   // e.g., "photobox/"
   ```
2. Install `github.com/minio/minio-go/v7`
3. Create `api/internal/adapter/storage/s3/repository.go`:
   ```go
   type S3Repository struct {
       client *minio.Client
       bucket string
       prefix string
   }
   ```
   Implement:
   - `UploadObject(key string, data []byte) error`
   - `DownloadObject(key string) ([]byte, error)`
   - `DeleteObject(key string) error`
   - `GetPresignedURL(key string, expiry time.Duration) (string, error)`
4. Create `api/internal/core/port/storage.go` interface:
   ```go
   type ObjectStorage interface {
       Upload(key string, data []byte) error
       Download(key string) ([]byte, error)
       Delete(key string) error
       PresignedURL(key string, expiry time.Duration) (string, error)
   }
   ```
5. Update `FilesystemRepository` / `FilesystemService` to:
   - Upload original photos to S3 on `SavePhoto`
   - Store `s3://bucket/prefix/photoId` in `Photo.FilesystemPath`
   - Download to local temp cache on `GetPhotoBin` if not cached
   - Upload thumbnails to S3 after generation
6. Update `GetPhotoBin` handler to return a 302 redirect to a presigned S3 URL instead of `ctx.File()` when `S3_ENABLED=true`
7. Update `DeletePhoto` to also delete from S3
8. Add a local LRU cache for recently accessed photos/thumbnails

---

### Option B: Local Primary + S3 Backup
Keep local filesystem as primary storage, sync to S3 asynchronously for backup and remote access.

**Pros:**
- Fast local access for common operations
- S3 acts as durable backup
- Works offline

**Cons:**
- Two sources of truth
- Sync logic adds complexity
- Local disk still the bottleneck

**Implementation Steps:**
1. Same S3 client setup as Option A
2. Add a background job (`scheduler.go`) that:
   - Scans local photos not yet in S3
   - Uploads them in batches
   - Updates `Photo` with `backup_path = s3://...`
3. On `DeletePhoto`, delete from both local and S3
4. On disaster recovery, add an endpoint to restore from S3

---

### Option C: MinIO On-Premise
Self-host MinIO on the same network for S3-compatible storage without cloud vendor lock-in.

**Pros:**
- Same API as S3
- No egress costs
- Full data control

**Cons:**
- You manage the storage hardware
- Still need backups of MinIO itself

**Implementation Steps:**
1. Same as Option A, just point `S3_ENDPOINT` at your MinIO instance
2. Run MinIO in Docker alongside Photobox

---

### Recommended: Option A (S3-Compatible Primary Storage)
**Rationale:** Photobox already moved thumbnails to the filesystem in Phase 2. The next logical step is to decouple from the local disk entirely. Using S3 as primary with local cache enables CDN integration, multi-region deployment, and eliminates disk space concerns. Start with MinIO for local testing, then switch to Cloudflare R2 (no egress fees) or AWS S3 for production.

**Files to Reference:**
| File | Change |
|------|--------|
| `api/internal/appconfig/appconfig.go` | Add all `S3_*` env vars |
| `api/internal/core/port/storage.go` (new) | `ObjectStorage` interface |
| `api/internal/adapter/storage/s3/repository.go` (new) | MinIO/S3 client implementation |
| `api/internal/adapter/storage/filesystem/repository/filesystem.go` | Add local cache layer |
| `api/internal/core/service/photo.go` | Upload to S3 on save, use presigned URLs |
| `api/internal/adapter/handler/http/photo.go` | Redirect to presigned URL for binaries |
| `api/internal/components/scheduler.go` | Background sync/backup job |

---

---

## 4. Serverless Architecture

### Overview
Decompose Photobox into serverless functions for on-demand scaling, reduced cost for low-traffic instances, and event-driven processing.

### Option A: Function-as-a-Service for Background Jobs (Recommended)
Keep the main API as a container (or lightweight VM), but move heavy background work to serverless functions.

**Jobs to serverless-ify:**
- Photo indexing (CPU + I/O intensive)
- Thumbnail generation (CPU intensive)
- AI analysis (GPU intensive)
- Trash cleanup / batch operations

**Pros:**
- API server stays lean and responsive
- Background jobs scale independently
- Pay only for actual compute time
- Can use GPU instances only when needed

**Cons:**
- Adds infrastructure complexity (event queue, function deployments)
- Cold start latency for infrequent jobs
- Debugging distributed systems is harder

**Implementation Steps:**
1. Add a message queue (SQS, NATS, or RabbitMQ) or use PostgreSQL as a job queue
2. Create a `jobs` table:
   ```sql
   CREATE TABLE photobox.jobs (
       id UUID PRIMARY KEY,
       type VARCHAR(50) NOT NULL,  -- 'index', 'thumbnail', 'ai'
       payload JSONB NOT NULL,
       status VARCHAR(20) DEFAULT 'pending',
       created_at TIMESTAMP DEFAULT NOW(),
       started_at TIMESTAMP,
       completed_at TIMESTAMP,
       error TEXT
   );
   ```
3. Refactor `PerformPhotoIndex` to:
   - Scan filesystem and enqueue one job per photo (or batch)
   - Return immediately with a job ID
4. Create standalone Go binaries (or Docker images) for each job type:
   - `cmd/indexer/main.go` — processes index jobs
   - `cmd/thumbnailer/main.go` — processes thumbnail jobs
   - `cmd/ai-worker/main.go` — processes AI analysis jobs
5. Each worker:
   - Polls the job queue
   - Updates job status to `running`
   - Processes the job
   - Updates job status to `completed` or `failed`
6. Deploy workers as:
   - AWS Lambda (with container images)
   - Google Cloud Run jobs
   - Self-hosted Knative / OpenFaaS
   - Or simply separate systemd services / Docker containers
7. Add a job status API endpoint for the frontend to poll progress

---

### Option B: Full Serverless API
Replace the entire Gin API with AWS Lambda + API Gateway (or equivalent).

**Pros:**
- Zero server management
- Automatic scaling to zero (cost savings for personal use)
- Built-in HTTPS and DDoS protection

**Cons:**
- Cold start latency on first request
- Lambda execution limits (15 min timeout, 10 GB memory)
- Vendor lock-in
- Harder to run locally for development

**Implementation Steps:**
1. Refactor `main.go` to use AWS Lambda adapter:
   ```go
   import "github.com/awslabs/aws-lambda-go-api-proxy/gin"
   ```
   Wrap the existing Gin router in a Lambda handler.
2. Move database connection to a Lambda extension or RDS Proxy (to avoid connection pool exhaustion)
3. Move file storage to S3 (required — Lambda has ephemeral storage only)
4. Deploy via Terraform / AWS SAM / Serverless Framework
5. Use CloudFront in front of API Gateway for caching

---

### Option C: Hybrid — API Container + Serverless Workers + Event Bus
Best of both worlds: API runs as a container, workers are serverless, communication via events.

**Pros:**
- API is fast and predictable
- Workers scale independently and cheaply
- Can mix cloud and local workers

**Cons:**
- Most complex setup
- Need to operate an event bus

**Implementation Steps:**
1. Implement Option A (serverless background jobs)
2. Keep the API as-is (Gin + container)
3. Use NATS or RabbitMQ as the event bus
4. Workers can run anywhere (cloud Lambda, local Docker, Kubernetes)
5. API publishes events; workers consume them

---

### Recommended: Option A (Serverless Background Jobs)
**Rationale:** Photobox's API is lightweight (Gin + PostgreSQL). The heavy lifting is indexing, thumbnail generation, and AI analysis. Moving just those to serverless functions gives the biggest bang for buck without sacrificing API performance or adding cold-start latency to user-facing endpoints. The main API can stay as a simple container on a small VM.

Start by extracting thumbnail generation into a standalone worker. It's already decoupled from the request path (on-demand generation). Make it a background job triggered by a message queue.

**Files to Reference:**
| File | Change |
|------|--------|
| `api/internal/core/domain/job.go` | Expand job model with `type`, `payload`, `status` |
| `api/internal/adapter/storage/database/repository/jobs.go` | Job queue operations |
| `api/internal/core/service/job.go` | Enqueue/dequeue logic |
| `api/cmd/indexer/main.go` (new) | Standalone index worker |
| `api/cmd/thumbnailer/main.go` (new) | Standalone thumbnail worker |
| `api/cmd/ai-worker/main.go` (new) | Standalone AI worker |
| `api/internal/components/scheduler.go` | Publish jobs instead of running inline |

---

---

## Summary & Prioritisation

| Idea | Recommended Option | Effort | Impact | Priority |
|------|-------------------|--------|--------|----------|
| **Multi-tenancy** | A — Row-Level Security | Medium | High (enables SaaS) | P1 |
| **AI Integration** | A — Local Ollama | Medium | High (differentiating feature) | P0 |
| **Cloud Storage** | A — S3 Primary + Local Cache | Medium-High | High (scalability) | P1 |
| **Serverless** | A — Serverless Background Jobs | High | Medium (cost/ops) | P2 |

### Suggested Roadmap

**Phase 1: AI Integration (Quick Win)**
- Set up Ollama locally
- Implement `AIService` with `moondream`
- Auto-tag new photos on index
- Background job to tag existing library

**Phase 2: Multi-tenancy (Enable Sharing)**
- Add `tenant_id` to all tables
- Update all queries with tenant filter
- Add tenant middleware and JWT claims

**Phase 3: Cloud Storage (Scale Beyond Disk)**
- Implement S3 repository with MinIO
- Upload photos/thumbnails to S3 on save
- Serve binaries via presigned URLs
- Add local LRU cache

**Phase 4: Serverless Jobs (Optimise Costs)**
- Extract thumbnail generation to standalone worker
- Add message queue (NATS or PostgreSQL)
- Deploy workers as serverless functions
- Add job status API for frontend progress

---

*Generated from `plans/ideas.md` — this document supersedes the ideas list with actionable options and implementation steps.*
