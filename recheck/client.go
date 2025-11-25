package recheck

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Client handles communication with the recheck binary
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	mutex  sync.Mutex
	nextID int
}

// NewClient starts the recheck binary in agent mode
func NewClient(binaryPath string) (*Client, error) {
	cmd := exec.Command(binaryPath, "agent")
	// Inherit stderr to see logs if any
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start recheck agent: %w", err)
	}

	return &Client{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdoutPipe),
		nextID: 1,
	}, nil
}

// Check sends a check request to the recheck agent
func (c *Client) Check(source string) (*Result, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	id := c.nextID
	c.nextID++

	req := RPCRequest{
		JSONRPC: "2.0+push",
		Method:  "check",
		Params: CheckParams{
			Source: source,
			Flags:  "",
			Params: map[string]interface{}{},
		},
		ID: id,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Write request to stdin with newline
	if _, err := c.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write to stdin: %w", err)
	}

	// Read response from stdout
	for c.stdout.Scan() {
		line := c.stdout.Bytes()
		var resp RPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			// Skip invalid lines or logs that aren't JSON-RPC
			continue
		}

		if resp.ID == id {
			if resp.Error != nil {
				return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
			}
			return resp.Result, nil
		}
	}

	if err := c.stdout.Err(); err != nil {
		return nil, fmt.Errorf("error reading stdout: %w", err)
	}

	return nil, fmt.Errorf("stream ended without response for id %d", id)
}

// Close terminates the recheck agent
func (c *Client) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}
