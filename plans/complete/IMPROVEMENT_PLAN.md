# Photobox Webapp Improvement Plan

## Overview

A comprehensive audit of the React webapp based on modern React conventions, performance analysis, bug identification, and code quality assessment. Items are prioritised by impact and risk.

---

## 1. Critical: TypeScript & Build Configuration

### 1.1 Entry point is `.jsx`, should be `.tsx`
- **File:** `src/index.jsx` → rename to `src/index.tsx`
- **Impact:** TypeScript is not checking the app entry point at all. No type safety on the root render, provider composition, or `globalThis` access.
- **Fix:** Rename file, remove `import React from 'react'` (see 1.2), update `index.html` `<script>` src reference.

### 1.2 Legacy JSX transform
- **File:** `tsconfig.json` has `"jsx": "react"`
- **Impact:** Requires `import React from 'react'` in every JSX file (~16 files). React 17+ supports automatic runtime which is faster and tree-shakes better.
- **Fix:** Change to `"jsx": "react-jsx"` and remove all `import React from 'react'` statements.

### 1.3 Legacy module resolution
- **File:** `tsconfig.json` has `"moduleResolution": "node"`
- **Impact:** The `node` resolution mode is deprecated for bundler-based projects (Vite). Can cause incorrect type resolution.
- **Fix:** Change to `"moduleResolution": "bundler"` (also add `"module": "preserve"`).

### 1.4 Missing `tsconfig.node.json`
- **Impact:** Vite config files run in Node but inherit the app tsconfig, risking type errors.
- **Fix:** Create `tsconfig.node.json` for `vite.config.ts` and `vitest.config.ts`.

### 1.5 Missing TypeScript strictness flags
- **Missing:** `noUnusedParameters`, `noUncheckedIndexedAccess`
- **Fix:** Add `"noUnusedParameters": true`, consider `"noUncheckedIndexedAccess": true`.

---

## 2. Critical: Dependency Management

### 2.1 Production/Dev dependency misclassification
Packages in `dependencies` that belong in `devDependencies`:
- `@rollup/plugin-babel`, `rollup`, all `rollup-plugin-*` — these are **unused entirely** (Vite bundles with its own internal Rollup)
- `@testing-library/*` packages
- `@types/react-router-dom` (v5 types for v7 router — wrong version)
- `@vitejs/plugin-react`, `@vitejs/plugin-react-swc`

### 2.2 Unused packages to remove
- `rollup` + all `rollup-plugin-*` — leftover from pre-Vite migration
- `react-icons` — project uses `@tabler/icons-react` exclusively
- `strict-event-emitter-types` — no usage found
- `prop-types` — unnecessary in TypeScript
- `regenerator-runtime` — not needed with modern Vite
- `acorn` — no direct usage

### 2.3 Wrong `@types/react-router-dom`
- **Current:** `^5.3.3` types, **Actual:** React Router `^7.2.0`
- **Fix:** Remove entirely (React Router v7 ships its own types).

---

## 3. High: State Management & Data Fetching

### 3.1 `useEffect` for data fetching (anti-pattern)
- **File:** `src/Components/AlbumCard/useAlbumCard.ts`
- **Issue:** Uses `useEffect` + `useState` for fetching thumbnail and photo count. No caching, no deduplication, no retry, no loading/error states. Inconsistent with the rest of the app which uses TanStack Query.
- **Fix:** Replace with two `useQuery` calls.

### 3.2 Direct state mutation
- **File:** `src/Components/SettingsView/SettingsView.tsx:44`
- **Issue:** `settings[i].value = event.target.value` mutates state directly. React requires immutable updates.
- **Fix:** `setSettings(prev => prev.map((s, idx) => idx === i ? { ...s, value } : s))`

### 3.3 `useState(Boolean)` instead of `useState(false)`
- **File:** `src/Components/AlbumGrid/AlbumGrid.tsx:24`
- **Issue:** Passes the `Boolean` constructor function as initial value, not `false`.
- **Fix:** Change to `useState(false)`.

### 3.4 Missing `staleTime` defaults on QueryClient
- **File:** `src/index.jsx:12`
- **Issue:** Default `staleTime: 0` means every mount triggers refetches. Causes unnecessary network requests during navigation.
- **Fix:** Set `defaultOptions: { queries: { staleTime: 30_000 } }`.

---

## 4. High: Model Layer Architecture

