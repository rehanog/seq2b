#!/bin/bash

# Get the last added file in ~/Documents/screenshots
SCREENSHOT_DIR="$HOME/Documents/screenshots"
LAST_FILE=$(ls -t "$SCREENSHOT_DIR" | head -1)

if [ -z "$LAST_FILE" ]; then
    echo "No files found in $SCREENSHOT_DIR"
    exit 1
fi

FULL_PATH="$SCREENSHOT_DIR/$LAST_FILE"

# Display the screenshot using Quick Look
qlmanage -p "$FULL_PATH" >/dev/null 2>&1 &
QL_PID=$!

# Get file extension
EXT="${LAST_FILE##*.}"

# Ask for filename prefix (while preview is still open)
read -p "Enter filename prefix: " PREFIX

# Kill the Quick Look process after getting prefix
kill $QL_PID 2>/dev/null || true

# Create date string in yyyymmdd format
DATE=$(date +%Y%m%d)

# Create full filename
NEW_FILENAME="${PREFIX}_${DATE}.${EXT}"

# Move and rename file to current directory
NEW_PATH="$HOME/sandbox/seq2b/tmp/screenshots/$NEW_FILENAME"
echo "Executing mvv $FULL_PATH $NEW_PATH"
mv "$FULL_PATH" "$NEW_PATH"

echo "Screenshot moved and renamed to:"
echo "${NEW_PATH}"
TRUNCATED_PATH="seq2b/${NEW_PATH#*/seq2b/}"
echo "$TRUNCATED_PATH"
echo "$TRUNCATED_PATH" | pbcopy
