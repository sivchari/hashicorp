package hashicorp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Interface is an interface for Hashicorp API.
type Interface interface {
	ListProductNames(ctx context.Context) (*ListProductNamesResponse, error)
}

type client struct {
	c *http.Client
}

// Client is an API client for Hashicorp.
type Client struct {
	*client
}

var _ Interface = (*client)(nil)

// ClientOption is defined type for functional option pattern.
type ClientOption func(*client)

// WithHTTPClient returns ClientOption for HTTP Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) {
		c.c = httpClient
	}
}

// New returns a Hashicorp API client.
func New(opts ...ClientOption) *Client {
	c := &client{
		c: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return &Client{
		c,
	}
}

// ListProductNames returns the names of all products on the releases site are returned.
func (c *client) ListProductNames(ctx context.Context) (*ListProductNamesResponse, error) {
	ep := "https://api.releases.hashicorp.com/v1/products"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list product names: %s", err.Error())
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list product names: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			APIName: "list product names",
			Status:  resp.Status,
			URL:     req.URL.String(),
		}
	}
	var arr []string
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, fmt.Errorf("list product names: %s", err.Error())
	}
	body := make(map[string]struct{}, len(arr))
	for _, a := range arr {
		body[a] = struct{}{}
	}
	return &ListProductNamesResponse{
		Products: body,
	}, nil
}
