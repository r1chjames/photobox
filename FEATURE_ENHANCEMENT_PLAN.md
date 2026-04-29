# Photobox Feature Enhancement & UX Improvement Plan

## Overview

This plan focuses on **user-facing features and UI polish** that elevate Photobox from a functional photo browser into a slick, self-hosted Google Photos competitor. It is intentionally separate from the technical improvement plan (`webapp/IMPROVEMENT_PLAN.md`) which addresses code quality and bugs.

Each item is tagged with:

| Dimension | Values |
|-----------|--------|
| **Priority** | P0 (essential), P1 (high), P2 (nice-to-have), P3 (future) |
| **Effort** | S (small, <2 days), M (medium, 2-5 days), L (large, 1-2 weeks), XL (multi-week) |
| **Impact** | High, Medium, Low — subjective UX value |

---

## 1. Navigation & Layout

### 1.1 Timeline Scrubber / Date Quick-Nav
**P0 · M · High**

The single highest-impact navigation missing. Google Photos' side date scrubber lets users jump to any month/year instantly. For large libraries, infinite scroll alone is painful.

- **UI:** Vertical date rail on the right edge of the photo grid showing year/month labels. Click/tap to jump.
- **Backend:** `GET /api/photos/timeline` returning `[{year, month, count}]` for the scrubber density display.
- **Frontend:** Scroll to first photo matching the selected month using the existing cursor-based pagination.

### 1.2 Breadcrumb Navigation in AppBar
**P1 · S · Medium**

Currently there is no indication of where you are beyond the sidebar highlight. Add breadcrumbs: `Photos > Album Name` or `Dashboard > Album > Photo`.

- Mantine's `Breadcrumbs` component already exists in the library.
- Show current album name, photo count, and back navigation.

### 1.3 Bottom Navigation Bar (Mobile)
**P1 · S · Medium**

The sidebar navigation pattern breaks down on mobile. A bottom tab bar with 4-5 primary destinations (Photos, Albums, Search, Favorites, Settings) would feel far more native.

- Use Mantine's `AppShell` with `navbar` breakpoint to swap between sidebar (desktop) and bottom tabs (mobile).
- Keep the existing `AppBar` header regardless.

### 1.4 Grid Density Control
**P1 · S · Medium**

Let users adjust thumbnail size with a slider in the grid toolbar. Three presets: compact, comfortable, large.

- Control column count in `JustifiedInfiniteGrid` via the `gap` and row height.
- Persist preference in localStorage.

### 1.5 View Mode Toggle (Grid / List / Map)
**P2 · M · Medium**

Beyond the justified grid, add:
- **List view:** Table with filename, date, size, dimensions, camera — useful for triage.
- **Map view:** See Section 8.2.

---

## 2. Photo Browsing Experience

### 2.1 Skeleton / Blurhash Loading Placeholders
**P0 · S · High**

Currently thumbnails pop in with no placeholder, causing layout shift. Generate a tiny blurhash or use a dominant-color placeholder during the indexing process.

- **Backend:** Store `blurhash` or `dominantColor` in the `Photo` model, computed during thumbnail generation.
- **Frontend:** Render a blurred `<div>` with the dominant color or `<Blurhash>` canvas as the `GridImageItem` background while the thumbnail loads.

### 2.2 Smart Keyboard Shortcuts
**P1 · S · High**

The lightbox already has arrow-key navigation. Extend to:
- `F` — toggle favorite
- `Delete` / `Backspace` — delete current photo (with confirmation)
- `E` — show EXIF / info panel
- `D` — download original
- `Esc` — close lightbox / go back
- `Cmd/Ctrl + A` — select all (in selection mode)
- `?` — show keyboard shortcut overlay
- Add a `useKeyboardShortcuts` hook with a central registry.

### 2.3 Swipe Gestures (Mobile Lightbox)
**P1 · M · Medium**

Swipe left/right to navigate between photos in the lightbox. Swipe down to dismiss. Use `@use-gesture/react` + `react-spring` or a lightweight touch handler.

