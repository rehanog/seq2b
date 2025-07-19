# ADR-001: Testing Strategy for Claude-Assisted Development

Date: 2025-01-18
Status: Accepted

## Context

We are developing a Logseq replacement using Go (backend) and JavaScript (frontend) with the Wails framework. Development is assisted by Claude, an AI assistant with specific constraints:

### Claude's Limitations:
- Cannot make HTTP requests or connect to localhost
- Cannot interact with running GUI applications
- Cannot take screenshots or see visual output
- Cannot control browser instances
- Operates in a sandboxed environment with file system access only

### Current Testing Landscape:
- Backend: Go unit tests
- Frontend: Vitest with jsdom for JavaScript testing
- Integration: Limited testing of Go-JS communication
- Manual: Developer must manually verify GUI behavior

## Decision

We will implement a **file-based testing strategy** that allows Claude to verify application behavior without direct access to the running application.

### Testing Approach:

1. **Test Output Artifacts**: The application will write key outputs to files that Claude can read:
   - Rendered DOM snapshots
   - API request/response logs
   - Navigation event logs
   - Application state exports

2. **Structured Logging**: Implement comprehensive logging that captures:
   - User interactions
   - State transitions
   - Rendering operations
   - Error conditions

3. **Mock-Based Testing**: Expand frontend tests to simulate user interactions with verifiable outputs

4. **State Verification**: Export application state at key points for analysis

## Implementation

### 1. Test Mode Flag
Add a `--test-mode` flag to the application that enables output capture:

```go
// When test mode is enabled
if testMode {
    captureDOM(page)
    logAPICall(request, response)
    exportState(appState)
}
```

### 2. Output Directory Structure
```
frontend/test-output/
├── dom/
│   ├── page-{name}-{timestamp}.html
│   └── navigation-{timestamp}.html
├── api/
│   ├── calls-{timestamp}.json
│   └── responses-{timestamp}.json
├── state/
│   ├── app-state-{timestamp}.json
│   └── navigation-history.json
└── logs/
    ├── user-actions.log
    └── errors.log
```

### 3. Verification Workflow
1. Developer runs: `wails dev --test-mode`
2. Performs test scenario (e.g., create link, navigate, edit)
3. Application writes outputs to files
4. Claude reads files and verifies expected behavior
5. Claude can suggest corrections based on analysis

### 4. Example Test Scenarios

**Navigation Test**:
```javascript
// User action: Click [[Page A]] link
// Expected outputs:
// - dom/navigation-*.html shows Page A content
// - api/calls-*.json shows GetPage("Page A") call
// - state/navigation-history.json shows ["Home", "Page A"]
```

## Consequences

### Positive:
- Claude can verify application behavior without manual testing
- Creates audit trail of all operations
- Enables regression testing by comparing outputs
- Improves debugging with detailed logs
- Can generate test reports from artifacts

### Negative:
- Additional code complexity for output capture
- Performance overhead in test mode
- Storage requirements for test artifacts
- Not real-time - requires app execution then analysis

### Neutral:
- Requires discipline to run in test mode during development
- Test artifacts need periodic cleanup
- May need to sanitize sensitive data in outputs

## Alternatives Considered

1. **Pure Unit Testing**: Insufficient for GUI behavior verification
2. **Manual Testing Only**: Too slow and error-prone
3. **External Testing Service**: Would require Claude to have network access
4. **Video Recording**: Claude cannot process video content

## References

- [Wails Testing Documentation](https://wails.io/docs/howdoesitwork#testing)
- [Testing Without a GUI](https://martinfowler.com/articles/nonDeterminism.html)
- Project testing files: `app_test.go`, `frontend/test/*.js`

## Review

This ADR should be reviewed if:
- Claude's capabilities change (e.g., gains network access)
- We adopt a different frontend framework
- Performance impact becomes unacceptable
- Better testing tools become available