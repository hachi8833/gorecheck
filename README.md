# gorecheck

`gorecheck` is a CLI tool written in Go that checks for ReDoS (Regular Expression Denial of Service) vulnerabilities. It wraps the [recheck](https://github.com/makenowjust-labs/recheck) library (specifically the native binary) to provide a dependency-free, easy-to-use interface.

## Features

- **Dependency-Free**: The `recheck` binary is embedded, so you only need the `gorecheck` executable.
- **Standard Input Support**: Pipe regex patterns directly into the tool.
- **YAML Output**: Results are formatted in YAML for easy reading and parsing.
- **Cross-Platform (Limited)**: Currently embeds the macOS arm64 binary.

## Installation

### From Source

Requirements: Go 1.25+

```bash
git clone <repository-url>
cd gorecheck
go build -o gorecheck main.go
```

## Usage

### Check a single pattern

```bash
./gorecheck "^(a|a)*$"
```

### Check multiple patterns via pipe

```bash
cat patterns.txt | ./gorecheck
```

### Output Format

The output is a YAML list of results:

```yaml
- source: ^(a|a)*$
  flags: ""
  status: vulnerable
  checker: automaton
  error: null
```

## Development

### Dependencies

- Node.js (to download the initial `recheck` binary)

### Setup

1. Install dependencies:
   ```bash
   npm install
   ```
2. Copy the binary for embedding (if updated):
   ```bash
   mkdir -p recheck/bin
   cp node_modules/recheck-macos-arm64/recheck recheck/bin/recheck
   ```
3. Build:
   ```bash
   go build -o gorecheck main.go
   ```

## Testing

Run all tests:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./recheck
```

Run specific test suites:
```bash
# Test complex Japanese patterns
go test -v ./recheck -run TestComplexJapanesePatterns

# Test timeout handling
go test -v ./recheck -run TestNoTimeout

# Test large response handling
go test -v ./recheck -run TestLargeResponse

# Test multiple patterns
go test -v ./recheck -run TestMultiplePatterns

# Test edge cases
go test -v ./recheck -run TestEdgeCases
```

The test suite includes:
- **TestClient**: Basic functionality tests
- **TestComplexJapanesePatterns**: Tests for the bug fix (complex Japanese regex patterns)
- **TestNoTimeout**: Verifies that complex patterns don't timeout
- **TestLargeResponse**: Tests handling of large responses from recheck
- **TestMultiplePatterns**: Tests sequential pattern checking
- **TestEdgeCases**: Tests edge cases (empty character class, Unicode, lookahead)

Current test coverage: **73.1%**

For more details on the bug fix, see [BUG_FIX.md](BUG_FIX.md).

## License

MIT
