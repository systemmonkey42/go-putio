package putio

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent = "go-putio"
	defaultMediaType = "application/json"
	defaultBaseURL   = "https://api.put.io"
	defaultUploadURL = "https://upload.put.io"
)

// errors
var (
	ErrNotExist        = fmt.Errorf("file does not exist")
	ErrPaymentRequired = fmt.Errorf("payment required")

	errRedirect   = fmt.Errorf("redirect attempt on a no-redirect client")
	errNegativeID = fmt.Errorf("file id cannot be negative")
)

// Client manages communication with Put.io v2 API.
type Client struct {
	// HTTP client used to communicate with Put.io API
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Services used for communicating with the API
	Files     *FilesService
	Transfers *TransfersService
	Zips      *ZipsService
	Friends   *FriendsService
	Account   *AccountService
}

// NewClient returns a new Put.io API client, using the htttpClient, which must
// be a new Oauth2 enabled http.Client. If httpClient is not defined, default
// HTTP client is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	c.Files = &FilesService{client: c}
	c.Transfers = &TransfersService{client: c}
	c.Zips = &ZipsService{client: c}
	c.Friends = &FriendsService{client: c}
	c.Account = &AccountService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided via
// relURL, which will be resolved to the BaseURL of the Client.
func (c *Client) NewRequest(method, relURL string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(relURL)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

// redirectOnceFunc follows the redirect only once, and copies the original
// request headers to the new one.
func redirectOnceFunc(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}

	if len(via) > 1 {
		return errRedirect
	}

	// merge headers with request headers
	for header, values := range via[0].Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}
	return nil
}
