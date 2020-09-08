package putio

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const uploadURL = "https://upload.put.io/files/"

// Uploader uses TUS upload protocol for sending files to put.io.
type Uploader struct {
	client  *http.Client
	timeout time.Duration
	token   string

	Log func(message string)
}

// NewUploader returns a new Uploader for given http.Client and request timeout.
func NewUploader(client *http.Client, timeout time.Duration, token string) *Uploader {
	return &Uploader{
		client:  client,
		timeout: timeout,
		token:   token,
	}
}

func (u *Uploader) log(message string) {
	if u.Log != nil {
		u.Log(message)
	}
}

// CreateUpload is used for beginning new upload. Use returned location in SendFile function.
func (u *Uploader) CreateUpload(ctx context.Context, filename string, parentID, length int64) (location string, err error) {
	u.log(fmt.Sprintf("Creating upload %q at parent=%d", filename, parentID))
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, nil)
	if err != nil {
		return
	}
	metadata := map[string]string{
		"name":       filename,
		"parent_id":  strconv.FormatInt(parentID, 10),
		"no-torrent": "true",
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(length, 10))
	req.Header.Set("Upload-Metadata", encodeMetadata(metadata))
	req.Header.Set("Authorization", "token "+u.token)

	resp, err := u.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	location = resp.Header.Get("Location")
	return
}

// SendFile sends the contents of the file to put.io.
// In case of an transmission error, you can resume upload but you have to get the correct offset from server by
// calling GetOffset and must seek to the new offset on io.Reader.
func (u *Uploader) SendFile(ctx context.Context, r io.Reader, location string, offset int64) (fileID int64, crc32 string, err error) {
	u.log(fmt.Sprintf("Sending file %q offset=%d", location, offset))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stop upload if speed is too slow.
	// Wrap reader so each read call resets the timer that cancels the request on certain duration.
	r = &timerResetReader{r: r, timer: time.AfterFunc(u.timeout, cancel), timeout: u.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, location, r)
	if err != nil {
		return
	}

	req.Header.Set("content-type", "application/offset+octet-stream")
	req.Header.Set("upload-offset", strconv.FormatInt(offset, 10))
	req.Header.Set("Authorization", "token "+u.token)
	resp, err := u.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	fileID, err = strconv.ParseInt(resp.Header.Get("putio-file-id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse putio-file-id header: %w", err)
		return
	}
	crc32 = resp.Header.Get("putio-file-crc32")
	return
}

// GetOffset returns the offset at the server.
func (u *Uploader) GetOffset(ctx context.Context, location string) (n int64, err error) {
	u.log(fmt.Sprintf("Getting upload offset %q", location))
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, location, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "token "+u.token)
	resp, err := u.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	n, err = strconv.ParseInt(resp.Header.Get("upload-offset"), 10, 64)
	u.log(fmt.Sprintln("uploadJob offset:", n))
	return n, err
}

// TerminateUpload removes incomplete file from the server.
func (u *Uploader) TerminateUpload(ctx context.Context, location string) (err error) {
	u.log(fmt.Sprintf("Terminating upload %q", location))
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, location, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "token "+u.token)
	resp, err := u.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	return nil
}

func encodeMetadata(metadata map[string]string) string {
	encoded := make([]string, 0, len(metadata))
	for k, v := range metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, base64.StdEncoding.EncodeToString([]byte(v))))
	}
	return strings.Join(encoded, ",")
}

type timerResetReader struct {
	r       io.Reader
	timer   *time.Timer
	timeout time.Duration
}

func (r *timerResetReader) Read(p []byte) (int, error) {
	r.timer.Reset(r.timeout)
	return r.r.Read(p)
}
