#!/bin/bash

# Find the Chrome/Chromium executable in Playwright's cache
CHROME_PATH=$(find ~/.cache/ms-playwright -name "chrome" -o -name "chromium" | grep "chrome-linux64/chrome" | head -n 1)

if [ -z "$CHROME_PATH" ]; then
    echo "Browser not found. Installing..."
    npx playwright install chromium
    CHROME_PATH=$(find ~/.cache/ms-playwright -name "chrome" -o -name "chromium" | grep "chrome-linux64/chrome" | head -n 1)
fi

echo "Starting browser from: $CHROME_PATH"

# Kill existing instance if any
pkill -f "remote-debugging-port=9222" || true

# Start Chrome with required flags for this environment
nohup "$CHROME_PATH" \
    --headless \
    --no-sandbox \
    --disable-setuid-sandbox \
    --remote-debugging-address=0.0.0.0 \
    --remote-debugging-port=9222 \
    > /tmp/chrome.log 2>&1 &

echo "Browser started on port 9222. Waiting for it to be ready..."
sleep 2

if lsof -i :9222 > /dev/null; then
    echo "Browser is UP and listening on port 9222."
else
    echo "Failed to start browser. Check /tmp/chrome.log"
    cat /tmp/chrome.log
    exit 1
fi
