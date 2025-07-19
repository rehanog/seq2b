#!/bin/bash
# Run the seq2b desktop application with the test library
# Usage: ./scripts/run_seq2b.sh [-dev]
#   -dev: Run in development mode using wails dev instead of the built binary

set -e

# Default to production mode
DEV_MODE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -dev|--dev)
            DEV_MODE=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [-dev]"
            exit 1
            ;;
    esac
done

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Library path - using the test library
LIBRARY_PATH="${PROJECT_ROOT}/testdata/library_test_0/pages"

if [ "$DEV_MODE" = true ]; then
    echo "Running seq2b in development mode..."
    echo "Library: $LIBRARY_PATH"
    cd "${PROJECT_ROOT}/desktop/wails"
    
    # Note: wails dev doesn't support command line arguments to the app
    # So we need to temporarily set an environment variable
    export SEQ2B_LIBRARY_PATH="$LIBRARY_PATH"
    
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
    
    "$WAILS_CMD" dev
else
    echo "Running seq2b in production mode..."
    echo "Library: $LIBRARY_PATH"
    
    # Determine the binary name based on OS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        BINARY="${PROJECT_ROOT}/bin/seq2b.app/Contents/MacOS/seq2b"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        BINARY="${PROJECT_ROOT}/bin/seq2b"
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" ]]; then
        BINARY="${PROJECT_ROOT}/bin/seq2b.exe"
    else
        echo "Error: Unsupported operating system: $OSTYPE"
        exit 1
    fi
    
    # Check if binary exists
    if [ ! -f "$BINARY" ]; then
        echo "Error: Binary not found at $BINARY"
        echo "Please build the application first with:"
        echo "  ./scripts/build_seq2b.sh"
        exit 1
    fi
    
    # Run the binary with library parameter
    "$BINARY" -library "$LIBRARY_PATH"
fi