### 2.4 Slideshow Mode
**P2 · M · Medium**

Full-screen, auto-advancing view with configurable interval (3s, 5s, 10s). Play/pause controls, subtle Ken Burns effect on still images.

- Build as a variant of the existing `PhotoCard` lightbox modal.
- Add a "Slideshow" button in the album/photo grid toolbar.

### 2.5 Photo Comparison / Side-by-Side View
**P3 · M · Low**

Select two photos and view them side-by-side for comparison. Useful for burst shots or similar compositions.

### 2.6 Fullscreen Immersive Mode
**P2 · S · Low**

Hide all UI chrome (header, sidebar) when in fullscreen photo view. Show on mouse move / tap. Similar to YouTube's cinema mode.

---

## 3. Photo Interaction & Management

### 3.1 Batch Selection & Bulk Operations
**P0 · L · High**

Currently there is zero multi-select capability. This is table-stakes for any photo manager.

- **Selection mode:** Long-press (mobile) or click a "Select" button to enter selection mode. Checkboxes appear on each thumbnail.
- **Select all / deselect all** toggle.
- **Bulk actions toolbar** (appears when 1+ selected):
  - Add to album / move to album
  - Delete
  - Download as zip
  - Add/remove tag
  - Favorite / unfavorite
- **Drag-to-select:** Click and drag a marquee rectangle to select a range. Use the existing `@egjs/react-infinitegrid` events to determine visible items.

### 3.2 Photo Deletion with Trash
**P1 · L · High**

Currently photos cannot be deleted at all. Needed both for the UI and for basic content management.

- **Backend:** `DELETE /api/photos/:id` (move to `.trash/` directory or flag in DB). `POST /api/photos/trash/restore/:id`. `DELETE /api/photos/trash/empty` (permanent). Auto-purge trash after 30 days via cron.
- **Frontend:** Delete button in photo detail and grid context menu. Confirmation dialog. Trash view accessible from sidebar. "Empty trash" button.

### 3.3 Photo Rotation & Flip
**P2 · M · Medium**

Rotate 90° CW/CCW, flip horizontal/vertical. Essential for photos imported from cameras that don't respect EXIF orientation.

- **Backend:** `POST /api/photos/:id/rotate?direction=cw|ccw`. Apply lossless JPEG rotation using `jpegtran`-style EXIF orientation manipulation (preferred) or re-encode with `imaging` library. Regenerate thumbnail.
- **Frontend:** Rotate/flip buttons in the lightbox toolbar and batch operations.

### 3.4 Favorites / Starring
**P1 · M · Medium**

A simple boolean flag on photos. Favorite button (star icon) in the lightbox and on hover over thumbnails. "Favorites" filter in the sidebar.

- **Backend:** `PATCH /api/photos/:id/favorite` (toggle). `GET /api/photos?favorites=true`.
- **Frontend:** Star icon overlay on `GridImageItem`. "Favorites" virtual album in sidebar.

### 3.5 Download
**P1 · S · High**

- **Single photo:** Download button in lightbox and photo detail. Uses `GET /api/photo/bin/:id` with `Content-Disposition: attachment`.
- **Album / selection:** Download as zip. Backend endpoint `POST /api/photos/download` that streams a zip archive. Show progress if archive generation takes time.

---

## 4. Search & Discovery

### 4.1 Full-Text Search
**P0 · L · High**

The single biggest gap versus Google Photos. Users need to find photos by filename, album name, or metadata.

- **Backend:** Add a PostgreSQL `tsvector` column on `photos` and `albums` (populated via trigger or in indexing). `GET /api/search?q=...` — ranks results by relevance, returns both photos and albums. Leverage existing EXIF metadata (camera model, lens, date, GPS city lookups if added).
- **Frontend:** Global search bar in the `AppBar` (replaces/clones the Mantine `Autocomplete` component). Results dropdown with sections for Photos and Albums. Full search results page with the same infinite grid filtered to search results.

### 4.2 Faceted Filtering
**P2 · L · Medium**

