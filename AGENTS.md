# Photobox — Agent Operating Instructions

These instructions apply to every session in this repository (Photobox: Go API + React webapp, deployed via Helm/ArgoCD to K3s).

For the full autonomous-implementation framework, read:
- `plans/AUTONOMOUS_IMPLEMENTATION_PLAN.md` — validation framework, per-issue recipes, workflow
- `plans/IMPLEMENTATION_DECISIONS.md` — access decisions, report templates, open questions

## Branch Flow (enforced by GitHub rulesets — never attempt to bypass)

```
feature/*  ──PR──▶  develop  ──(Rich)──▶  main
 (agent)           (agent merges)      (Rich, batched stable releases)
```

- Work on feature branches named `feat/issue-NNN-description` or `fix/issue-NNN-description`.
- Open PRs targeting **`develop`**. The `Protect develop` ruleset blocks direct pushes — all changes land via pull request.
- Merge your own PRs into `develop` after validation passes (no review required by the ruleset; Rich may still review if a checkpoint is wanted).
- **NEVER push to `main` and never merge a PR into `main`.** The `Protect main` ruleset grants only Rich (`r1chjames`) bypass. Rich merges `develop` → `main` in batches for stable releases.
- Never disable, weaken, or request exceptions to the rulesets.

## Scope & Access

- **GitHub (this repo)**: full access — branches, PRs, issues.
- **GitLab**: `r1chjames/kubernetes-helm-charts` — Helm chart changes.
- **K3s**: cluster credentials are in scope — use for Layer 4 deployment validation.
- **Secrets**: never commit tokens, kubeconfigs, or credentials; never echo secrets into logs, issues, or PRs.

## Implementing an Issue

1. Read the full plan docs (above) and the relevant per-issue validation recipe.
2. Create branch → implement → validate locally (see Validation).
3. Open a PR against `develop` with the three required sections in the description: **ADR**, **Test Plan**, **Test Outcome** (templates in `plans/IMPLEMENTATION_DECISIONS.md` §5).
4. Merge into `develop` once validation passes.
5. Update the issue status as work progresses (Open → In Progress → PR Open → Merged) and reference the issue with `Closes #NNN`.

## Validation (required before merging to develop)

- **Layer 1 — always**:
  - `cd api && golangci-lint run ./...`
  - `cd api && go test -v -race -count=1 ./...`
  - `cd webapp && npm run lint && npm test -- --coverage && npm run build`
- **Layer 2 — backend/API changes**: `docker compose -f api/test/docker-compose.test.yml up -d postgres valkey`, then `cd api && go test -v -tags=integration ./test/integration/...`
- **Layer 3 — UI changes / new features**: `docker compose -f docker-compose.local.yml up --build -d`, then verify in the browser with the chrome-devtools MCP (see Testing Changes below)
- **Layer 4 — Helm / Docker / deploy changes**: `helm lint` + `helm template`, then K3s verification if deploying (`kubectl -n media get pods -l app=photobox`)

Include concrete evidence in the PR: test output, coverage numbers, Storybook stories, and screenshots of the feature working (see Testing Changes).

## Testing Changes (UI)

Required before merging any PR that touches UI:

- **Storybook for new/changed React components** — every new React component must ship with a co-located story (`Component.stories.tsx` under `webapp/src/`; `.storybook/main.ts` auto-discovers `**/*.stories.tsx`). Cover the main states: loading, empty, populated, error. When modifying a component, update its story.
- **Storybook build must pass**: `cd webapp && npm run build-storybook`
- **Browser verification via chrome-devtools MCP** — with the app running (Layer 3), use the chrome-devtools MCP tools to:
  - navigate to the feature and exercise the user flow end-to-end (upload → process → view where relevant)
  - check for console errors via `list_console_messages`
  - confirm network requests return the expected responses (`list_network_requests`)
- **Screenshots as evidence** — capture screenshots of the feature working (chrome-devtools `take_screenshot`, or a Storybook snapshot) and attach them to the PR. A UI PR must not merge without a visual record.

## Commit Conventions

- Conventional commits referencing the issue, e.g. `feat(issue-92): add virtual scrolling to photo grid`
- Close issues with `Closes #NNN` in the PR body or final commit

## Session Reporting

At the end of each session, post a report to the issue/PR: completed items, validation results, blockers, and next steps.
