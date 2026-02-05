package aria2

import (
	"context"
	"fmt"

	"github.com/siku2/arigo"
)

// Client wraps the aria2 RPC client
type Client struct {
	rpc arigo.Client
	ctx context.Context
}

// NewClient creates a new aria2 RPC client
func NewClient(rpcURL, secret string) (*Client, error) {
	// Default: http://localhost:6800/jsonrpc
	if rpcURL == "" {
		rpcURL = "http://localhost:6800/jsonrpc"
	}

	// Convert http:// to ws:// for arigo WebSocket connection
	wsURL := rpcURL
	if len(rpcURL) >= 7 && rpcURL[:7] == "http://" {
		wsURL = "ws://" + rpcURL[7:]
	} else if len(rpcURL) >= 8 && rpcURL[:8] == "https://" {
		wsURL = "wss://" + rpcURL[8:]
	}

	var client arigo.Client
	var err error

	if secret != "" {
		client, err = arigo.Dial(wsURL, secret)
	} else {
		client, err = arigo.Dial(wsURL, "")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to aria2: %w", err)
	}

	return &Client{
		rpc: client,
		ctx: context.Background(),
	}, nil
}

// Close closes the RPC connection
func (c *Client) Close() error {
	return c.rpc.Close()
}

// Ping checks if the aria2 RPC server is accessible
func (c *Client) Ping() error {
	_, err := c.rpc.GetVersion()
	return err
}
