# Photobox Remaining Implementation Plan

## Overview

This document contains the remaining phases to be implemented after Phase 18 (Video Support). Phases 1–18 have been completed and committed to the `feature/react-query-a` branch.

Each phase is designed to be completed in sequence. Between each phase: run tests (`go test ./...` in `/api`, `npx vitest run --pool=forks` for affected frontend tests), ensure `vite build` and `npx tsc --noEmit` pass, then commit with the message format: `Phase N: <summary>`.

---

## Phase 19 — Faceted Filtering (4.2)
**Priority:** P2 · **Effort:** L · **Impact:** Medium

Extend search and photo grids with structured filter chips for narrowing results by date range, camera/lens, file type, location, and favorites.

### Tasks
- **Backend:**
  - Extend `ListPhotos` and search queries to accept filter params: `startDate`, `endDate`, `camera`, `lens`, `fileType`, `hasGPS`, `favorite`.
  - Add indexes on frequently filtered columns (camera, file type, favorite + created_epoch).
  - Add `GET /api/photos/filters` endpoint returning distinct values for cameras, lenses, and file types to populate filter dropdowns.
- **Frontend:**
  - Create `FilterBar` component with filter chips below the search bar on PhotoGrid.
  - Support date-range picker with presets (Last 7 days, Last 30 days, This year).
  - Camera/lens dropdowns populated from `/api/photos/filters`.
  - File type toggle (Image, Video, HEIC, RAW).
  - Favorite-only toggle.
  - Sync active filters to URL query params so filtered views are shareable/bookmarkable.
  - Update `usePhotoGrid` to pass filter state to adapter methods.
- **Adapters:**
  - Add filter params to `getAllPhotosInfo`, `getPhotosInfoInAlbum`, `searchPhotos`.
  - Add `getFilterOptions()` method.

### Acceptance Criteria
- [ ] User can filter photos by date range, camera, file type, and favorite status.
- [ ] Filtered URLs can be shared and reproduce the same results.
- [ ] Filter chips are removable individually and all at once.
- [ ] Filter state is reflected in the URL query string.
- [ ] Go tests pass; frontend tests pass; `vite build` passes.

---

## Phase 20 — Inline Rename for Photos & Albums (10.5)
**Priority:** P2 · **Effort:** S · **Impact:** Low

Click-to-edit filenames and album names directly in the grid or header without opening a modal.

### Tasks
- **Backend:**
  - `PATCH /api/photos/:id/rename` — update photo `name` field and optionally rename the underlying filesystem file.
  - Ensure album rename (`PATCH /api/albums/:id`) also renames the directory on disk (already exists; verify and extend if needed).
- **Frontend:**
  - Create `InlineEditableText` reusable component: shows text normally, switches to an input on click, commits on Enter/Blur, cancels on Escape.
  - Use in `PhotoGrid` list view for photo names.
  - Use in album detail header for album name.
  - Add optimistic UI update with error revert.
- **Adapters:**
  - Add `renamePhoto(photoId, newName)` to `IPhotosAdapter`.

### Acceptance Criteria
- [ ] User can click a photo filename in list view to edit it inline.
- [ ] User can click an album title to rename it inline.
- [ ] Changes persist after refresh.
- [ ] Go tests pass; frontend tests pass; `vite build` passes.

---

## Phase 21 — HEIC / RAW Preview Generation (8.2)
**Priority:** P2 · **Effort:** M · **Impact:** Low

Generate JPEG previews for HEIC and RAW camera files during indexing so they display correctly in the grid and lightbox.

### Tasks
- **Backend:**
  - Detect HEIC (`.heic`, `.heif`) and RAW (`.dng`, `.cr2`, `.nef`, `.arw`, `.orf`, `.raf`) files during filesystem scan.
  - Use external tools to generate JPEG previews:
    - HEIC → `libheif-examples` (`heif-convert`) or `ffmpeg` with libheif.
    - RAW → `dcraw` or `libraw` CLI (`rawtherapee-cli`, `darktable-cli`).
  - Store generated JPEG preview as the `Thumbnail` blob (or alongside the original) and set `ThumbnailUrl` accordingly.
  - Original file remains untouched for download.
- **Frontend:**
  - No frontend changes required if backend serves generated JPEG via existing thumbnail/bin endpoints.
  - Optionally show a "RAW" or "HEIC" badge on thumbnails.
- **Build / Dev:**
  - Document HEIC/RAW dependencies in README (e.g., `brew install libheif dcraw`).

### Acceptance Criteria
- [ ] HEIC files from iPhones display thumbnails and open in the lightbox.
- [ ] RAW files from DSLRs display thumbnails and open in the lightbox.
- [ ] Original file is preserved for download.
- [ ] Go tests pass (mock filesystem for HEIC/RAW paths); frontend tests pass; `vite build` passes.