Once search exists, add faceted filters for narrowing results:
- **Date range:** Date picker or "Last 7 days / 30 days / this year" presets.
- **Camera / Lens:** Populated from EXIF data.
- **File type:** JPEG, PNG, HEIC, RAW, Video.
- **Location:** If GPS is present.
- **Favorite:** Starred only.

Filter chips appear below the search bar. URL query params sync so filters are shareable/bookmarkable.

### 4.3 Tag Management
**P2 · M · Medium**

The `tags` field exists on both `Album` and `Photo` models but is unused. Build a simple tagging system:
- Inline tag editor in photo detail sidebar (add/remove tags via chips).
- Tag autocomplete from existing tags.
- "Browse by tag" view — tag cloud or list in sidebar.
- Batch add/remove tags in selection mode.

---

## 5. Album Management

### 5.1 Album Editing (Update / Delete)
**P1 · M · High**

The API has `GET /api/album/:id` and `GET /api/albums` but no update or delete. Albums can only be created via directory scanning or upload.

- **Backend:** `PATCH /api/albums/:id` (rename, description, cover photo, tags). `DELETE /api/albums/:id` (delete album, keep photos or delete photos — option via query param).
- **Frontend:** Edit button on album card → inline rename or modal. Delete with confirmation. "Set as cover" on photos within the album.

### 5.2 Manual Album Creation (Empty)
**P1 · S · Medium**

Currently albums can only be created with photos (upload flow) or by directory scanning. Add a simple "New Album" button that creates an empty album, ready for photos to be added later.

### 5.3 Drag Photos to Album
**P2 · M · Medium**

In the main photo grid or selection mode, drag selected photos onto an album card in the sidebar to add them. Uses HTML5 drag-and-drop API.

### 5.4 Album Sorting & Ordering
**P2 · S · Low**

Currently albums appear in creation order. Add sort options: alphabetical, most recent photo, photo count. Allow manual reorder via drag-and-drop on the albums page.

### 5.5 Album Cover Customisation
**P2 · S · Low**

Let users pick which photo in the album is the cover. Store `cover_photo_id` on the `Album` model.

---

## 6. Sharing & Collaboration

### 6.1 Shareable Photo Links
**P1 · L · High**

Generate a unique, unguessable URL that displays a single photo or album to anyone with the link (no login required).

- **Backend:** `POST /api/share` (create share link with optional expiry and password). `GET /api/shared/:token` (public endpoint, no auth). Store in `shared_links` table with token, resource type, resource ID, expiry, password hash, view count.
- **Frontend:** "Share" button in photo lightbox and album view. Modal with copyable link, expiry picker, optional password field. "Manage shares" page listing all active shares with ability to revoke.

### 6.2 Shared Album (View-Only)
**P2 · L · Medium**

Extension of shareable links — a shared album page that:
- Shows a clean, branded gallery without the full app chrome.
- Supports the same infinite scroll, lightbox, and download.
- Shows EXIF data if the sharer opted in.
- "Save to my library" button if the viewer is logged in.

### 6.3 User Management UI
**P2 · M · Medium**

The backend has full user CRUD (`POST /api/user/register`, `POST /api/user/update`) but the frontend only has login/register. An admin-only "Users" page with:
- Table of all users (username, role, status, created date).
- Toggle approve/disable.
- Change role dropdown.
- Delete user.

---

## 7. Smart & AI Features

### 7.1 Memories / "On This Day"
**P2 · M · Medium**

Show photos taken on this day in previous years. A `Memories` card on the Dashboard or a dedicated page. Uses EXIF `DateTimeOriginal` which is already extracted.

- **Backend:** `GET /api/photos/on-this-day?month=X&day=Y`. Simple query on parsed EXIF date.
- **Frontend:** "On This Day" section on the dashboard showing a carousel of photos from 1, 2, 3+ years ago. Show notification/badge when memories are available.

### 7.2 Map View
**P1 · L · Medium**

Many EXIF records contain GPS coordinates. Plot them on an interactive map.

