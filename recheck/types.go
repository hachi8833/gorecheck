package recheck

// RPCRequest represents a JSON-RPC 2.0 request
type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// RPCResponse represents a JSON-RPC 2.0 response
type RPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  *Result   `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CheckParams represents the parameters for the "check" method
type CheckParams struct {
	Source string            `json:"source"`
	Flags  string            `json:"flags"`
	Params CheckParamsConfig `json:"params"`
}

// CheckParamsConfig represents the configuration parameters for checking
type CheckParamsConfig struct {
	Timeout interface{} `json:"timeout,omitempty"` // null for no timeout, or number in milliseconds
}

// Result represents the result of a recheck analysis
type Result struct {
	Source     string        `json:"source"`
	Flags      string        `json:"flags"`
	Status     string        `json:"status"`
	Checker    interface{}   `json:"checker"` // Can be null or object
	Error      *RecheckError `json:"error,omitempty"`
	Complexity interface{}   `json:"complexity,omitempty" yaml:",omitempty"`
	Attack     interface{}   `json:"attack,omitempty" yaml:",omitempty"`
	Hotspot    interface{}   `json:"hotspot,omitempty" yaml:",omitempty"`
}

// RecheckError represents an error from the recheck logic (not JSON-RPC error)
type RecheckError struct {
	Kind    string `json:"kind"`
	Message string `json:"message,omitempty"`
}