### 4.1 Class-based models with private fields and getters
- **Files:** `src/Models/Album.ts`, `Photo.ts`, `User.ts`, `Setting.ts`
- **Issues:**
  - Class instances with identical data aren't referentially equal (`!==`), breaking `React.memo`, `useMemo` deps, and React Query's structural sharing
  - Can't be serialised to JSON for caching (getters are lost)
  - `Photo.metadata` type is wrong: declared as `Record<string, any>[]` but code treats it as a single `Record<string, any>`
- **Fix:** Replace all models with plain TypeScript `interface` / `type` declarations. Zero-runtime-cost, massive improvement to React compatibility.

### 4.2 Legacy tslint comments
- **Files:** All models have `// tslint:disable:variable-name`
- **Fix:** Remove — tslint is deprecated and not configured.

---

## 5. High: Performance Optimisations

### 5.1 Missing `useCallback` / `useMemo`
- **Key targets:** `PhotoGrid.tsx` (`onImageClick`, `handlePreviousPhoto`, `handleNextPhoto`), `SettingsView.tsx` (`handleValueChange`, `handleModalSave`), `LoginCard.tsx` (`handleSubmit`), `PhotoCard.tsx` (`fetchImage`, `fetchAlbumName`)
- **Fix:** Wrap event handlers in `useCallback` and derived data in `useMemo`.

### 5.2 Inline function components cause full re-renders
- **File:** `SettingsView.tsx` — `tableRow()`, `addCancelButton()`, `editingButton()` defined as methods inside the component
- **Fix:** Extract into `React.memo`-wrapped sub-components.

### 5.3 `photosRef.current = photos` during render
- **File:** `src/Components/PhotoGrid/PhotoGrid.tsx:66-67`
- **Issue:** Writing to ref during render body can cause issues with Concurrent Mode.
- **Fix:** Move to `useEffect`.

### 5.4 Unnecessary array operations on every render
- **Files:** `AlbumGrid.tsx:20`, `PhotoGrid.tsx:21` — `maxDisplayed: 20000000` (sentinel value calling `.slice(0, 20000000)` on every render)
- **Fix:** Remove the `maxDisplayed`/slice pattern entirely, or use `undefined` default with conditional slicing.

---

## 6. Medium: Error Handling & Resilience

### 6.1 No error handling in login
- **File:** `src/Components/LoginCard/LoginCard.tsx:86-95`
- **Issue:** `handleSubmit` has no try/catch. Failed login throws unhandled rejection.
- **Fix:** Wrap in try/catch, show error notification to user.

### 6.2 Floating promise in file upload
- **File:** `src/Components/CreateAlbumView/CreateAlbumView.tsx:49`
- **Issue:** `uploadPhotoToApi(...).then(...)` — upload failure is unhandled.
- **Fix:** Use `try/catch` around the await or add `.catch()`.

### 6.3 No error boundary at route level
- **File:** `src/Routing/Router.tsx`
- **Issue:** Single route error tears down entire app including navigation.
- **Fix:** Add per-route error boundaries.

### 6.4 Missing loading states for mutations
- **Files:** `LoginCard`, `CreateAlbumView`, `SettingsView`
- **Issue:** No loading indicators during async operations. Users can double-submit.
- **Fix:** Add `isLoading` state to disable buttons during submission.

### 6.5 `InfoSnackbar` dead code and ignores `show` prop
- **File:** `src/Components/Snackbar/InfoSnackbar.tsx`
- **Issue:** `show` prop accepted but never checked. Commented-out MUI code remains. Component always renders.
- **Fix:** Gate rendering on `{props.show && <Notification ... />}`.

### 6.6 Blob URL lifecycle issues
- **File:** `src/Components/PhotoDetail/PhotoDetail.tsx:38-42`
- **Issue:** Cleanup revokes blob URL on dependency change. If query refetches, displayed image's URL gets revoked while still rendered.
- **Fix:** Keep ref to current blob URL, only revoke on final unmount.

---

## 7. Medium: Security

### 7.1 JWT in localStorage + direct localStorage access from HTTP adapter
- **Files:** `src/Routing/AuthContext.tsx`, `src/Adapters/RestApiAdapter.ts`
- **Issue:** localStorage tokens are XSS-vulnerable. `RestApiAdapter.handleUnauthorized()` clears localStorage directly, coupling HTTP layer to browser storage.
- **Fix:** Centralise all token management in `AuthContext`. Ideally move to httpOnly cookies (requires backend change). Add CSP header.

### 7.2 No input sanitisation
- **Files:** `LoginCard.tsx`, `SettingModal.tsx`, `CreateAlbumView.tsx`
- **Fix:** Add client-side validation (length limits, character restrictions, trimming).

---

## 8. Medium: Code Quality

