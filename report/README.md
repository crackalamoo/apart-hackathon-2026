# Report

The submission writeup. Written in Markdown (`report.md`), built to `report.pdf`.

## Build

```
./build.sh
```

This runs `report.md` through pandoc (Markdown → styled HTML, using `style.css`)
and then headless Chromium (HTML → PDF). Mermaid code blocks become diagrams via
`mermaid-filter`.

## Requirements

- `pandoc`
- `mermaid-filter` (`npm install -g mermaid-filter`) and its bundled Chromium

## Editing

- Prose goes in `report.md`. Bracketed `[notes]` mark what each section is for —
  replace them as you write.
- Diagrams are the ` ```mermaid ` blocks — edit labels in place.
- Appearance (fonts, spacing, colors) is in `style.css`. The body font is
  Charter.
- `report.html`, `report.pdf`, and `.puppeteer.json` are build outputs and are
  gitignored.
