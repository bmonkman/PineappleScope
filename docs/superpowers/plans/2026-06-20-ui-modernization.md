# PineappleScope UI Modernization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the dead Material Design Lite (MDL) frontend with a handwritten, self-hosted modern UI so the app renders cleanly and never breaks from CDN rot again.

**Architecture:** Frontend-only change. Rewrite `resources/css/styles.css` (handwritten, no framework, CSS custom properties), rewrite the four templates to use semantic HTML + the new classes, and replace the MDL-snackbar code in `pineappleScope.js` with a small self-contained toast. No Go/backend changes — routes already serve `/css`, `/js`, `/images` and parse the templates at startup.

**Tech Stack:** Go (Gin + multitemplate, unchanged), HTML templates, handwritten CSS, vanilla JS. Chart.js/moment/hammerjs/zoom kept on their working CDNs.

## Global Constraints

- No Go/backend changes. Do not edit any `.go` file. (Verified: no `.go` file references MDL/Material classes.)
- Preserve all template variables and their exact names: `{{ .version }}`, `{{ .title }}`, `{{ .currentTemp }}`, `{{ .deviceCheckedIn }}`, `{{ .firings }}`, `{{ .firing }}`, `{{ .currentFiringThreshold }}`, `{{ .temperatureReadings }}`, `{{ .peakTemperature }}`, `{{ .stats }}`.
- Preserve all form `name=` attributes, form `action`/`method`, and JS hook IDs/handlers: `#myChart`, `#raw-data`, `#toast-container`, `addOuterData()`, `renderChart(...)`, `deleteFiring(...)`, `editform.submit()`.
- Remove all CDN dependencies on `code.getmdl.io` and Google Fonts (`fonts.googleapis.com`). Keep Chart.js, moment, chartjs-adapter-moment, hammerjs, chartjs-plugin-zoom.
- Theme: light only, neutral grays + one warm orange accent `#fe8b36`. System font stack. Inline SVG icons (no icon font).
- App runs with `make run` on `http://localhost:1111`. Templates parse at startup, so a template syntax error shows as a startup failure, not a build failure.

---

## Task 1: CSS foundation + base layout shell

**Files:**
- Modify (full rewrite): `resources/css/styles.css`
- Modify (full rewrite): `resources/html/base.html`

**Interfaces:**
- Produces (CSS classes consumed by later tasks): `.container`, `.card`, `.btn`, `.btn--accent`, `.btn--danger`, `.app-header`, `.temp-badge`, `.icon` (inline SVG sizing), `.firing-list`, `.firing-row`, `.firing-row__title`, `.firing-row__meta`, `.info-panel`, `.info-row`, `.info-row__label`, `.info-row__value`, `.data-table`, `.form-field`, `.form-field label`, `.toast`. CSS custom properties on `:root`: `--accent`, `--bg`, `--surface`, `--text`, `--muted`, `--border`, `--radius`, `--space`.
- Consumes: nothing.

- [ ] **Step 1: Write the new `styles.css`**

Replace the entire file with:

