package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	baseURL     = "http://localhost:8080"
	testAddress = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
)

type subscribeRequest struct {
	Address string `json:"address"`
}

type subscribeResponse struct {
	Success bool `json:"success"`
}

type blockResponse struct {
	CurrentBlock int `json:"current_block"`
}

type testResponse struct {
	statusCode int
	body       []byte
}

// Helper functions
func subscribeAddress(t *testing.T, address string) testResponse {
	reqBody := subscribeRequest{Address: address}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/subscribe", baseURL),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return testResponse{
		statusCode: resp.StatusCode,
		body:       body,
	}
}

func getCurrentBlock(t *testing.T) testResponse {
	resp, err := http.Get(fmt.Sprintf("%s/block", baseURL))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return testResponse{
		statusCode: resp.StatusCode,
		body:       body,
	}
}

func getTransactions(t *testing.T, address string) testResponse {
	resp, err := http.Get(fmt.Sprintf("%s/transactions?address=%s", baseURL, address))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return testResponse{
		statusCode: resp.StatusCode,
		body:       body,
	}
}

func TestParserIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test in non-integration mode")
	}

	// Run test scenarios
	t.Run("success_scenarios", func(t *testing.T) {
		testSuccessScenarios(t)
	})

	t.Run("failure_scenarios", func(t *testing.T) {
		testFailureScenarios(t)
	})
}

func testSuccessScenarios(t *testing.T) {
	t.Run("valid_subscription", func(t *testing.T) {
		resp := subscribeAddress(t, testAddress)
		if resp.statusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.statusCode)
		}

		var subResp subscribeResponse
		if err := json.Unmarshal(resp.body, &subResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if !subResp.Success {
			t.Error("Expected subscription to succeed")
		}
	})

	t.Run("get_current_block", func(t *testing.T) {
		resp := getCurrentBlock(t)
		if resp.statusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.statusCode)
		}

		var blockResp blockResponse
		if err := json.Unmarshal(resp.body, &blockResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if blockResp.CurrentBlock <= 0 {
			t.Error("Expected current block to be greater than 0")
		}
	})

	t.Run("get_transactions", func(t *testing.T) {
		// Wait for transactions to be processed
		time.Sleep(30 * time.Second)

		resp := getTransactions(t, testAddress)
		if resp.statusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.statusCode)
		}

		var transactions []entity.Transaction
		if err := json.Unmarshal(resp.body, &transactions); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		t.Logf("Found %d transactions for address %s", len(transactions), testAddress)

		for _, tx := range transactions {
			if tx.Hash == "" {
				t.Error("Transaction hash should not be empty")
			}
			if tx.BlockNumber <= 0 {
				t.Error("Block number should be greater than 0")
			}
		}
	})
}

func testFailureScenarios(t *testing.T) {
	t.Run("invalid_ethereum_address", func(t *testing.T) {
		invalidAddresses := []string{
			"0xinvalid",
			"not-an-address",
			"0x742d35Cc6634C0532925a3b844Bc454e4438f44",   // Too short
			"0x742d35Cc6634C0532925a3b844Bc454e4438f44e1", // Too long
		}

		for _, addr := range invalidAddresses {
			resp := subscribeAddress(t, addr)
			if resp.statusCode != http.StatusOK {
				t.Errorf("expected status OK, got %v", resp.statusCode)
			}

			var subResp subscribeResponse
			if err := json.Unmarshal(resp.body, &subResp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if subResp.Success {
				t.Errorf("Expected subscription to fail for invalid address: %s", addr)
			}
		}
	})

	t.Run("missing_address_parameter", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/transactions", baseURL))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status BadRequest, got %v", resp.StatusCode)
		}
	})

	t.Run("invalid_json_subscription", func(t *testing.T) {
		invalidJSON := []byte(`{"address": invalid}`)
		resp, err := http.Post(
			fmt.Sprintf("%s/subscribe", baseURL),
			"application/json",
			bytes.NewBuffer(invalidJSON),
		)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status BadRequest, got %v", resp.StatusCode)
		}
	})
}
