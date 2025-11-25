package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorecheck/recheck"

	"gopkg.in/yaml.v3"
)

func main() {
	// 1. Determine input
	var patterns []string

	// Check command line args
	if len(os.Args) > 1 {
		patterns = append(patterns, os.Args[1:]...)
	}

	// Check stdin if no args or if piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				patterns = append(patterns, line)
			}
		}
	}

	if len(patterns) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gorecheck <pattern>... or echo <pattern> | gorecheck")
		os.Exit(1)
	}

	// 2. Initialize client
	// Use embedded binary
	binaryPath, err := recheck.GetBinaryPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting binary: %v\n", err)
		os.Exit(1)
	}
	
	client, err := recheck.NewClient(binaryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting recheck: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// 3. Process patterns
	var results []*recheck.Result
	for _, p := range patterns {
		res, err := client.Check(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking pattern '%s': %v\n", p, err)
			continue
		}
		results = append(results, res)
	}

	// 4. Output YAML
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	if err := encoder.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding YAML: %v\n", err)
		os.Exit(1)
	}
}