```css
:root {
  --accent: #fe8b36;
  --accent-dark: #e5781f;
  --secondary: #4caf50;
  --bg: #f5f5f4;
  --surface: #ffffff;
  --text: #1f2933;
  --muted: #6b7280;
  --border: #e2e2e0;
  --danger: #d64545;
  --radius: 8px;
  --space: 1rem;
  --maxw: 56rem;
}

* { box-sizing: border-box; }

html, body {
  margin: 0;
  padding: 0;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.5;
}

a { color: var(--accent-dark); text-decoration: none; }
a:hover { text-decoration: underline; }

.container {
  max-width: var(--maxw);
  margin: 0 auto;
  padding: var(--space);
}

/* Header */
.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: var(--accent);
  color: #fff;
  padding: 0.75rem 1rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
}
.app-header a { color: #fff; }
.app-header__title { font-size: 1.15rem; font-weight: 600; }
.app-header__spacer { flex: 1; }
.temp-badge {
  font-size: 0.9rem;
  background: rgba(255,255,255,0.2);
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
}
.app-header__link { font-weight: 500; }
.icon { width: 24px; height: 24px; vertical-align: middle; display: inline-block; fill: currentColor; }

/* Cards */
.card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  padding: var(--space);
  margin-bottom: var(--space);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  background: var(--surface);
  color: var(--text);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 0.5rem 1rem;
  font-size: 0.95rem;
  font-family: inherit;
  cursor: pointer;
  margin: 0.25rem 0.25rem 0.25rem 0;
}
.btn:hover { background: #f0f0ef; text-decoration: none; }
.btn--accent { background: var(--accent); color: #fff; border-color: var(--accent); }
.btn--accent:hover { background: var(--accent-dark); }
.btn--danger { background: var(--surface); color: var(--danger); border-color: var(--danger); }
.btn--danger:hover { background: var(--danger); color: #fff; }

/* Firing list */
.firing-list { list-style: none; margin: 0; padding: 0; }
.firing-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
}
.firing-row:hover { background: #fafafa; }
.firing-row__flame { color: var(--accent); flex-shrink: 0; }
.firing-row__body { display: flex; flex-direction: column; }
.firing-row__title { font-weight: 600; color: var(--text); }
.firing-row__meta { color: var(--muted); font-size: 0.85rem; }

/* Info panel (firing detail) */
.firing-detail { display: flex; flex-wrap: wrap; gap: var(--space); align-items: flex-start; }
.info-panel { flex: 1 1 320px; }
.info-row { display: flex; justify-content: space-between; gap: 1rem; padding: 0.4rem 0; border-bottom: 1px solid var(--border); }
.info-row:last-child { border-bottom: none; }
.info-row__label { color: var(--muted); }
.info-row__value { text-align: right; }
.chart-card { flex: 2 1 420px; }
.chart-wrap { position: relative; height: 360px; }

/* Disclosure (notifications) */
details.notify { margin-top: 0.5rem; }
details.notify summary { cursor: pointer; color: var(--accent-dark); }

/* Tables */
.data-table { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
.data-table th, .data-table td { padding: 0.5rem 0.75rem; text-align: left; }
.data-table thead th { border-bottom: 2px solid var(--border); color: var(--muted); }
.data-table tbody tr:nth-child(even) { background: #fafafa; }
.table-scroll { overflow-x: auto; }

/* Forms */
.form-field { display: flex; flex-direction: column; margin-bottom: 1rem; }
.form-field label { font-size: 0.85rem; color: var(--muted); margin-bottom: 0.25rem; }
.form-field input {
  font-size: 1rem;
  font-family: inherit;
  padding: 0.5rem;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  color: var(--text);
}
.form-field input:focus { outline: 2px solid var(--accent); border-color: var(--accent); }

/* Toast */
.toast {
  position: fixed;
  left: 50%;
  bottom: 1.5rem;
  transform: translateX(-50%) translateY(200%);
  background: #323232;
  color: #fff;
  padding: 0.75rem 1.25rem;
  border-radius: var(--radius);
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.25s ease;
  z-index: 1000;
}
.toast.toast--show { transform: translateX(-50%) translateY(0); }
.toast button {
  background: none;
  border: none;
  color: var(--accent);
  font-weight: 600;
  cursor: pointer;
  font-size: 0.95rem;
}

@media (max-width: 600px) {
  .data-table { font-size: 0.8rem; }
  .info-row__value { text-align: right; }
}
```

- [ ] **Step 2: Write the new `base.html`**

Replace the entire file with:

```html
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="mobile-web-app-capable" content="yes">

        <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.1/chart.umd.min.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/moment.js/2.13.0/moment.min.js"></script>
        <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-moment@^1"></script>
        <script src="/js/pineappleScope.js?v{{ .version }}"></script>
        <script src="https://cdn.jsdelivr.net/npm/hammerjs@2.0.8"></script>
        <script src="https://cdn.jsdelivr.net/npm/chartjs-plugin-zoom@2.0.1"></script>

        <link rel="icon" sizes="192x192" href="/images/android-desktop.png">
        <link rel="stylesheet" href="/css/styles.css">
        <link rel="icon" href="/favicon.ico" type="image/x-icon">
        <title>PineappleScope</title>
    </head>

    <body>
        <header class="app-header">
            <a href="/" title="Home" aria-label="Home">
                <svg class="icon" viewBox="0 0 24 24"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/></svg>
            </a>
            <span class="app-header__title">{{ .title }}</span>
            <span class="temp-badge">{{ .currentTemp }}˚C</span>
            <span class="app-header__spacer"></span>
            {{ if eq .deviceCheckedIn "0" }}
            <a href="/stats" title="Device offline" aria-label="Device offline">
                <svg class="icon" viewBox="0 0 24 24"><path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/></svg>
            </a>
            {{ end }}
            <a class="app-header__link" href="/stats">Stats</a>
        </header>

        <main class="container">
            {{ template "content" . }}
        </main>
    </body>
</html>
```

