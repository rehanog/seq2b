#!/bin/bash

# Script to download Logseq test graphs for compatibility testing
# Usage: ./download-logseq-graphs.sh

set -e

# Create tmp directory if it doesn't exist
GRAPHS_DIR="tmp/logseq-test-graphs"
mkdir -p "$GRAPHS_DIR"

echo "Downloading Logseq test graphs to $GRAPHS_DIR..."

# List of repositories to clone
declare -a repos=(
    "candideu/Logseq-Demo-Graph"
    "logseq/docs"
    # Add more repositories here as needed
)

# Clone or update each repository
for repo in "${repos[@]}"; do
    repo_name=$(basename "$repo")
    repo_path="$GRAPHS_DIR/$repo_name"
    
    if [ -d "$repo_path" ]; then
        echo "Updating $repo_name..."
        cd "$repo_path"
        git pull --quiet
        cd - > /dev/null
    else
        echo "Cloning $repo_name..."
        git clone --quiet --depth 1 "https://github.com/$repo.git" "$repo_path"
    fi
done

# Download specific test files from gists or other sources
# Example:
# curl -s https://gist.github.com/username/gistid/raw/file.md > "$GRAPHS_DIR/specific-test.md"

echo "Download complete!"
echo ""
echo "Available test graphs:"
ls -la "$GRAPHS_DIR"
echo ""
echo "To analyze a graph, use:"
echo "  go run cmd/seq2b/main.go $GRAPHS_DIR/<graph-name>"