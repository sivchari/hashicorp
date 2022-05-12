package hashicorp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Interface is an interface for Hashicorp API.
type Interface interface {
	ListProductNames(ctx context.Context) (*ListProductNamesResponse, error)
	ListReleases(ctx context.Context, product string, param ...*ListReleasesParam) (*ListReleasesResponse, error)
	SpecificRelease(ctx context.Context, product string, version string, param ...*SpecificReleaseParam) (*Release, error)
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

// ListReleases returns the release metadata for multiple releases. This endpoint uses pagination for products with many releases. Results are ordered by release creation time from newest to oldest.
func (c *client) ListReleases(ctx context.Context, product string, param ...*ListReleasesParam) (*ListReleasesResponse, error) {
	ep := fmt.Sprintf("https://api.releases.hashicorp.com/v1/releases/%s", product)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list product names: %s", err.Error())
	}
	var p ListReleasesParam
	switch len(param) {
	case 0:
		// do nothing
	case 1:
		p = *param[0]
	default:
		return nil, errors.New("list releases: only one option is allowed")
	}
	q := req.URL.Query()
	var limit int
	if p.Limit > 20 {
		return nil, errors.New("the limit parameter must be  20 or less")
	}
	if p.Limit != 0 {
		limit = p.Limit
	}
	q.Add("limit", strconv.Itoa(limit))
	if p.LicenseClass != "" {
		q.Add("license_class", string(p.LicenseClass))
	}
	if p.After != "" {
		q.Add("after", p.After)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list releases: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			APIName: "list releases",
			Status:  resp.Status,
			URL:     req.URL.String(),
		}
	}
	var rs []*Release
	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, fmt.Errorf("list releases: %s", err.Error())
	}
	return &ListReleasesResponse{
		Releases: rs,
	}, nil
}

// SpecificRelease returns an all metadata for a single product release is returned.
func (c *client) SpecificRelease(ctx context.Context, product string, version string, param ...*SpecificReleaseParam) (*Release, error) {
	ep := fmt.Sprintf("https://api.releases.hashicorp.com/v1/releases/%s/%s", product, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("specific release: %s", err.Error())
	}
	var p SpecificReleaseParam
	switch len(param) {
	case 0:
		// do nothing
	case 1:
		p = *param[0]
	default:
		return nil, errors.New("specific release: %s")
	}
	q := req.URL.Query()
	if p.LicenseClass != "" {
		q.Add("license_class", string(p.LicenseClass))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("specific release: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			APIName: "specific release",
			Status:  resp.Status,
			URL:     req.URL.String(),
		}
	}
	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("specific release: %s", err.Error())
	}
	return &r, nil
}
