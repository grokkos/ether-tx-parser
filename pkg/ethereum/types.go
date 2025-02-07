package ethereum

type JSONRPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type JSONRPCResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  string      `json:"result"`
	ID      int         `json:"id"`
	Error   interface{} `json:"error,omitempty"`
}