- **Frontend:** Use `react-leaflet` or `maplibre-gl` (both open-source, no API key needed with OpenStreetMap tiles). Cluster markers by proximity. Click a cluster to zoom in, click a marker to see the photo thumbnail.
- **Backend:** `GET /api/photos/geodata` returning `[{id, lat, lng, thumbnailUrl, dateTaken}]` for photos with GPS data.
- Filter by bounding box to only load photos in the current map viewport.

### 7.3 Duplicate Detection
**P2 · M · Medium**

During indexing or as a background job, detect duplicate photos by perceptual hash (e.g., `dHash`). Flag duplicates for manual review or auto-hide.

- **Backend:** Compute a 64-bit perceptual hash during indexing, store in `Photo.metadata`. Run a `GROUP BY` on hash to find groups of 2+ photos. `GET /api/photos/duplicates` returns groups.
- **Frontend:** "Duplicates" view showing side-by-side groups with "Keep" / "Delete" actions.

### 7.4 AI-Powered Auto-Tagging (Future)
**P3 · XL · Low**

Run a local ML model (e.g., ONNX with a small vision transformer or CLIP) during indexing to generate descriptive tags. This is a significant undertaking but is the key differentiator for a Google Photos competitor.

- Consider: `clip.cpp`, `llama.cpp` with a vision adapter, or calling an optional external API (Ollama, OpenAI-compatible).
- Store tags in the existing `tags` field.
- This item is aspirational — it belongs in a separate plan.

---

## 8. Media & Import

### 8.1 Video Support
**P1 · L · High**

Support video files in the library. At minimum: playback in the lightbox, thumbnail generation (first frame or midpoint), and basic metadata (duration, resolution, codec).

- **Backend:** Use `ffmpeg` (or a Go binding like `ffmpeg-go`) to generate video thumbnails and extract metadata during indexing. Support common formats: MP4, MOV, WebM, AVI.
- **Frontend:** HTML5 `<video>` player in the lightbox with standard controls (play, seek, volume, fullscreen). Show duration badge on video thumbnails. Filter by "Video" type.

### 8.2 HEIC / RAW Preview
**P2 · M · Low**

Many users import directly from iPhones (HEIC) or cameras (RAW). Generate JPEG previews for HEIC (`libheif`) and RAW files (`libraw` or `dcraw`) during indexing.

- **Backend:** Detect file type, generate JPEG preview, store as the thumbnail. Original file remains available for download.

### 8.3 Upload Progress & Queue
**P1 · M · Medium**

The current dropzone upload gives no progress feedback and processes files one-by-one silently. This is critical UX for any upload flow.

- **Frontend:** Upload queue panel showing each file with a progress bar, filename, and status (pending / uploading / done / error). Use `axios` `onUploadProgress`. Allow cancel/retry per file. Show aggregate progress (e.g., "3 of 12 uploaded").
- Drag-and-drop directly onto album cards to upload into that album.

### 8.4 Watch Folder / Auto-Import
**P2 · M · Low**

In addition to scheduled indexing, watch the photo directory for new files using filesystem events (`fsnotify` in Go). Instantly index new photos as they appear.

---

## 9. Mobile & Responsiveness

### 9.1 Pull-to-Refresh
**P1 · S · Medium**

Standard mobile gesture. Swipe down on the photo grid to trigger a refetch.

- Use `@tanstack/react-query`'s built-in `refetch` via a touch handler or wrap the grid in a pull-to-refresh container.

### 9.2 Responsive Photo Detail
**P1 · S · Medium**

`PhotoDetail` currently shows a sidebar with metadata. On mobile, collapse the sidebar into a bottom sheet or tabbed panel below the photo.

### 9.3 Native Share API
**P2 · S · Low**

On mobile, use the Web Share API (`navigator.share`) to share photos to other apps (Messages, WhatsApp, etc.). Falls back to copy-link on desktop.

### 9.4 PWA Support
**P2 · M · Low**

Add a service worker and web manifest so Photobox can be "installed" on mobile home screens. Enables offline caching of recently viewed thumbnails.

