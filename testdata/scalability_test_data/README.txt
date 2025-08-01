SCALABILITY TEST DATA
====================

This directory contains auto-generated test data for performance and scalability testing.

DO NOT EDIT THESE FILES MANUALLY - They are generated by tools and may be regenerated at any time.

Contents:
- test-vault-100/   - 100 test pages
- test-vault-1000/  - 1,000 test pages  
- test-vault-5000/  - 5,000 test pages (may only contain README)

Generator Tool:
These vaults are generated using: tools/generate-test-vault/main.go

Usage:
  go run tools/generate-test-vault/main.go -pages 1000 -output testdata/scalability_test_data/test-vault-1000

The generated vaults simulate realistic Logseq libraries with:
- Various page types (Daily Notes, Project Pages, Meeting Notes, etc.)
- Interconnected pages with [[wiki links]]
- Configurable link density
- Date-based journal pages

These are used by our benchmarking tools to test parser performance at scale.