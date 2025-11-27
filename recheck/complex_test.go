package recheck_test

import (
	"testing"

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
