// Package onntrackclient provides a client for the Onntrack tracking dashboard REST API.
package onntrackclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Onntrack API.
	DefaultBaseURL = "https://platform.onntrack.nl/v3/new/"

	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// Client is a client for the Onntrack API.
type Client struct {
	// BaseURL is the base URL for API requests.
	BaseURL *url.URL

	// APIKey is the API key used for authentication.
	APIKey string

	// HTTPClient is the HTTP client used to communicate with the API.
	HTTPClient *http.Client

	// Common service fields
	common service

	// Services used for communicating with different parts of the Onntrack API.
	Auth    *AuthService
	Devices *DevicesService
}

type service struct {
	client *Client
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client) error

// NewClient returns a new Onntrack API client.
func NewClient(options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
	}

	// Apply options
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// Initialize services
	c.common.client = c
	c.Auth = (*AuthService)(&c.common)
	c.Devices = (*DevicesService)(&c.common)

	return c, nil
}

// WithBaseURL sets the base URL for the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.BaseURL = u
		return nil
	}
}

// WithHTTPClient sets the HTTP client for the client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.HTTPClient = httpClient
		return nil
	}
}

// WithAPIKey sets the API key for authentication.
// The apiKey parameter is the authentication token required for API requests. This must be an JWT token.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.APIKey = apiKey
		return nil
	}
}

// NewRequest creates an API request.
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

// ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"message"`
	Code     string `json:"code"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Code)
}
