# Bug Investigator Persistent Memory

## Project Context
**Project**: Seq2b - Logseq replacement in Go  
**Focus**: High-performance, secure, reliable sync, AI-first  
**Testing Philosophy**: TDD for bug fixes, comprehensive regression prevention  

## Known Bug Patterns

### Frontend Testing Issues
- **DOM Mocking Gaps**: Vitest tests fail due to incomplete DOM structure mocking
- **Wails API Mocking**: Missing critical functions like `IsTestMode()` in test setup
- **Test Status**: 6/7 JavaScript unit tests failing as of 2025-08-23

### Parser Bugs
- **Image Handling**: Recent fixes for image rendering with asset serving API
- **Property Parsing**: New test file suggests property parsing issues

## Test Coverage Status

### Well-Tested Areas ✅
- `app_test.go`: Core app functionality
- `integration_test.go`: File system operations
- Block manipulation (add, update, delete)
- Page navigation and retrieval

### Needs Coverage ⚠️
- Frontend UI interactions (Playwright needed)
- Wiki link navigation flows
- Image/PDF rendering pipelines
- Concurrent editing scenarios
- Search functionality

## Active Investigations

### Current Focus: Frontend Testing Framework
**Date**: 2025-08-28  
**Issue**: Deciding between Playwright, Rod, CDP direct, or runtime mocking for UI testing  
**Analysis**: 
- **Playwright**: Best for rapid development, auto-waiting, debugging tools
- **CDP Direct**: Maximum control but 5-10x more code, maintenance burden
- **Rod**: Go-native middle ground, smaller community
- **Runtime Mocking**: Fast but doesn't test actual frontend

**Final Recommendation**: Playwright for E2E testing
- Wails apps are web frontends - Playwright is purpose-built for this
- Faster test development (5-10 lines vs 50+ for CDP)
- Superior debugging (screenshots, videos, trace viewer)
- Auto-handles timing issues that plague CDP direct
- Can drop to CDP for specific features via `newCDPSession()`

**Next Steps**: 
1. Fix mock setup in `frontend/test/setup.js`
2. Install Playwright and create initial test suite
3. Focus on critical user flows: navigation, editing, wiki links
4. Add visual regression tests for UI consistency

## Bug History

### Fixed Bugs
1. **Image Rendering** (commit: 1586207)
   - Issue: Assets not serving correctly
   - Fix: Implemented asset serving API
   - Test: Added image editing tests

## Testing Strategies

### Preferred Patterns
1. **Table-Driven Tests**: For multiple input scenarios
2. **Subtests**: For grouped test cases with shared setup
3. **Benchmark Tests**: For performance-critical paths
4. **Integration Tests**: For multi-component workflows

### Test File Naming
- Unit tests: `*_test.go` in same package
- Integration: `integration_test.go`
- E2E: `e2e/` directory with descriptive names
- Benchmarks: `*_bench_test.go`

## Code Smells to Watch

### Frontend
- Missing null checks on DOM queries
- Unhandled promise rejections
- Race conditions in async navigation

### Backend
- File handle leaks
- Concurrent map access without locks
- Missing error propagation

## Useful Commands

```bash
# Run all Go tests
go test ./...

# Run with race detection
go test -race ./...

# Run JavaScript tests (currently failing)
cd desktop/wails/frontend && npm test

# Build and run GUI
./scripts/build_seq2b.sh
cd desktop/wails && wails dev

# Run specific test
go test -run TestUpdateBlockAfterAddBlock ./desktop/wails
```

## Test Template Usage

### Concurrent Test
Use `templates/concurrent_test.go` when testing:
- Simultaneous block edits
- Multi-user scenarios
- Race conditions

### Benchmark Test
Use `templates/benchmark_test.go` for:
- Parser performance
- Search operations
- Large file handling

### Table Test
Use `templates/table_test.go` for:
- Input validation
- Markdown parsing variations
- Edge case coverage

## Notes for Future Investigations

- Monitor `testdata/test_image_bug.md` for regression patterns
- Check if property parsing issues relate to Logseq compatibility
- Consider adding mutation testing for parser logic
- Profile memory usage during large graph operations

---
*Last Updated: 2025-08-23*  
*Active Agent: bug-investigator-tester-go*