- Add `vite-plugin-pwa` which integrates with Workbox.
- Cache thumbnail and image bin responses for offline viewing.
- Show "Add to Home Screen" prompt.

---

## 10. UI Polish & Delight

### 10.1 Empty States
**P1 · S · Medium**

Replace blank pages with illustrated empty states:
- **No photos:** "Your library is empty. Upload photos or configure your photo directory." + CTA button.
- **No albums:** "Albums appear automatically from your photo folders, or create one manually." + CTA.
- **No search results:** "No photos match your search. Try different keywords."
- **Settings:** Already has empty state, but no "Getting started" guidance.

Use Mantine's `@mantine/notifications` or custom empty state component with an illustration/icon.

### 10.2 Toast Notifications
**P1 · S · High**

Replace `InfoSnackbar` (which is currently broken per the improvement plan) with Mantine's `@mantine/notifications`. Notifications for:
- Upload success / failure
- Deletion (with undo via trash)
- Share link copied
- Settings saved
- Indexing started / completed
- Errors (with retry action)

### 10.3 Confirmation Dialogs
**P1 · S · High**

Add Mantine `ModalsProvider` confirmation modals for destructive actions:
- Delete photo / album
- Empty trash
- Remove user
- Revoke share link

Consistent pattern: red "Delete" button, cancel option, clear description of consequences. Prevent accidental data loss.

### 10.4 Smooth Transitions & Micro-Animations
**P2 · M · Medium**

- **Page transitions:** Subtle fade between routes using `framer-motion` or CSS transitions on the `<Outlet>`.
- **Lightbox enter/exit:** Scale-up animation from the thumbnail position (shared element transition — complex but impactful).
- **Hover effects:** Subtle scale + shadow on album/photo cards.
- **Loading shimmer:** Skeleton shimmer animation.

### 10.5 Inline Rename (Albums & Photos)
**P2 · S · Low**

Click-to-edit filename or album name directly in the grid/header. No modal required.

### 10.6 Photo Info Overlay (Hover)
**P2 · S · Low**

Show a subtle info overlay on hover: date taken, camera model, favorite star. Avoids needing to open the lightbox for basic info.

---

## 11. Accessibility

### 11.1 Keyboard Navigation Parity
**P1 · M · Medium**

Ensure every interactive element is reachable via keyboard:
- Tab through photos in the grid, Enter to open.
- Tab through sidebar navigation.
- Focus trap in modals (lightbox, settings, confirmations).
- Visible focus indicators (Mantine provides these by default).

### 11.2 ARIA Labels & Roles
**P1 · M · Medium**

Audit and add:
- `aria-label` on icon-only buttons (delete, share, favorite, etc.).
- `role="img"` and `alt` text for photo thumbnails (use filename, date, or AI-generated caption).
- `aria-live` regions for status updates (upload progress, indexing status).
- Proper `role="dialog"` and `aria-modal` on lightbox and modals.

### 11.3 Screen Reader Announcements
**P2 · S · Low**

Use `aria-live="polite"` to announce:
- "X photos uploaded successfully"
- "Photo deleted"
- "Navigated to album name"

---

## 12. Performance at Scale

### 12.1 Virtualised Grid (for 100k+ photo libraries)
**P2 · L · Medium**

The current `JustifiedInfiniteGrid` renders all fetched items in the DOM. For libraries exceeding ~10k photos, DOM node count becomes a bottleneck.

- Evaluate `@tanstack/react-virtual` or `react-virtuoso` — though grid virtualization is harder than list virtualization.
- Alternatively, increase the page size and use `content-visibility: auto` CSS to reduce rendering cost of off-screen tiles.

### 12.2 Progressive Image Loading
**P2 · S · Low**

For the lightbox, load a low-res thumbnail first, then swap to the full-resolution image. Feels instant.

- Already have thumbnails (600x600). Load thumbnail as `<img>` with a blur filter, then load full image behind it, crossfade when loaded.

### 12.3 Prefetch on Hover / Idle
**P3 · S · Low**

