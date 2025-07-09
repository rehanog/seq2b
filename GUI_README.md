# Logseq Go GUI

A simple GUI viewer for Logseq-style markdown files.

## Usage

```bash
# Default: opens Page A from testdata/pages
go run cmd/gui/main.go

# Specify a different starting page
go run cmd/gui/main.go -page "Page B"

# Use a different directory
go run cmd/gui/main.go ~/my-logseq-files

# Both directory and page
go run cmd/gui/main.go -page "My Note" ~/my-logseq-files
```

## Building

```bash
# Build the GUI application
go build -o logseq-gui cmd/gui/main.go

# Run the built application
./logseq-gui
./logseq-gui -page "Page C"
```

## Features

- Single page view (no directory browser)
- Shows block hierarchy with bullet points
- Displays backlinks at the bottom
- Clean, focused interface
- Command-line arguments for page selection

## Navigation (Coming Soon)

Next update will add clickable [[page links]] for navigation between pages.