---

## Phase 22 — Drag Photos to Album (5.3)
**Priority:** P2 · **Effort:** M · **Impact:** Medium

In selection mode, drag selected photos onto an album card in the sidebar to add them.

### Tasks
- **Backend:**
  - `POST /api/photos/move` or `POST /api/photos/album` accepting `{ photoIds[], albumId }`.
  - Update `album_id` on photos and physically move files into the target album directory.
- **Frontend:**
  - Make `GridImageItem` draggable in selection mode (`draggable="true"`).
  - Make sidebar album cards drop targets with visual feedback (highlight on drag over).
  - Handle drop event: call bulk move adapter method.
  - Show toast notification: "Moved N photos to Album Name".
- **Adapters:**
  - Add `movePhotosToAlbum(photoIds, albumId)` to `IPhotosAdapter`.

### Acceptance Criteria
- [ ] User can select multiple photos and drag them onto a sidebar album.
- [ ] Album card highlights when a valid drag is hovering over it.
- [ ] Photos are removed from current grid and moved to target album.
- [ ] Go tests pass; frontend tests pass; `vite build` passes.

---

## Phase 23 — Virtualised Grid (12.1)
**Priority:** P2 · **Effort:** L · **Impact:** Medium

Replace or augment `JustifiedInfiniteGrid` with a virtualised renderer so libraries with 10k+ photos don't overload the DOM.

### Tasks
- **Frontend:**
  - Evaluate `@tanstack/react-virtual` grid virtualization or `react-virtuoso` grid mode.
  - Alternatively, keep `JustifiedInfiniteGrid` but add `content-visibility: auto` CSS to off-screen rows.
  - Preserve infinite-scroll pagination (page size may increase to reduce refetch frequency).
  - Preserve density control and view-mode toggle behaviour.
- **Backend:**
  - No backend changes required.
- **Performance:**
  - Benchmark DOM node count before/after with 5,000+ photo dataset.
  - Ensure scroll performance stays at 60fps.

### Acceptance Criteria
- [ ] Grid remains performant with 10,000+ photos loaded (smooth scrolling, no jank).
- [ ] DOM node count is bounded and does not grow linearly with total photos.
- [ ] Infinite scroll, density control, and selection mode still work.
- [ ] Frontend tests pass; `vite build` passes.

---

## Phase 24 — AI Auto-Tagging (7.4)
**Priority:** P3 · **Effort:** XL · **Impact:** Low

Run a local ML model during indexing to generate descriptive tags automatically. This is aspirational and may be split into sub-phases if undertaken.

### Research & Options
- `clip.cpp` or `llama.cpp` with a vision adapter for zero-shot classification.
- ONNX Runtime with a small vision transformer (e.g., MobileViT).
- Optional external API integration (Ollama, OpenAI-compatible) for users who opt in.

### Tasks (if pursued)
- **Backend:**
  - Integrate an ONNX or GGUF vision model.
  - During `PerformPhotoIndex`, generate tags for each photo and store in the existing `tags` field.
  - Make AI tagging optional via a settings flag.
  - Batch process to avoid blocking indexing.
- **Frontend:**
  - Show AI-generated tags in a distinct style (e.g., dotted border) in the tag editor.
  - Allow users to remove individual AI tags or disable AI tagging globally in Settings.
- **Ops:**
  - Document model download and setup.
  - Provide a fallback (no AI tags) if model is not installed.

### Acceptance Criteria
- [ ] Photos are automatically tagged with descriptive labels during indexing.
- [ ] AI tags appear alongside user tags and are removable.
- [ ] Feature is opt-in via settings.
- [ ] System works offline (local model, no cloud dependency).

---

## Notes

- **Completed phases (1–18)** are on the `feature/react-query-a` branch. Do not re-implement them.
- **TypeScript strictness:** The project has `strict: true` and `noUnusedLocals: true`. All new code must compile cleanly with `tsc --noEmit`.
- **Tests:** Fix any test mocks when new hooks/components are introduced. Use `npx vitest run --pool=forks <path>` for individual test files if full test suite hits `ENFILE` on macOS.
- **Build:** Run `vite build` before each commit.
- **Go tests:** Run `go test ./...` in `/api` before each commit.
- **Storybook:** Add `.stories.tsx` files for all new reusable components where practical.

## Original Plan Reference

Items in this plan are drawn from `FEATURE_ENHANCEMENT_PLAN.md`. The original document also contains design principles and a historical roadmap (Phases 1–6) that has since been superseded by the actual Phases 1–18 implementation.