- [ ] **Step 3: Verify no dead dependencies remain**

Run: `grep -rn "getmdl\|fonts.googleapis\|mdl-\|MaterialSnackbar" resources/html/base.html`
Expected: no output (exit code 1). The MDL CSS/JS, Google Fonts, and all `mdl-` classes are gone from base.html.

- [ ] **Step 4: Verify Go still builds**

Run: `make build`
Expected: builds with no errors (no Go changes were made).

- [ ] **Step 5: Verify the app starts and the header renders**

Run: `make run` (in a background shell), then `curl -s http://localhost:1111/ | grep -c "app-header"`
Expected: `1` or more (the new header markup is present, templates parsed successfully). Stop the server afterward.

- [ ] **Step 6: Commit**

```bash
git add resources/css/styles.css resources/html/base.html
git commit -m "Replace MDL base layout with handwritten self-hosted CSS"
```

---

## Task 2: Firings list page

**Files:**
- Modify (full rewrite): `resources/html/list.html`

**Interfaces:**
- Consumes from Task 1: `.firing-list`, `.firing-row`, `.firing-row__flame`, `.firing-row__body`, `.firing-row__title`, `.firing-row__meta`, `.icon`.
- Consumes template data: `.firings` (each has `.ID`, `.Name`, `.ConeNumber`, `.StartDate`, `.Duration`, `.EndDate`), `.currentFiringThreshold`.

- [ ] **Step 1: Write the new `list.html`**

Replace the entire file with:

```html
{{ define "content" }}
    <ul class="firing-list">
    {{ range .firings }}
        <a href="/firing/{{ .ID }}">
            <li class="firing-row">
                {{ if gt .EndDate.Unix $.currentFiringThreshold.Unix }}
                <svg class="icon firing-row__flame" viewBox="0 0 24 24"><path d="M13.5.67s.74 2.65.74 4.8c0 2.06-1.35 3.73-3.41 3.73-2.07 0-3.63-1.67-3.63-3.73l.03-.36C5.21 7.51 4 10.62 4 14c0 4.42 3.58 8 8 8s8-3.58 8-8C20 8.61 17.41 3.8 13.5.67z"/></svg>
                {{ end }}
                <span class="firing-row__body">
                    <span class="firing-row__title">{{ .Name }} (∆{{ .ConeNumber }})</span>
                    <span class="firing-row__meta">{{ .StartDate.Format "Mon, Jan _2 3:04PM" }} ({{ .Duration }} hours)</span>
                </span>
            </li>
        </a>
    {{ end }}
    </ul>
{{ end }}
```

- [ ] **Step 2: Verify no MDL classes remain**

Run: `grep -rn "mdl-\|material-icons" resources/html/list.html`
Expected: no output (exit code 1).

- [ ] **Step 3: Verify the list page renders**

Run: with `make run` active, `curl -s http://localhost:1111/ | grep -c "firing-list"`
Expected: `1`. Also open `http://localhost:1111/` in a browser and confirm firings appear as styled rows, with a flame icon on any currently-active firing.

- [ ] **Step 4: Commit**

```bash
git add resources/html/list.html
git commit -m "Modernize firings list page"
```

---

## Task 3: Firing detail page

**Files:**
- Modify (full rewrite): `resources/html/firing.html`

**Interfaces:**
- Consumes from Task 1: `.firing-detail`, `.info-panel`, `.card`, `.info-row`, `.info-row__label`, `.info-row__value`, `details.notify`, `.chart-card`, `.chart-wrap`, `.btn`, `.btn--accent`, `.btn--danger`, `.data-table`.
- Consumes template data: `.firing` (`.StartDate`, `.EndDate`, `.StartDateAmbientTemp`, `.ConeNumber`, `.Notes`, `.ID`, `.LowNotificationTemp`, `.HighNotificationTemp`, `.LowNotificationSent`, `.HighNotificationSent`), `.peakTemperature`, `.temperatureReadings` (`.CreatedDate`, `.Inner`, `.Outer`).
- Preserves JS hooks: `#myChart`, `#raw-data`, `addOuterData()`, `renderChart(innerOnlyData)`, and the inline `innerOnlyData`/`allData` script blocks (carried over verbatim).

- [ ] **Step 1: Write the new `firing.html`**

Replace the entire file with:

```html
{{ define "content" }}
    <div class="firing-detail">
        <div class="info-panel card">
            <div class="info-row">
                <span class="info-row__label">Start / End</span>
                <span class="info-row__value">{{ .firing.StartDate.Format "Mon, Jan _2 3:04PM" }}<br/>{{ .firing.EndDate.Format "Mon, Jan _2 3:04PM" }}</span>
            </div>
            <div class="info-row">
                <span class="info-row__label">Starting ambient temp</span>
                <span class="info-row__value">{{ .firing.StartDateAmbientTemp }}˚C</span>
            </div>
            <div class="info-row">
                <span class="info-row__label">Peak temp</span>
                <span class="info-row__value">{{ .peakTemperature }}˚C</span>
            </div>
            <div class="info-row">
                <span class="info-row__label">Cone</span>
                <span class="info-row__value">{{ .firing.ConeNumber }}</span>
            </div>
            <div class="info-row">
                <span class="info-row__label">Notes</span>
                <span class="info-row__value">{{ .firing.Notes }}</span>
            </div>
            {{ if or (gt .firing.LowNotificationTemp 0.0) (gt .firing.HighNotificationTemp 0.0) }}
            <details class="notify">
                <summary>Notification settings</summary>
                <a href="/firing/{{ .firing.ID }}/edit">
                    {{ if gt .firing.LowNotificationTemp 0.0 }}<div class="info-row"><span class="info-row__label">Low temp</span><span class="info-row__value">{{ .firing.LowNotificationTemp }}{{ if .firing.LowNotificationSent }} (sent){{ end }}</span></div>{{ end }}
                    {{ if gt .firing.HighNotificationTemp 0.0 }}<div class="info-row"><span class="info-row__label">High temp</span><span class="info-row__value">{{ .firing.HighNotificationTemp }}{{ if .firing.HighNotificationSent }} (sent){{ end }}</span></div>{{ end }}
                </a>
            </details>
            {{ end }}
        </div>

        <div class="chart-card card">
            <div class="chart-wrap">
                <canvas id="myChart"></canvas>
            </div>
        </div>
    </div>

    <div>
        <button class="btn" onclick="document.getElementById('raw-data').style.display='';">Raw Data</button>
        <button class="btn" onclick="addOuterData();">Ambient Temp</button>
        <a class="btn btn--accent" href="/firing/{{ .firing.ID }}/edit">Edit</a>
        <button class="btn btn--danger" onclick="deleteFiring({{ .firing.ID }})">Delete</button>
    </div>

    <table id="raw-data" class="data-table" style="display: none;">
    {{ range .temperatureReadings }}
        <tr><td>{{ . }}</td></tr>
    {{ end }}
    </table>

    <script>
        var myChart;
        var innerOnlyData = {
                labels: [ {{ range .temperatureReadings }}moment("{{ .CreatedDate.Format "2006-01-02 15:04:05" }}"),{{ end }} ],
                datasets: [{
                    fill: false,
                    label: 'Temperature',
                    data: [{{ range .temperatureReadings }}{{ .Inner }},{{ end }}],
                    borderColor: '#fe8b36',
                    backgroundColor: '#fe8b36',
                    lineTension: 0.1,
                    yAxisID: 'inner'
                }]
            };

            var allData = {
                labels: [ {{ range .temperatureReadings }}moment("{{ .CreatedDate.Format "2006-01-02 15:04:05" }}"),{{ end }} ],
                datasets: [{
                    fill: false,
                    label: 'Inner Temperature',
                    data: [{{ range .temperatureReadings }}{{ .Inner }},{{ end }}],
                    borderColor: '#fe8b36',
                    backgroundColor: '#fe8b36',
                    lineTension: 0.1,
                    yAxisID: 'inner'
                },
                {
                    fill: false,
                    label: 'Outer Temperature',
                    data: [{{ range .temperatureReadings }}{{ .Outer }},{{ end }}],
                    borderColor: '#4caf50',
                    backgroundColor: '#4caf50',
                    lineTension: 0.1,
                    yAxisID: 'outer'
                }]
            };

        function addOuterData() {
            if (myChart.options.scales.outer.display == false) {
                myChart.data = allData
            } else {
                myChart.data = innerOnlyData
            }
            myChart.options.scales.outer.display = !myChart.options.scales.outer.display;
            myChart.update();
        }

        window.onload = renderChart(innerOnlyData);
    </script>
{{ end }}
```

- [ ] **Step 2: Verify no MDL classes remain**

Run: `grep -rn "mdl-\|material-icons\|mdl-shadow" resources/html/firing.html`
Expected: no output (exit code 1).