### 8.1 Inconsistent styling approach
- **Current:** CSS Modules + plain CSS + Mantine inline styles mixed together
- **Recommendation:** Consolidate on Mantine's built-in styling. CSS Modules only for complex custom layouts.

### 8.2 Inconsistent export patterns
- **Observed:** Mix of named, default, and both. `usePhotoGrid` is default, `App` is both.
- **Recommendation:** Standardise on named exports for components.

### 8.3 Filename `TypeUtils.js.ts`
- **File:** `src/utils/TypeUtils.js.ts` — double extension is confusing
- **Fix:** Rename to `TypeUtils.ts`.

### 8.4 Hardcoded user info in `AppBar`
- **File:** `src/Components/AppBar/AppBar.tsx:148-150`
- **Issue:** `src={"user.image"}`, `alt={"user.name"}`, `{"Rich"}` — all placeholder strings
- **Fix:** Wire up to `AuthContext` which should store a user profile.

### 8.5 Comma-dangle in theme font family
- **File:** `src/theme.ts:4`
- **Issue:** `'BlinkMacSystemFont Segoe UI Monaco, sans-serif'` — commas are inside the string literal
- **Fix:** `'BlinkMacSystemFont, Segoe UI, Monaco, sans-serif'`

### 8.6 Commented-out code
- **Files:** `AlbumGrid.tsx` (lines 31-33, 74-82), `InfoSnackbar.tsx` (lines 16-29)
- **Fix:** Remove dead code.

### 8.7 `React.FunctionComponent` everywhere
- **Fix:** Switch to `import type { FC } from 'react'` and use `FC<Props>`.

### 8.8 `IProps` generic naming
- **Fix:** Use named props types like `AlbumCardProps`, `LoginCardProps`.

---

## 9. Medium: Testing

### 9.1 Dual test setup files
- **Files:** `src/setupTests.js` (legacy Jest) and `src/test/setup.ts` (Vitest)
- **Fix:** Remove `src/setupTests.js`.

### 9.2 Test utilities don't support adapter injection
- **File:** `src/test/test-utils.tsx`
- **Issue:** `AllTheProviders` doesn't include `AdapterProvider`. Components using `useAdapters()` will fail in tests.
- **Fix:** Add `AdapterProvider` with mock adapters.

### 9.3 Low test coverage
- **Missing tests for:** `AuthContext`, custom hooks (`usePhotoGrid`, `useAlbumGrid`, `useAlbumCard`), `ErrorBoundary`, routing integration
- **Recommendation:** Prioritise hook tests and AuthContext tests.

---

## 10. Low: Build & HTML

### 10.1 HTML metadata
- **File:** `index.html`
- **Issues:** Placeholder description (`"Web site"`), no Open Graph / Twitter Card tags, no canonical link
- **Fix:** Add proper metadata including OG tags for photo sharing.

### 10.2 Remove `appGlobals.ts` global mutation pattern
- **File:** `src/appGlobals.ts`
- **Issue:** `globalThis.app` mutation is an anti-pattern. Vite's `import.meta.env` provides env vars natively.
- **Fix:** Remove the file. Use `import.meta.env.VITE_API_URL` directly.

---

## Recommended Implementation Order

| # | Task | Effort | Risk |
|---|------|--------|------|
| 1 | Fix `index.jsx` → `index.tsx` | Tiny | Low |
| 2 | Fix tsconfig (JSX transform, moduleResolution) | Small | Medium |
| 3 | Reclassify dependencies + remove unused | Small | Low |
| 4 | Replace class models with interfaces | Medium | High |
| 5 | Fix `useState(Boolean)` bug | Tiny | Low |
| 6 | Fix direct state mutation in SettingsView | Tiny | Low |
| 7 | Convert `useAlbumCard` to TanStack Query | Small | Low |
| 8 | Add error handling to LoginCard + CreateAlbumView | Small | Low |
| 9 | Fix InfoSnackbar (dead code + show prop) | Tiny | Low |
| 10 | Add QueryClient staleTime defaults | Tiny | Low |
| 11 | Memoization pass | Medium | Low |
| 12 | Extract inline components | Medium | Low |
| 13 | Rename TypeUtils.js.ts | Tiny | Low |
| 14 | Fix blob URL lifecycle | Small | Medium |
| 15 | Add per-route error boundaries | Small | Low |
| 16 | Consolidate test setup | Small | Low |
| 17 | Fix hardcoded user in AppBar | Medium | Medium |
| 18 | Fix HTML metadata | Tiny | Low |
| 19 | Remove appGlobals.ts | Tiny | Low |
| 20 | Standardise exports & naming | Medium | Low |

---

*Generated from full codebase audit. All file references use paths relative to `webapp/`.*
