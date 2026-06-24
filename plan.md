
Here's the full summary:
Portfolio Website Plan
Aesthetic

Neo-brutalism with a comic/tabloid magazine flavor. Key elements:

    Hard black borders, raw box-shadows (box-shadow: 7px 7px 0px #000), intentional inconsistency
    Display typography (Anton, Bebas Neue, or Black Han Sans) + monospace body (Courier Prime, JetBrains Mono)
    White/cream base, black borders, one acid accent (yellow #FFE600, hot pink, or safety orange)
    Layout as comic panels — hard borders, slight rotations, distinct blocks
    Signature element: a stamp/sticker overlay ("AVAILABLE FOR HIRE", "★ FEATURED PROJECT ★")

Tech Stack

    Go templates — server-rendered HTML, fast, deployable as a single binary
    Plain CSS — raw and deliberate, fights Tailwind's normalizing tendency
    HTMX (optional) — one or two interactions max (form submit, load more, etc.)
    No React

Dev vs. Release Embedding

Use build tags to switch between live disk files in dev and a fully embedded binary in release:
go

//go:build release
var staticFiles = http.FS(embeddedFS) // prod: embedded

go

//go:build !release
var staticFiles = http.Dir("static")  // dev: live from disk

Templates re-parsed on each request in dev for hot reload with no extra tooling.
Directory Structure

/
├── main.go
├── templates/
│   ├── base.html       # layout shell
│   ├── index.html      # hero "cover page"
│   ├── projects.html   # comic panel grid
│   └── contact.html
├── static/
│   ├── style.css
│   └── htmx.min.js     # optional
└── content/
    └── projects.json   # project data


