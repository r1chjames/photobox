# Photobox Next Stage Plan

## Bug Fixes (High Priority)

- [x] 0. AppBar Reloads on First Visit to Lazy-Loaded Routes
  - File: `src/Routing/Router.tsx`, `src/Components/AppBar/AppBar.tsx`
  - Root cause: `<Suspense>` wrapped the entire `<Routes>` tree, so lazy-loading a page replaced AppBar with the loading fallback.
  - Fix: Restructured to use React Router nested routes with `AppShellLayout`. AppBar is now rendered as a layout route **outside** the `Suspense` boundary. Only `<Outlet />` (page content) is suspended.

- [x] 1. Registration Form Shows Errors by Default
  - File: `src/Components/LoginCard/LoginCard.tsx`
  - Removed hardcoded `error="Invalid username"` and `error="Invalid email"` on registration inputs; validate on submit/blur instead.

- [x] 2. Accessibility — Color Contrast Failures
  - Replaced default Mantine `indigo` palette with custom darker `indigoDark` palette (primary: `#4263eb`, >4.5:1 on white).
  - Changed breadcrumb and mobile nav inactive items from `dimmed` (`#868e96`, 3.32:1) to `gray.7` (meets WCAG AA).
  - Files: `src/theme.ts`, `src/Components/AppBar/AppBar.tsx`.

- [x] 3. Password Field Not in a `<form>`
  - Login `PasswordInput` not wrapped in `<form>`, breaking browser password managers and autofill.
  - Wrapped inputs in `<form onSubmit={handleSubmit}>`.

- [x] 4. `React.createRef` Called Every Render
  - File: `src/Components/PhotoCard/PhotoCard.tsx:38`
  - Replaced `React.createRef()` with `useRef<HTMLImageElement>(null)`.

- [x] 5. Stale Closure in Thumbnail Cleanup
  - Files: `MapView.tsx`, `TrashView.tsx`
  - `useEffect` cleanup references stale `thumbnailUrls` state, leaving orphaned blob URLs.
  - Fixed by using a ref to track latest URLs in cleanup.

- [x] 6. `forwardRef` Console Error
  - `forwardRef render functions accept exactly two parameters: props and ref.`
  - Resolved after LoginCard and other fixes (was triggered by component re-renders).

- [x] 7. Settings View Empty State Is Confusing
  - File: `src/Components/SettingsView/SettingsView.tsx`
  - Added skeleton loader while loading and an actionable empty state on failure.

## UX & Polish (High Priority)

- [x] 8. Photo Detail Page Lacks Navigation
  - File: `src/Components/PhotoDetail/PhotoDetail.tsx`
  - Added arrow-key and on-screen prev/next navigation to the standalone detail page.

- [x] 9. Slideshow Ignores Videos
  - File: `src/Components/Slideshow/Slideshow.tsx`
  - Added conditional rendering for `mediaType === 'video'` using `<video>` element.

- [x] 10. Missing "Add to Album" Bulk Action
  - File: `src/Components/PhotoGrid/PhotoGrid.tsx`
  - Added "Add to album" bulk action with album picker modal.

- [x] 11. Duplicates View Is Read-Only
  - File: `src/Components/DuplicatesView/DuplicatesView.tsx`
  - Added "Delete all but one" and "Delete all" actions per group.

- [x] 12. Mobile Bottom Nav Is Incomplete
  - File: `src/Components/AppBar/AppBar.tsx`
  - Replaced last item with "More" button that opens a bottom drawer with overflow routes.

## Feature Suggestions (Medium Priority)

- [x] 13. Drag-and-Drop Upload
  - Added `useDropzone`-based upload areas to `Dashboard.tsx` and `PhotoGrid.tsx`.
  - Supports drag-and-drop or click-to-upload for `image/*` files.
  - Uploads to "General" album via `photosAdapter.uploadPhoto()` with progress bar and notifications.

- [x] 14. Offline / Backend Unavailable Indicator
  - Added `NetworkStatusBanner` component that monitors `navigator.onLine` and pings backend every 30s.
  - Shows red Alert when offline, yellow Alert when backend unreachable. Dismissible with 5-min cooldown.
  - Integrated into `App.tsx` above the router.

- [x] 15. PWA Service Worker Caches Only Static Shell
  - File: `public/service-worker.js`
  - Updated SW to consume `self.__WB_MANIFEST` (from Workbox/vite-plugin-pwa) and precache Vite build JS/CSS chunks alongside the static shell.

- [x] 16. Missing Apple Touch Icon / Maskable Icon
  - `index.html` already had `<link rel="apple-touch-icon" href="/logo192.png" />`. Verified present.

- [x] 17. Search by Date Range / Facets
  - Added native date inputs ("From"/"To") and tag chips to `SearchView.tsx`.
  - Fetches tags via `photosAdapter.getAllTags()` and displays as clickable `Chip` components.
  - Added "Clear filters" button to reset all filters.
  - Tags are passed to `PhotoGrid` via its existing `tags` prop.

- [x] 18. Timeline Scrubber Hidden on Mobile
  - Made `TimelineScrubber` responsive: desktop keeps the side panel; mobile shows a FAB (floating action button) that opens a bottom `Drawer` with the year/month list.
  - Removed `!isMobile` guard in `PhotoGrid.tsx` — the scrubber now handles responsiveness internally.

- [x] 19. Auto-Retry Failed Thumbnails
  - Added exponential backoff (1s, then 2s) to `fetchThumbnailWithAuth` in `ThumbnailUtils.ts`.
  - Added retry button with `IconRefresh` in `GridImageItem` error state. Revokes old blob URL and re-fetches on click.

- [x] 20. Keyboard Shortcuts Discovery
  - Created `useShortcutsHint` hook that shows a one-time blue notification toast: "Press ? anytime to see keyboard shortcuts".
  - Tracks seen state in `localStorage`. Integrated into persistent `AppBar` layout.

## Architecture / Performance (Medium Priority)

- [x] 21. Centralize Blob URL Lifecycle
  - Created `useBlobUrl` hook at `src/hooks/useBlobUrl.ts` with module-level ref-counted cache.
  - Refactored `ThumbnailUtils.ts` to delegate to centralized cache. Exports imperative helpers for non-React code.

- [x] 22. Reduce Query Client Cache Mutation Duplication
  - Created `optimisticallyUpdatePhoto(queryClient, photoId, partial)` helper at `src/utils/queryClientHelpers.ts`.
  - Refactored `PhotoDetail.tsx` and `PhotoCard.tsx` to use it. `PhotoGrid.tsx` uses a bulk variant (Set of IDs) so it keeps its inline implementation.

- [x] 23. Add Loading Skeletons to All Views
  - Added `<Skeleton>` loading states to `ShareManagement`, `UserManagement`, and `MapView`.
  - Skeletons match the layout structure of loaded content (table rows, map area).

- [x] 24. Tests Failing from File-Descriptor Exhaustion
  - Fixed by adding `test.alias` in `vite.config.ts` to redirect `@tabler/icons-react` to CJS barrel export during tests.
  - Added `test.server.deps.inline` to prevent Vite from transforming the package.
  - All 15 test suites now pass (146/146 tests).
