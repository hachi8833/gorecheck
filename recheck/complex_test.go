package recheck_test

import (
	"testing"
	"time"

	"gorecheck/recheck"
)

func TestComplexJapanesePatterns(t *testing.T) {
	// Get binary path
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	// Create client
	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test vulnerable pattern (should detect as vulnerable, not timeout)
	t.Run("Vulnerable Japanese Pattern", func(t *testing.T) {
		pattern := `(?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])`
		res, err := client.Check(pattern)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status != "vulnerable" {
			t.Errorf("Expected status 'vulnerable', got '%s'", res.Status)
		}
		t.Logf("Pattern correctly identified as vulnerable using checker: %v", res.Checker)
	})

	// Test safe pattern (fixed version with quantifier limit)
	t.Run("Safe Japanese Pattern", func(t *testing.T) {
		pattern := `(?<![がうち人達])が[、，](?=([^。．:：！？」\n]){1,29}?[^方]が[、，])`
		res, err := client.Check(pattern)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status != "safe" {
			t.Errorf("Expected status 'safe', got '%s'", res.Status)
		}
		t.Logf("Pattern correctly identified as safe using checker: %v", res.Checker)
	})
}

// TestNoTimeout verifies that complex patterns don't timeout
func TestNoTimeout(t *testing.T) {
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// This pattern would timeout with the default 10-second timeout
	// but should complete successfully with timeout disabled
	t.Run("Complex Pattern Should Not Timeout", func(t *testing.T) {
		pattern := `(?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])`

		start := time.Now()
		res, err := client.Check(pattern)
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}

		// Should not return timeout status
		if res.Status == "timeout" {
			t.Errorf("Pattern should not timeout, but got status 'timeout'")
		}

		// Should return vulnerable
		if res.Status != "vulnerable" {
			t.Errorf("Expected status 'vulnerable', got '%s'", res.Status)
		}

		t.Logf("Analysis completed in %v (status: %s)", duration, res.Status)
	})
}

// TestLargeResponse verifies that large responses are handled correctly
func TestLargeResponse(t *testing.T) {
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Patterns that generate large attack strings
	testCases := []struct {
		name    string
		pattern string
		status  string
	}{
		{
			name:    "Japanese Pattern with Large Response",
			pattern: `(?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])`,
			status:  "vulnerable",
		},
		{
			name:    "Classic ReDoS Pattern",
			pattern: `^(a|a)*$`,
			status:  "vulnerable",
		},
		{
			name:    "Nested Quantifiers",
			pattern: `^(a+)+$`,
			status:  "vulnerable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Check(tc.pattern)
			if err != nil {
				t.Fatalf("Check failed: %v", err)
			}

			if res.Status != tc.status {
				t.Errorf("Expected status '%s', got '%s'", tc.status, res.Status)
			}

			// Verify we got a result (not a buffer overflow error)
			if res.Source == "" {
				t.Error("Response should contain source pattern")
			}

			t.Logf("Successfully processed pattern with status: %s, checker: %v", res.Status, res.Checker)
		})
	}
}

// TestMultiplePatterns verifies that multiple patterns can be checked sequentially
func TestMultiplePatterns(t *testing.T) {
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	patterns := []struct {
		pattern string
		status  string
	}{
		{`^a$`, "safe"},
		{`^(a|a)*$`, "vulnerable"},
		{`(?<![がうち人達])が[、，](?=([^。．:：！？」\n]){1,29}?[^方]が[、，])`, "safe"},
		{`(?<![がうち人達])が[、，](?=([^。．:：！？」\n])+?[^方]が[、，])`, "vulnerable"},
		{`^[a-z]+$`, "safe"},
	}

	for i, p := range patterns {
		t.Run(p.pattern, func(t *testing.T) {
			res, err := client.Check(p.pattern)
			if err != nil {
				t.Fatalf("Check %d failed: %v", i, err)
			}

			if res.Status != p.status {
				t.Errorf("Pattern %d: expected status '%s', got '%s'", i, p.status, res.Status)
			}

			t.Logf("Pattern %d: %s -> %s (checker: %v)", i, p.pattern, res.Status, res.Checker)
		})
	}
}

// TestEdgeCases verifies edge cases are handled correctly
func TestEdgeCases(t *testing.T) {
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	t.Run("Empty Character Class", func(t *testing.T) {
		pattern := `^[]$`
		res, err := client.Check(pattern)
		// This might be an error or unknown status depending on recheck behavior
		if err == nil && res.Status == "timeout" {
			t.Error("Should not timeout on simple patterns")
		}
		t.Logf("Empty character class result: %v, error: %v", res, err)
	})

	t.Run("Unicode Pattern", func(t *testing.T) {
		pattern := `^[あ-ん]+$`
		res, err := client.Check(pattern)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status == "timeout" {
			t.Error("Should not timeout on simple unicode patterns")
		}
		t.Logf("Unicode pattern status: %s", res.Status)
	})

	t.Run("Complex Lookahead", func(t *testing.T) {
		pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$`
		res, err := client.Check(pattern)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status == "timeout" {
			t.Error("Should not timeout on lookahead patterns")
		}
		t.Logf("Lookahead pattern status: %s", res.Status)
	})
}
