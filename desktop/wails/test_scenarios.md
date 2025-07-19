# Test Scenarios for seq2b

This document describes test scenarios that can be run with test mode enabled to verify application behavior.

## Running Test Mode

```bash
cd desktop/wails
go run github.com/wailsapp/wails/v2/cmd/wails@latest dev -- -test-mode
```

Or if wails is installed:
```bash
wails dev -- -test-mode
```

Test outputs will be captured in `frontend/test-output/`

## Test Scenario 1: Page Navigation with Case Sensitivity

**Goal**: Verify that page navigation is case-insensitive

**Steps**:
1. Start the app in test mode
2. Navigate to the browser at http://localhost:34115
3. Create a link: Type `[[page a]]` in any block
4. Click the link
5. Verify navigation succeeds

**Expected Outputs**:
- `logs/user-actions-{date}.log`: Should contain "Navigate to page: page a"
- `api/api-calls-{date}.json`: Should show GetPage("page a") call
- `dom/page-page-a-{timestamp}.html`: Should show the rendered page
- `state/navigation-history.json`: Should include "page a" in history

## Test Scenario 2: Edit Persistence

**Goal**: Verify edits are saved before navigation

**Steps**:
1. Start the app in test mode
2. Click on any block to edit
3. Add text including a link: `Test content with [[new page]]`
4. Click the [[new page]] link WITHOUT pressing Enter/Escape
5. Click back button

**Expected Outputs**:
- `logs/user-actions-{date}.log`: Should show edit action before navigation
- `api/api-calls-{date}.json`: Should show UpdateBlockAtPath call
- Navigation should preserve the edit

## Test Scenario 3: Backlink Updates

**Goal**: Verify backlinks update correctly

**Steps**:
1. Start the app in test mode
2. Navigate to any page
3. Add a link to another page: `[[Page B]]`
4. Navigate to Page B
5. Check backlinks section

**Expected Outputs**:
- `dom/page-page-b-{timestamp}.html`: Should show backlinks from the source page
- Backlinks should be visible in the rendered DOM

## Verification Process

After running a test scenario:

1. Check the log files were created:
   ```bash
   ls -la frontend/test-output/logs/
   ```

2. Verify API calls:
   ```bash
   cat frontend/test-output/api/api-calls-*.json | jq '.'
   ```

3. Check captured DOM:
   ```bash
   ls -la frontend/test-output/dom/
   ```

4. Verify navigation history:
   ```bash
   cat frontend/test-output/state/navigation-history.json
   ```

## Automated Verification

Run the verification script:
```bash
go run verify_test_output.go
```

This will analyze the captured outputs and report any issues.