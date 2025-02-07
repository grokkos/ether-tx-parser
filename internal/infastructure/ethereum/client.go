package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	ethtypes "github.com/grokkos/ether-tx-parser/pkg/ethereum"
	"net/http"
)

type Client struct {
	rpcURL string
}

func NewClient(rpcURL string) *Client {
	return &Client{rpcURL: rpcURL}
}

func (c *Client) MakeRPCCall(method string, params []interface{}) (*ethtypes.JSONRPCResponse, error) {
	request := ethtypes.JSONRPCRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(c.rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var response ethtypes.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}
