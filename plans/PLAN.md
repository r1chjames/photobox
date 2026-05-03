# Photobox Next Stage Plan

## Bug Fixes (High Priority)

- [x] 0. AppBar Reloads on First Visit to Lazy-Loaded Routes
  - File: `src/Routing/Router.tsx`, `src/Components/AppBar/AppBar.tsx`
  - Root cause: `<Suspense>` wrapped the entire `<Routes>` tree, so lazy-loading a page replaced AppBar with the loading fallback.
  - Fix: Restructured to use React Router nested routes with `AppShellLayout`. AppBar is now rendered as a layout route **outside** the `Suspense` boundary. Only `<Outlet />` (page content) is suspended.

- [x] 1. Registration Form Shows Errors by Default
  - File: `src/Components/LoginCard/LoginCard.tsx`
  - Removed hardcoded `error="Invalid username"` and `error="Invalid email"` on registration inputs; validate on submit/blur instead.

- [ ] 2. Accessibility — Color Contrast Failures
  - Active nav-link label/description contrast too low (`#4c6ef5` on `#edf1fe` = 3.83:1, `#868e96` on `#ffffff` = 3.32:1)
  - "Create album" / "Upload photos" button labels fail contrast (4.32:1, needs 4.5:1)
  - Use Mantine theme tokens or darken text to meet WCAG AA.

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

- [ ] 13. Drag-and-Drop Upload
  - Add a `react-dropzone` upload area to the dashboard and album pages.

- [ ] 14. Offline / Backend Unavailable Indicator
  - Add a global network-status banner that shows when requests fail.

- [ ] 15. PWA Service Worker Caches Only Static Shell
  - File: `public/service-worker.js`
  - Update SW to cache Vite build assets using `self.__WB_MANIFEST` from Workbox, or rely fully on `vite-plugin-pwa`.

- [ ] 16. Missing Apple Touch Icon / Maskable Icon
  - Add `<link rel="apple-touch-icon">` and a maskable PNG for iOS PWA support.

- [ ] 17. Search by Date Range / Facets
  - Add date-range pickers, camera model filters, and tag chips to the search UI.

- [ ] 18. Timeline Scrubber Hidden on Mobile
  - Collapse the scrubber into a bottom-sheet or year-picker on mobile.

- [ ] 19. Auto-Retry Failed Thumbnails
  - Add a retry button or automatic retry with exponential backoff for failed thumbnails.

- [ ] 20. Keyboard Shortcuts Discovery
  - Show a small "Press ? for shortcuts" toast on first visit.

## Architecture / Performance (Medium Priority)

- [ ] 21. Centralize Blob URL Lifecycle
  - Create a `useBlobUrl(adapter, photoId)` hook to standardize creation, caching, and cleanup across `PhotoCard`, `MapView`, `TrashView`, `Slideshow`.

- [ ] 22. Reduce Query Client Cache Mutation Duplication
  - Extract an `optimisticallyUpdatePhoto(id, partial)` helper to replace copy-pasted `queryClient.setQueriesData` patterns in `PhotoGrid`, `PhotoDetail`, `PhotoCard`.

- [ ] 23. Add Loading Skeletons to All Views
  - Use `<Skeleton>` consistently across `Shares`, `Users`, `Map`, `Duplicates` views.

- [ ] 24. Tests Failing from File-Descriptor Exhaustion
  - 5 test suites fail with `ENFILE: file table overflow` loading `@tabler/icons-react` ESM icons.
  - Switch Vitest config to use CJS barrel export, increase `ulimit`, or mock the icon library in tests.