- [ ] **Step 3: Verify the firing detail page renders and chart works**

Run: with `make run` active, open `http://localhost:1111/firing/1` (use a valid firing ID from the list page). Confirm: info panel and chart show side-by-side on desktop, the chart renders, "Ambient Temp" toggles the second axis, "Raw Data" reveals the table, and mouse-wheel zoom/pan works on the chart.

- [ ] **Step 4: Commit**

```bash
git add resources/html/firing.html
git commit -m "Modernize firing detail page"
```

---

## Task 4: Stats page

**Files:**
- Modify (full rewrite): `resources/html/stats.html`

**Interfaces:**
- Consumes from Task 1: `.table-scroll`, `.data-table`.
- Consumes template data: `.stats` (`.CreatedDate`, `.Temperature`, `.AmbientTemperature`, `.Humidity`, `.CPUTemperature`, `.Uptime`, `.FreeMemory`).

- [ ] **Step 1: Write the new `stats.html`**

Replace the entire file with:

```html
{{ define "content" }}
    <div class="table-scroll">
    <table class="data-table">
        <thead>
            <tr>
                <th>Date</th>
                <th>Kiln Temp</th>
                <th>Ambient Temp</th>
                <th>Humidity</th>
                <th>CPU Temp</th>
                <th>Uptime</th>
                <th>Free Memory</th>
            </tr>
        </thead>
        <tbody>
        {{ range .stats }}
            <tr>
                <td>{{ .CreatedDate.Format "Mon, Jan _2 3:04:05PM" }}</td>
                <td>{{ .Temperature }}˚</td>
                <td>{{ .AmbientTemperature }}˚</td>
                <td>{{ .Humidity }}%</td>
                <td>{{ .CPUTemperature }}˚</td>
                <td>{{ .Uptime }}</td>
                <td>{{ .FreeMemory }}</td>
            </tr>
        {{ end }}
        </tbody>
    </table>
    </div>
{{ end }}
```

- [ ] **Step 2: Verify no leftover bare-table styling references**