Use TanStack Query's `queryClient.prefetchQuery` to preload neighboring photos when a user hovers over a thumbnail or idles on a photo. Makes lightbox navigation feel instant.

---

## Implementation Roadmap

### Phase 1 — Core UX Polish (Weeks 1-3)

| # | Item | Section | Effort |
|---|------|---------|--------|
| 1 | Toast notifications (replace InfoSnackbar) | 10.2 | S |
| 2 | Confirmation dialogs for destructive actions | 10.3 | S |
| 3 | Skeleton/blurhash loading placeholders | 2.1 | S |
| 4 | Empty state illustrations | 10.1 | S |
| 5 | Smart keyboard shortcuts | 2.2 | S |
| 6 | Download (single + album zip) | 3.5 | S |
| 7 | Breadcrumb navigation | 1.2 | S |

### Phase 2 — Interaction & Management (Weeks 4-7)

| # | Item | Section | Effort |
|---|------|---------|--------|
| 8 | Batch selection + bulk operations toolbar | 3.1 | L |
| 9 | Photo deletion with trash | 3.2 | L |
| 10 | Favorites / starring | 3.4 | M |
| 11 | Upload progress & queue | 8.3 | M |
| 12 | Album editing (update / delete / cover) | 5.1 | M |
| 13 | Manual empty album creation | 5.2 | S |

### Phase 3 — Search & Discovery (Weeks 8-11)

| # | Item | Section | Effort |
|---|------|---------|--------|
| 14 | Full-text search (backend + frontend) | 4.1 | L |
| 15 | Timeline scrubber / date quick-nav | 1.1 | M |
| 16 | Tag management | 4.3 | M |
| 17 | Faceted filtering | 4.2 | L |

### Phase 4 — Smart Features & Sharing (Weeks 12-16)

| # | Item | Section | Effort |
|---|------|---------|--------|
| 18 | Shareable photo/album links | 6.1 | L |
| 19 | Map view | 7.2 | L |
| 20 | Memories / "On This Day" | 7.1 | M |
| 21 | Video support (playback + thumbnails) | 8.1 | L |
| 22 | Duplicate detection | 7.3 | M |

### Phase 5 — Polish & Platform (Weeks 17-20)

| # | Item | Section | Effort |
|---|------|---------|--------|
| 23 | Mobile bottom nav + pull-to-refresh | 1.3, 9.1 | S |
| 24 | Swipe gestures (mobile lightbox) | 2.3 | M |
| 25 | Slideshow mode | 2.4 | M |
| 26 | PWA support | 9.4 | M |
| 27 | Accessibility audit (keyboard + ARIA) | 11.1-11.2 | M |
| 28 | HEIC/RAW preview generation | 8.2 | M |
| 29 | Photo rotation & flip | 3.3 | M |
| 30 | User management UI | 6.3 | M |

### Phase 6 — Delight & Future (Ongoing)

| Item | Section | Effort |
|------|---------|--------|
| Smooth page transitions & micro-animations | 10.4 | M |
| Grid density control & view mode toggle | 1.4-1.5 | M |
| Drag photos to album | 5.3 | M |
| Inline rename | 10.5 | S |
| Photo info hover overlay | 10.6 | S |
| AI auto-tagging | 7.4 | XL |
| Virtualised grid for 100k+ libraries | 12.1 | L |
| Shared album view-only page | 6.2 | L |
| Progressive image loading in lightbox | 12.2 | S |

---

## Design Principles (for all new work)

1. **Information density scales with screen size.** Desktop shows more columns, mobile shows larger thumbnails.
2. **Every action has an undo path.** Deletions go to trash, settings changes are reversible, dialog confirms destruction.
3. **Keyboard, touch, and mouse all work equally well.** No feature should require a specific input method.
4. **Performance is a feature.** Large libraries (50k+ photos) must feel as fast as small ones.
5. **Progressive disclosure.** Show the most common actions prominently, tuck advanced options behind a "..." menu.
6. **The filesystem is the source of truth.** Albums mirror directories. Deletions move to trash, not purge. The user's disk is never touched without asking.
