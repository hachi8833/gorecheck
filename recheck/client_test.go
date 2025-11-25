package recheck_test

import (
	"testing"

	"gorecheck/recheck"
)

func TestClient(t *testing.T) {
	// 1. Get binary path
	binPath, err := recheck.GetBinaryPath()
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	// 2. Create client
	client, err := recheck.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 3. Test Safe Pattern
	t.Run("Safe Pattern", func(t *testing.T) {
		res, err := client.Check("^a$")
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status != "safe" {
			t.Errorf("Expected status 'safe', got '%s'", res.Status)
		}
	})

	// 4. Test Vulnerable Pattern
	t.Run("Vulnerable Pattern", func(t *testing.T) {
		res, err := client.Check("^(a|a)*$")
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}
		if res.Status != "vulnerable" {
			t.Errorf("Expected status 'vulnerable', got '%s'", res.Status)
		}
	})

	// 5. Test Invalid Pattern (should return error or unknown status depending on recheck behavior)
	t.Run("Invalid Pattern", func(t *testing.T) {
		res, err := client.Check("(")
		if err != nil {
			// It might return an RPC error or a result with error
			t.Logf("Got expected error: %v", err)
		} else {
			if res.Status != "unknown" && res.Error == nil {
				t.Errorf("Expected error or unknown status for invalid pattern, got status '%s'", res.Status)
			}
		}
	})
}
