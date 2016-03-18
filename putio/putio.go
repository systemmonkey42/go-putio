// Package putio is the Put.io API v2 client for Go.
package putio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultUserAgent = "go-putio"
	defaultMediaType = "application/json"
	defaultBaseURL   = "https://api.put.io"
	defaultUploadURL = "https://upload.put.io"
)

var errRedirect = fmt.Errorf("redirect attempt on a no-redirect client")

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

// Download retrieves the download URL for the given file id. Callers can pass
// additional useTunnel parameter to fetch the file from the nearest tunnel
// server, not from the main storage servers.
func (c *Client) Download(id int, useTunnel bool) (string, error) {
	if id < 0 {
		return "", fmt.Errorf("id cannot be negative")
	}

	notunnel := "notunnel=1"
	if useTunnel {
		notunnel = "notunnel=0"
	}
	req, err := c.NewRequest("HEAD", "/v2/files/"+strconv.Itoa(id)+"/download?"+notunnel, nil)
	if err != nil {
		return "", err
	}

	// our HTTP client follows redirect by default but file download URL is in
	// the first requests Location header, and this header exists on the first
	// request.
	c.client.CheckRedirect = noRedirectFunc
	defer func() {
		c.client.CheckRedirect = nil
	}()

	resp, err := c.client.Do(req)
	if urlErr, ok := err.(*url.Error); ok && urlErr.Err == errRedirect {
		err = nil
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("file could not be found. File id: %v, Status: %v", id, resp.Status)
	}

	downloadURL := resp.Header.Get("Location")
	if downloadURL == "" {
		return "", fmt.Errorf("could not retrieve download URL")
	}

	return downloadURL, nil
}

// CreateFolder creates a new folder under parent.
func (c *Client) CreateFolder(name string, parent int) error {
	if parent < 0 {
		return fmt.Errorf("parent id cannot be negative")
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("parent_id", strconv.Itoa(parent))

	req, err := c.NewRequest("POST", "/v2/files/create-folder", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create-folder request failed. HTTP Status: %v", resp.Status)
	}

	return nil
}

// Delete deletes given files.
func (c *Client) Delete(files ...int) error {
	if len(files) == 0 {
		return fmt.Errorf("no file id's are given")
	}

	var ids []string
	for _, id := range files {
		if id < 0 {
			return fmt.Errorf("file id cannot be negative")
		}
		ids = append(ids, strconv.Itoa(id))
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))

	req, err := c.NewRequest("POST", "/v2/files/delete", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete request failed. HTTP Status: %v", resp.Status)
	}

	return nil
}

// Rename renames the file to name for the given file id.
func (c *Client) Rename(id int, name string) error {
	if id < 0 {
		return fmt.Errorf("id cannot be negative")
	}
	if name == "" {
		return fmt.Errorf("new filename cannot be empty")
	}

	params := url.Values{}
	params.Set("file_id", strconv.Itoa(id))
	params.Set("name", name)

	req, err := c.NewRequest("POST", "/v2/files/rename", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rename request failed. HTTP Status: %v", resp.Status)
	}

	return nil
}

// Move moves files to the given destination.
func (c *Client) Move(parent int, files ...int) error {
	if len(files) == 0 {
		return fmt.Errorf("no file id's are given")
	}

	var ids []string
	for _, id := range files {
		if id < 0 {
			return fmt.Errorf("file id cannot be negative")
		}
		ids = append(ids, strconv.Itoa(id))
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))
	params.Set("parent", strconv.Itoa(parent))

	req, err := c.NewRequest("POST", "/v2/files/move", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("move request failed. HTTP Status: %v", resp.Status)
	}

	return nil
}

// Upload reads from filepath and uploads the file contents to Put.io servers
// under the parent directory with the name filename. This method reads the
// file contents into the memory, so it should be used for <150MB files.
func (c *Client) Upload(filepath, filename string, parent int) error {
	if parent < 0 {
		return fmt.Errorf("parent id cannot be negative")
	}

	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	err = mw.WriteField("parent_id", strconv.Itoa(parent))
	if err != nil {
		return err
	}

	formfile, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(formfile, f)
	if err != nil {
		return err
	}

	err = mw.Close()
	if err != nil {
		return err
	}

	u, _ := url.Parse(defaultUploadURL)
	c.BaseURL = u

	req, err := c.NewRequest("POST", "/v2/files/upload", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return fmt.Errorf("upload request failed with HTTP status: %v", resp.Status)
	}

	return nil
}

func (c *Client) search(query string, page int) ([]File, error) {
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

// noRedirectFunc prevents http client to follow redirects. This is needed for
// Put.io Download method to grab the download URL of a file.
func noRedirectFunc(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}
	return errRedirect
}
