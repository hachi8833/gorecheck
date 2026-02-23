# gorecheck

`gorecheck` is a CLI tool written in Go that checks for ReDoS (Regular Expression Denial of Service) vulnerabilities. It wraps the [recheck](https://github.com/makenowjust-labs/recheck) library (specifically the native binary) to provide a dependency-free, easy-to-use interface.

## Features

- **Dependency-Free**: The `recheck` binary is embedded, so you only need the `gorecheck` executable.
- **Standard Input Support**: Pipe regex patterns directly into the tool.
- **YAML Output**: Results are formatted in YAML for easy reading and parsing.
- **Cross-Platform (Limited)**: Currently embeds the macOS arm64 binary.
- **Ruby Regex Normalization**: Ruby-specific syntax is automatically normalized before checking (see below).

## Ruby Regex Normalization

Ruby supports regex syntax that recheck does not understand natively. `gorecheck` automatically normalizes the following constructs before passing patterns to recheck:

### Unicode code point syntax

Ruby allows `\u{HHHHH}` for Unicode code points. These are converted to the standard `\uHHHH` form:

| Input | Normalized |
|---|---|
| `\u{2029}` | `\u2029` |
| `\u{1F600}` | `\u1F600` |

### Inline flag groups

Ruby supports `(?flags:X)` to enable or disable match flags within a group. The flags and surrounding syntax are stripped, leaving only the inner expression:

| Input | Normalized |
|---|---|
| `(?i:ancesfor)` | `ancesfor` |
| `(?imxs:pattern)` | `pattern` |
| `(?-i:pattern)` | `pattern` |
| `(?i-m:pattern)` | `pattern` |

Supported flags: `i` (case-insensitive), `m` (multiline), `x` (extended), `s` (single-line).

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

- [mise](https://mise.jdx.dev/) — manages Go and Node.js versions, and provides task runner
- Node.js (managed via mise) — used to install and update the `recheck` npm package
- Go 1.25+ (managed via mise)

### Setup

1. Install mise: https://mise.jdx.dev/getting-started.html

2. Install the required Go and Node.js versions:
   ```bash
   mise install
   ```

3. Install npm dependencies:
   ```bash
   npm install
   ```

4. Build:
   ```bash
   mise run build
   ```

### Available Tasks

| Task | Description |
|---|---|
| `mise run build` | Download the Linux recheck binary (if needed) and build both binaries |
| `mise run test` | Run all tests |
| `mise run recheck-update` | Update the `recheck` npm package to the latest version |
| `mise run recheck-version` | Show the currently installed `recheck` version |
| `mise run clean` | Remove build artifacts (`gorecheck`, `gorecheck-linux-amd64`) |

> **Note:** `mise run build` skips the download and build if all binaries are already present.
> To force a rebuild, run `mise run clean` first.

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
