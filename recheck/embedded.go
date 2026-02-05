package recheck

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

// GetBinaryPath returns the path to the executable recheck binary.
// It extracts the embedded binary to a temporary file if needed.
func GetBinaryPath() (string, error) {
	// 1. Read embedded content
	data, err := embeddedBinary.ReadFile(binaryName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded binary: %w", err)
	}

	// 2. Determine temp path based on hash to avoid unnecessary extraction
	hash := sha256.Sum256(data)
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("gorecheck-recheck-%x", hash[:8])
	exePath := filepath.Join(tmpDir, fileName)

	// 3. Check if file exists and is correct
	info, err := os.Stat(exePath)
	if err == nil && info.Size() == int64(len(data)) {
		// Already exists, assume it's good (could check hash but size + name hash is decent optimization)
		return exePath, nil
	}

	// 4. Extract file
	// Create with 0755 permissions (rwxr-xr-x)
	if err := os.WriteFile(exePath, data, 0755); err != nil {
		return "", fmt.Errorf("failed to extract binary to %s: %w", exePath, err)
	}

	return exePath, nil
}
