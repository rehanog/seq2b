#!/bin/bash
# Build the seq2b desktop application
# This script builds the Wails app and copies it to the bin/ directory

set -e  # Exit on error

echo "Building seq2b desktop application..."

# Change to the desktop/wails directory
cd "$(dirname "$0")/../desktop/wails"

# Find wails executable
WAILS_CMD=""
if command -v wails &> /dev/null; then
    WAILS_CMD="wails"
elif [ -x "$(go env GOPATH)/bin/wails" ]; then
    WAILS_CMD="$(go env GOPATH)/bin/wails"
else
    echo "Error: wails CLI not found. Please install it with:"
    echo "  go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Build the application
echo "Running wails build..."
# Set environment variable to indicate we're in build mode
SEQ2B_BUILD_MODE=1 "$WAILS_CMD" build

# Create bin directory at project root if it doesn't exist
echo "Creating bin directory..."
mkdir -p ../../bin

# Copy the built binary to bin/ directory
echo "Copying binary to bin/ directory..."

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    if [ -d "build/bin/seq2b.app" ]; then
        echo "Copying macOS app bundle..."
        rm -rf ../../bin/seq2b.app
        cp -r build/bin/seq2b.app ../../bin/seq2b.app
        echo "✓ Binary available at: bin/seq2b.app"
    else
        echo "Error: macOS app bundle not found at build/bin/seq2b.app"
        exit 1
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if [ -f "build/bin/seq2b" ]; then
        echo "Copying Linux binary..."
        cp build/bin/seq2b ../../bin/seq2b
        chmod +x ../../bin/seq2b
        echo "✓ Binary available at: bin/seq2b"
    else
        echo "Error: Linux binary not found at build/bin/seq2b"
        exit 1
    fi
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" ]]; then
    # Windows
    if [ -f "build/bin/seq2b.exe" ]; then
        echo "Copying Windows executable..."
        cp build/bin/seq2b.exe ../../bin/seq2b.exe
        echo "✓ Binary available at: bin/seq2b.exe"
    else
        echo "Error: Windows executable not found at build/bin/seq2b.exe"
        exit 1
    fi
else
    echo "Error: Unsupported operating system: $OSTYPE"
    exit 1
fi

echo ""
echo "Build complete! 🎉"
echo ""
echo "To run seq2b:"
echo "  ./scripts/run_seq2b.sh         # Run with test library"
echo "  ./bin/seq2b -library /path     # Run with custom library"