package deepseek

import (
	"fmt"
	utils "github.com/cohesion-org/deepseek-go/utils"
	"net/http"
	"net/url"
	"time"
)

// BaseURL is the base URL for the Deepseek API
const BaseURL string = "https://api.deepseek.com/v1"

// HTTPDoer is an interface for the Do method of http.Client
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the main struct for interacting with the Deepseek API.
type Client struct {
	AuthToken string        // The authentication token for the API
	BaseURL   string        // The base URL for the API
	Timeout   time.Duration // The timeout for the current Client
	ApiType   utils.ApiType
}

// NewClient creates a new client with an authentication token and an optional custom baseURL.
// If no baseURL is provided, it defaults to "https://api.deepseek.com/".
func NewClient(AuthToken string, baseURL ...string) *Client {
	if AuthToken == "" {
		return nil
	}
	// check if this is a valid URL
	if len(baseURL) > 0 {
		_, err := url.ParseRequestURI(baseURL[0])
		if err != nil {
			fmt.Printf("Invalid URL: %s. \nIf you are using options please use NewClientWithOptions", baseURL[0])
			return nil
		}
	}
	url := "https://api.deepseek.com/"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return &Client{
		AuthToken: AuthToken,
		BaseURL:   url,
		ApiType:   utils.ApiTypeDeepSeek,
	}
}

// Option configures a Client instance
type Option func(*Client) error

// NewClientWithOptions creates a new client with required authentication token and optional configurations.
// Defaults:
// - BaseURL: "https://api.deepseek.com/"
// - Timeout: 5 minutes
func NewClientWithOptions(authToken string, opts ...Option) (*Client, error) {
	client := &Client{
		AuthToken: authToken,
		BaseURL:   "https://api.deepseek.com/",
		Timeout:   5 * time.Minute,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return client, nil
}

// WithBaseURL sets the base URL for the API client
func WithBaseURL(url string) Option {
	return func(c *Client) error {
		c.BaseURL = url
		return nil
	}
}

func WithApiType(apiType utils.ApiType) Option {
	return func(c *Client) error {
		c.ApiType = apiType
		return nil
	}
}

// WithTimeout sets the timeout for API requests
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d < 0 {
			return fmt.Errorf("timeout must be a positive duration")
		}
		c.Timeout = d
		return nil
	}
}

// WithTimeoutString parses a duration string and sets the timeout
// Example valid values: "5s", "2m", "1h"
func WithTimeoutString(s string) Option {
	return func(c *Client) error {
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid timeout duration %q: %w", s, err)
		}
		return WithTimeout(d)(c)
	}
}