Run: `grep -rn "statsTable\|mdl-" resources/html/stats.html`
Expected: no output (exit code 1). (The old `.statsTable` CSS rule was dropped in Task 1's rewrite.)

- [ ] **Step 3: Verify the stats page renders**

Run: with `make run` active, open `http://localhost:1111/stats`. Confirm a styled, zebra-striped, horizontally-scrollable-on-mobile table with all seven columns.

- [ ] **Step 4: Commit**

```bash
git add resources/html/stats.html
git commit -m "Modernize stats table page"
```

---

## Task 5: Edit form page (with Cone duplicate-field bug fix)

**Files:**
- Modify (full rewrite): `resources/html/new-firing.html`

**Interfaces:**
- Consumes from Task 1: `.card`, `.form-field`, `.btn`, `.btn--accent`, `.btn--danger`.
- Consumes template data: `.firing` (`.ID`, `.Name`, `.Notes`, `.ConeNumber`, `.LowNotificationTemp`, `.HighNotificationTemp`).
- Preserves: form `action="/firing/{{ .firing.ID }}"` `method="POST"` `name="editform"`, field `name=`s (`name`, `notes`, `coneNumber`, `lowNotificationTemp`, `highNotificationTemp`), `editform.submit()`, `deleteFiring(...)`.

- [ ] **Step 1: Write the new `new-firing.html`**

Replace the entire file with (note: exactly **one** Cone field — the original had it duplicated):

```html
{{ define "content" }}
    <div class="card">
        <form action="/firing/{{ .firing.ID }}" method="POST" name="editform">
            <div class="form-field">
                <label for="name">Name</label>
                <input type="text" id="name" name="name" value="{{ .firing.Name }}">
            </div>
            <div class="form-field">
                <label for="notes">Notes</label>
                <input type="text" id="notes" name="notes" value="{{ .firing.Notes }}">
            </div>
            <div class="form-field">
                <label for="coneNumber">Cone</label>
                <input type="text" id="coneNumber" name="coneNumber" value="{{ .firing.ConeNumber }}">
            </div>
            <div class="form-field">
                <label for="lowNotificationTemp">Low notification temp</label>
                <input type="text" id="lowNotificationTemp" name="lowNotificationTemp" value="{{ .firing.LowNotificationTemp }}">
            </div>
            <div class="form-field">
                <label for="highNotificationTemp">High notification temp</label>
                <input type="text" id="highNotificationTemp" name="highNotificationTemp" value="{{ .firing.HighNotificationTemp }}">
            </div>
        </form>
        <button class="btn btn--accent" onclick="editform.submit();">Save</button>
        <button class="btn btn--danger" onclick="deleteFiring({{ .firing.ID }})">Delete</button>
    </div>
{{ end }}
```

- [ ] **Step 2: Verify exactly one Cone field and no MDL classes**

Run: `grep -c 'name="coneNumber"' resources/html/new-firing.html`
Expected: `1` (the duplicate is removed).

Run: `grep -rn "mdl-" resources/html/new-firing.html`
Expected: no output (exit code 1).

- [ ] **Step 3: Verify the edit form renders and saves**

Run: with `make run` active, open `http://localhost:1111/firing/1/edit`. Confirm one of each field (single Cone), change the Name, click Save, and verify it persists (redirects / shows the updated value on the detail page).

- [ ] **Step 4: Commit**

```bash
git add resources/html/new-firing.html
git commit -m "Modernize edit form and remove duplicate Cone field"
```

---

## Task 6: Toast replacement for delete flow

**Files:**
- Modify: `resources/js/pineappleScope.js` (replace `deleteFiring`, keep `renderChart`)
- Modify: `resources/html/base.html` (add toast container element)

**Interfaces:**
- Consumes from Task 1: `.toast`, `.toast--show` CSS.
- Produces: `deleteFiring(id)` using a self-contained toast instead of `MaterialSnackbar`. Preserves the DELETE request to `/firing/<id>` and redirect to `/` on success.

- [ ] **Step 1: Add the toast container to `base.html`**

In `resources/html/base.html`, immediately before the closing `</body>` tag, add:

```html
        <div id="toast-container" class="toast" role="status" aria-live="polite">
            <span id="toast-message"></span>
            <button id="toast-action" type="button" style="display:none;"></button>
        </div>
```

- [ ] **Step 2: Replace `deleteFiring` in `pineappleScope.js`**

Replace the `deleteFiring` function (lines 1-29 of the original) with the following. Leave `renderChart` (the rest of the file) unchanged.

```javascript
function showToast(message, actionText, actionHandler, timeout) {
    var container = document.getElementById('toast-container');
    var msg = document.getElementById('toast-message');
    var action = document.getElementById('toast-action');

    msg.textContent = message;

    if (actionText) {
        action.textContent = actionText;
        action.style.display = '';
        action.onclick = function () {
            container.classList.remove('toast--show');
            actionHandler();
        };
    } else {
        action.style.display = 'none';
        action.onclick = null;
    }

    container.classList.add('toast--show');

    if (timeout) {
        setTimeout(function () { container.classList.remove('toast--show'); }, timeout);
    }
}

function deleteFiring(id) {
    var doDelete = function () {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                showToast('Deleted firing.', null, null, 2000);
                setTimeout(function () { document.location = "/"; }, 1500);
            } else if (this.readyState == 4) {
                showToast('Error while trying to delete..', null, null, 4000);
            }
        };
        xhttp.open("DELETE", "/firing/" + id, true);
        xhttp.send();
    };

    showToast('Are you sure?', 'Yes', doDelete, 5000);
}
```

- [ ] **Step 3: Verify no MDL JS references remain**

Run: `grep -rn "MaterialSnackbar\|mdl-" resources/js/pineappleScope.js`
Expected: no output (exit code 1).

- [ ] **Step 4: Verify the delete flow works end-to-end**

Run: with `make run` active, open a firing detail page, click Delete. Confirm: a toast slides up reading "Are you sure?" with a "Yes" button; clicking Yes deletes the firing, shows "Deleted firing.", and redirects to the list. (Test on a disposable firing, or reseed from `resources/test_data.sql`.)

- [ ] **Step 5: Commit**

```bash
git add resources/js/pineappleScope.js resources/html/base.html
git commit -m "Replace MDL snackbar with self-contained toast for delete flow"
```

---

## Final verification

- [ ] **No MDL or Google Fonts references anywhere in resources**

Run: `grep -rn "getmdl\|fonts.googleapis\|mdl-\|material-icons\|MaterialSnackbar" resources/`
Expected: no output (exit code 1).

- [ ] **App builds and all pages render styled**

Run: `make build` (succeeds), then `make run` and walk every page: `/`, `/firing/<id>`, `/firing/<id>/edit`, `/stats`. Each is styled, responsive when the window is narrowed, with no broken/unstyled MDL remnants.

- [ ] **Chart, toggles, zoom, edit-save, and delete-toast all function** (covered per-task above).
