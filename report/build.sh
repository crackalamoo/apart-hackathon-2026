#!/usr/bin/env bash
# Build report.md -> report.pdf using pandoc (markdown -> styled HTML) and
# headless Chromium (HTML -> PDF). Diagrams come from mermaid-filter.
# Requires: pandoc, mermaid-filter (npm -g), and puppeteer's Chromium.
set -euo pipefail
cd "$(dirname "$0")"

CHROME="$HOME/.cache/puppeteer/chrome/mac_arm-150.0.7871.24/chrome-mac-arm64/Google Chrome for Testing.app/Contents/MacOS/Google Chrome for Testing"
export PUPPETEER_EXECUTABLE_PATH="$CHROME"
echo '{"args":["--no-sandbox"]}' > .puppeteer.json

pandoc report.md \
  --from markdown \
  --to html5 \
  --standalone \
  --embed-resources \
  --css style.css \
  --filter mermaid-filter \
  -o report.html

"$CHROME" --headless --disable-gpu --no-sandbox \
  --no-pdf-header-footer \
  --print-to-pdf="report.pdf" "file://$PWD/report.html"

echo "built report.pdf"
