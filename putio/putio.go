// Package putio is the Put.io API v2 client for Go.
package putio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultBaseURL   = "https://api.put.io"
	defaultUserAgent = "go-putio"
	defaultMediaType = "application/json"
)

var errRedirectAttempt = errors.New("redirect attempt through a redirect-prevented HTTP client")

// Client manages communication with Put.io v2 API.
type Client struct {
	// HTTP client used to communicate with Put.io API
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string
}

// NewClient returns a new Put.io API client. It is possible to pass a custom
// http.Client. If httpClient is not defined, default HTTP client is used.
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	return c, nil
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

// Get fetches a single file for given file id from Put.io API.
func (c *Client) Get(id int) (File, error) {
	if id < 0 {
		return File{}, fmt.Errorf("id cannot be negative")
	}

	req, err := c.NewRequest("GET", "/v2/files/"+strconv.Itoa(id), nil)
	if err != nil {
		return File{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return File{}, fmt.Errorf("get request failed with status: %v", resp.Status)
	}

	var getResponse struct {
		File   File   `json:"file"`
		Status string `json:"status"`
	}
	err = json.NewDecoder(resp.Body).Decode(&getResponse)
	if err != nil {
		return File{}, err
	}
	return getResponse.File, nil
}

// List fetches a list of files for given directory id from Put.io API.
func (c *Client) List(id int) (FileList, error) {
	if id < 0 {
		return FileList{}, fmt.Errorf("id cannot be negative")
	}
	req, err := c.NewRequest("GET", "/v2/files/list?parent_id="+strconv.Itoa(id), nil)
	if err != nil {
		return FileList{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return FileList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FileList{}, fmt.Errorf("list request failed. HTTP Status: %v", resp.Status)
	}

	var listResponse struct {
		Files  []File `json:"files"`
		Parent File   `json:"parent"`
		Status string `json:"status"`
	}

	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	if err != nil {
		return FileList{}, err
	}

	return FileList{
		Files:  listResponse.Files,
		Parent: listResponse.Parent,
	}, nil
}

// Download retrieves the download URL for the given file id.
func (c *Client) Download(id int) (string, error) {
	if id < 0 {
		return "", fmt.Errorf("id cannot be negative")
	}

	req, err := c.NewRequest("HEAD", "/v2/files/"+strconv.Itoa(id)+"/download?notunnel=1", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if urlErr, ok := err.(*url.Error); ok && urlErr.Err == errRedirectAttempt {
		err = nil
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("File (%v) could not be found", id)
	}

	downloadURL := resp.Header.Get("Location")
	if downloadURL == "" {
		return "", fmt.Errorf("location header is empty")
	}

	return downloadURL, nil
}

// createFolder creates a new folder under parent.
func (c *Client) createFolder(name string, parent int) error {
	panic("not implemented yet")
}

func (c *Client) search(query string, page int) ([]File, error) {
	panic("not implemented yet")
}

func (c *Client) upload(file, filename string, parent int) error {
	panic("not implemented yet")
}

func (c *Client) delete_(files ...int) error {
	if len(files) == 0 {
		return fmt.Errorf("no file id's are given")
	}
	panic("not implemented yet")
}

func (c *Client) rename(id int, name string) error {
	if id < 0 {
		return fmt.Errorf("id cannot be negative")
	}
	if name == "" {
		return fmt.Errorf("new filename cannot be empty")
	}
	panic("not implemented yet")
}

func (c *Client) move(parent int, files ...int) error {
	if parent < 0 {
		return fmt.Errorf("parent id cannot be negative")
	}
	if len(files) == 0 {
		return fmt.Errorf("no file id's are given")
	}
	panic("not implemented yet")
}

// File represents a Put.io file.
type File struct {
	ID                int    `json:"id"`
	Filename          string `json:"name"`
	Filesize          int64  `json:"size"`
	ContentType       string `json:"content_type"`
	CreatedAt         string `json:"created_at"`
	FirstAccessedAt   string `json:"first_accessed_at"`
	ParentID          int    `json:"parent_id"`
	Screenshot        string `json:"screenshot"`
	OpensubtitlesHash string `json:"opensubtitles_hash"`
	IsMP4Available    bool   `json:"is_mp4_available"`
	Icon              string `json:"icon"`
	CRC32             string `json:"crc32"`
}

// FileList represents a list of files of a Put.io directory.
type FileList struct {
	Files  []File `json:"file"`
	Parent File   `json:"parent"`
}
