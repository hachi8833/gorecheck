# Bug Fix: Timeout Issues with Complex Regex Patterns

## Problem

The CLI was incorrectly reporting complex regex patterns as "timeout" instead of properly analyzing them as "vulnerable". For example:

**Vulnerable pattern:**
```
(?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])
```

This pattern was instantly detected as vulnerable in the recheck playground (https://makenowjust-labs.github.io/recheck/playground/), but the CLI would timeout.

**Safe pattern (for comparison):**
```
(?<![がうち人達])が[、，](?=([^。．:：！？」\n]){1,29}?[^方]が[、，])
```

## Root Causes

There were two issues:

### 1. Default Timeout Too Short

The recheck library has a default timeout of 10 seconds (10000ms). Complex regex patterns can take longer to analyze, especially when the fuzzing checker needs to generate attack strings.

The CLI was not passing any timeout parameter, so it used the default 10-second timeout. The playground likely uses a much longer timeout or no timeout at all.

### 2. Scanner Buffer Too Small

When recheck finds a vulnerability, it returns a large response containing attack strings. The default `bufio.Scanner` buffer size is 64KB, which was insufficient for these large responses, causing a "token too long" error.

## Solution

### Fix 1: Disable Timeout

Modified `recheck/types.go` to add a `CheckParamsConfig` struct with a `Timeout` field:

```go
type CheckParamsConfig struct {
    Timeout interface{} `json:"timeout,omitempty"` // null for no timeout, or number in milliseconds
}
```

Updated `recheck/client.go` to set `Timeout: nil` (which becomes `null` in JSON), telling recheck to never timeout:

```go
Params: CheckParamsConfig{
    Timeout: nil, // null = no timeout limit
},
```

### Fix 2: Increase Scanner Buffer

Modified `recheck/client.go` to increase the scanner buffer from 64KB to 10MB:

```go
scanner := bufio.NewScanner(stdoutPipe)
const maxScanTokenSize = 10 * 1024 * 1024 // 10MB
buf := make([]byte, maxScanTokenSize)
scanner.Buffer(buf, maxScanTokenSize)
```

## Verification

After the fix, both patterns are correctly analyzed:

**Vulnerable pattern:**
```yaml
- source: (?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])
  flags: ""
  status: vulnerable
  checker: fuzz
  error: null
```

**Safe pattern:**
```yaml
- source: (?<![がうち人達])が[、，](?=([^。．:：！？」\n]){1,29}?[^方]が[、，])
  flags: ""
  status: safe
  checker: automaton
  error: null
```

## Test Suite

A comprehensive test suite was added in `recheck/complex_test.go` with 73.1% code coverage:

### 1. `TestComplexJapanesePatterns`
Tests the specific Japanese regex patterns from the bug report:
- Vulnerable pattern: correctly identified as `vulnerable`
- Safe pattern (with quantifier limit): correctly identified as `safe`

### 2. `TestNoTimeout`
Verifies that complex patterns complete analysis without timing out:
- Confirms the pattern doesn't return `timeout` status
- Measures and logs analysis duration
- Validates the pattern is correctly identified as `vulnerable`

### 3. `TestLargeResponse`
Tests handling of large responses (attack strings) from recheck:
- Japanese pattern with large response
- Classic ReDoS pattern: `^(a|a)*$`
- Nested quantifiers: `^(a+)+$`
- Verifies the scanner buffer can handle large JSON responses

### 4. `TestMultiplePatterns`
Verifies sequential checking of multiple patterns:
- Tests 5 different patterns in sequence
- Ensures the client can handle multiple requests without issues
- Validates both safe and vulnerable patterns

### 5. `TestEdgeCases`
Tests edge cases to ensure robustness:
- Empty character class patterns
- Unicode patterns (Japanese characters)
- Complex lookahead patterns
- Ensures no timeouts on simple patterns

All tests pass successfully:
```
=== RUN   TestClient
--- PASS: TestClient (0.07s)
=== RUN   TestComplexJapanesePatterns
--- PASS: TestComplexJapanesePatterns (0.27s)
=== RUN   TestNoTimeout
--- PASS: TestNoTimeout (0.26s)
=== RUN   TestLargeResponse
--- PASS: TestLargeResponse (0.25s)
=== RUN   TestMultiplePatterns
--- PASS: TestMultiplePatterns (0.25s)
=== RUN   TestEdgeCases
--- PASS: TestEdgeCases (0.24s)
PASS
ok      gorecheck/recheck       1.868s  coverage: 73.1% of statements
```
