# PDF.js Setup Instructions

To add PDF support, we need to:

1. Download PDF.js from: https://cdnjs.com/libraries/pdf.js
2. Add to frontend/index.html
3. Create PDF viewer component

For now, we'll implement a simple approach:
- Detect PDF links in markdown
- Open PDFs in a modal/panel
- Basic controls (prev/next page, zoom)