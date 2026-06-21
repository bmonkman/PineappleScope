# PineappleScope UI Modernization — Design

**Date:** 2026-06-20
**Status:** Approved (design), pending implementation plan

## Problem

The frontend is built on **Material Design Lite (MDL)** loaded from `code.getmdl.io`,
plus Roboto and Material Icons from Google Fonts. MDL is a dead Google project and its
CDN is unreliable, so the app currently renders broken/unstyled. The goal is a modern,
clean look for this small kiln-monitoring app, with no recurrence of CDN rot.

## Scope

Frontend only. Affected files:

- `resources/html/base.html`
- `resources/html/list.html`
- `resources/html/firing.html`
- `resources/html/stats.html`
- `resources/html/new-firing.html` (the edit form)
- `resources/css/styles.css`
- `resources/js/pineappleScope.js` (snackbar replacement only)

**No Go/backend changes.** Confirmed: no `.go` file references MDL/Material classes.
Routes already serve `/css`, `/js`, `/images` statically and parse the templates
(`cmd/pineapplescope/main.go:48-94`). The `{{ .version }}`, `{{ .title }}`,
`{{ .currentTemp }}`, `{{ .deviceCheckedIn }}` template vars stay as-is.

## Decisions

- **Styling:** handwritten, self-hosted CSS. No framework, no CSS CDN. (Removes the
  rot risk entirely.)
- **Theme:** neutral grays with a single warm orange accent (`#fe8b36`, already the
  chart's inner-temperature color). Light mode only.
- **Fonts:** native system font stack — no Google Fonts.
- **Icons:** small inline SVGs for the few needed (home, flame, bell, warning).
  No Material Icons font.

## Dependencies

**Remove from `base.html`:**

- `code.getmdl.io/1.3.0/material.green-orange.min.css`
- `code.getmdl.io/1.3.0/material.min.js`
- Google Fonts Roboto stylesheet(s)
- Material Icons font stylesheet(s)

**Keep (functional, working CDNs — not styling):**

- Chart.js 4.4.1
- moment.js 2.13.0
- chartjs-adapter-moment
- hammerjs 2.0.8
- chartjs-plugin-zoom 2.0.1

## Styling approach

- Self-hosted `resources/css/styles.css`, rewritten.
- CSS custom properties for colors and spacing (e.g. `--accent: #fe8b36`,
  neutral gray scale, `--radius`, spacing steps) so the theme is tweakable in one place.
- Card-based layout, responsive. The firing-detail page currently uses hardcoded
  `width:400px` and `inline-block`; replace with a flex layout that stacks on narrow
  screens.

## Per-page changes

### Header (base.html)
- Sticky top bar: home link (SVG), title `{{ .title }}`, current-temp badge
  `({{ .currentTemp }}˚C)`, and a plain "Stats" link.
- Keep the device-offline warning icon shown when `{{ .deviceCheckedIn }}` is `"0"`,
  linking to `/stats` (as a warning SVG).
- **Remove:** the non-functional search box, the `more_vert` overflow menu, and the
  commented-out drawer markup.

### List (list.html)
- Firings rendered as clean rows/cards, each linking to `/firing/{{ .ID }}`.
- Active firing (`EndDate > currentFiringThreshold`) marked with a flame SVG.
- Preserve existing data: name, cone number (∆), start date, duration.

### Firing detail (firing.html)
- Responsive info panel (start/end date, ambient temp, peak temp, cone, notes) that
  stacks on mobile.
- Chart in a card (`<canvas id="myChart">` unchanged so `renderChart` keeps working).
- Action buttons: Raw Data (toggles `#raw-data`), Ambient Temp (`addOuterData()`),
  Edit, Delete — restyled, same handlers/IDs.
- Notifications shown as a simple disclosure (keep low/high notification temp display
  and the link to `/firing/{{ .firing.ID }}/edit`).
- Keep the inline `<script>` data blocks and `addOuterData()` / `renderChart` wiring.

### Stats (stats.html)
- Modern, zebra-striped, responsive table. Same columns: Date, Kiln Temp,
  Ambient Temp, Humidity, CPU Temp, Uptime, Free Memory.

### Edit form (new-firing.html)
- Stacked labeled fields, restyled. Same form action/method, same field `name`s so the
  POST handler is unaffected.
- **Bug fix:** the form currently duplicates the Cone field — two inputs both
  `id/name="coneNumber"` (lines 13–20). Remove the duplicate; keep one Cone field.
- Keep Save (`editform.submit()`) and Delete (`deleteFiring(...)`) actions.

## JavaScript (pineappleScope.js)

- `deleteFiring(id)` currently uses `MaterialSnackbar`, which is gone with MDL JS.
  Replace with a small self-contained toast: a confirm step ("Are you sure? / Yes")
  and result messages ("Deleted firing." / "Error while trying to delete.."),
  preserving the existing DELETE-request + redirect-to-`/` behavior.
- Add minimal toast CSS to `styles.css` and a small toast helper in the JS.
- `renderChart` is unchanged.

## Out of scope

- No new features.
- No Go/backend changes.
- No dark mode.
- Not self-hosting the Chart.js stack (those CDNs work).

## Verification

- App builds and runs (`make` / existing run path); all 4 pages render styled.
- No requests to `code.getmdl.io` or Google Fonts remain.
- Delete-firing confirm/toast flow works.
- Chart renders, Ambient Temp toggle and zoom/pan still work.
- Edit form saves; only one Cone field present.
- Layout is reasonable on a narrow (mobile) viewport.
