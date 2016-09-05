package putio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FilesService is a general service to gather information about user files,
// such as listing, searching, creating new ones, or just fetching a single
// file.
type FilesService struct {
	client *Client
}

// Get fetches file metadata for given file ID.
func (f *FilesService) Get(ctx context.Context, id int) (File, error) {
	if id < 0 {
		return File{}, errNegativeID
	}

	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id), nil)
	if err != nil {
		return File{}, err
	}

	var r struct {
		File File `json:"file"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return File{}, err
	}
	return r.File, nil
}

// List fetches children for given directory ID.
func (f *FilesService) List(ctx context.Context, id int) ([]File, File, error) {
	if id < 0 {
		return nil, File{}, errNegativeID
	}
	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/list?parent_id="+strconv.Itoa(id), nil)
	if err != nil {
		return nil, File{}, err
	}

	var r struct {
		Files  []File `json:"files"`
		Parent File   `json:"parent"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, File{}, err
	}

	return r.Files, r.Parent, nil
}

// Download fetches the contents of the given file. Callers can pass additional
// useTunnel parameter to fetch the file from nearest tunnel server. Storage
// servers accept Range requests, so a range header can be provided by headers
// parameter.
//
// Download request is done by the client which is provided to the NewClient
// constructor. Additional client tunings are taken into consideration while
// downloading a file, such as Timeout etc.
func (f *FilesService) Download(ctx context.Context, id int, useTunnel bool, headers http.Header) (io.ReadCloser, error) {
	if id < 0 {
		return nil, errNegativeID
	}

	notunnel := "notunnel=1"
	if useTunnel {
		notunnel = "notunnel=0"
	}

	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id)+"/download?"+notunnel, nil)
	if err != nil {
		return nil, err
	}
	// merge headers with request headers
	for header, values := range headers {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	// follow the redirect only once. copy the original request headers to
	// redirect request.
	f.client.client.CheckRedirect = redirectOnceFunc
	defer func() {
		f.client.client.CheckRedirect = nil
	}()

	resp, err := f.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// CreateFolder creates a new folder under parent.
func (f *FilesService) CreateFolder(ctx context.Context, name string, parent int) (File, error) {
	if name == "" {
		return File{}, fmt.Errorf("empty folder name")
	}

	if parent < 0 {
		return File{}, errNegativeID
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("parent_id", strconv.Itoa(parent))

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/create-folder", strings.NewReader(params.Encode()))
	if err != nil {
		return File{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var r struct {
		File File `json:"file"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return File{}, err
	}

	return r.File, nil
}

// Delete deletes given files.
func (f *FilesService) Delete(ctx context.Context, files ...int) error {
	if len(files) == 0 {
		return fmt.Errorf("no file id is given")
	}

	var ids []string
	for _, id := range files {
		if id < 0 {
			return errNegativeID
		}
		ids = append(ids, strconv.Itoa(id))
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/delete", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}
	return nil
}

// Rename change the name of the file to newname.
func (f *FilesService) Rename(ctx context.Context, id int, newname string) error {
	if id < 0 {
		return errNegativeID
	}
	if newname == "" {
		return fmt.Errorf("new filename cannot be empty")
	}

	params := url.Values{}
	params.Set("file_id", strconv.Itoa(id))
	params.Set("name", newname)

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/rename", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}

	return nil
}

// Move moves files to the given destination.
func (f *FilesService) Move(ctx context.Context, parent int, files ...int) error {
	if parent < 0 {
		return errNegativeID
	}

	if len(files) == 0 {
		return fmt.Errorf("no files given")
	}

	var ids []string
	for _, file := range files {
		if file < 0 {
			return errNegativeID
		}
		ids = append(ids, strconv.Itoa(file))
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))
	params.Set("parent_id", strconv.Itoa(parent))

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/move", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}
	return nil
}

// Upload reads from given io.Reader and uploads the file contents to Put.io
// servers under directory given by parent. If parent is negative, user's
// prefered folder is used.
//
// If the uploaded file is a torrent file, Put.io will interpret it as a
// transfer and Transfer field will be present to represent the status of the
// tranfer. Likewise, if the uploaded file is a regular file, Transfer field
// would be nil and the uploaded file will be represented by the File field.
//
// This method reads the file contents into the memory, so it should be used for
// <150MB files.
func (f *FilesService) Upload(ctx context.Context, r io.Reader, filename string, parent int) (Upload, error) {
	if filename == "" {
		return Upload{}, fmt.Errorf("filename cannot be empty")
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// negative parent means use user's prefered download folder.
	if parent >= 0 {
		err := mw.WriteField("parent_id", strconv.Itoa(parent))
		if err != nil {
			return Upload{}, err
		}
	}

	formfile, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return Upload{}, err
	}

	_, err = io.Copy(formfile, r)
	if err != nil {
		return Upload{}, err
	}

	err = mw.Close()
	if err != nil {
		return Upload{}, err
	}

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/upload", &buf)
	if err != nil {
		return Upload{}, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	var response struct {
		Upload
	}
	_, err = f.client.Do(req, &response)
	if err != nil {
		return Upload{}, err
	}
	return response.Upload, nil
}

// Search makes a search request with the given query. Servers return 50
// results at a time. The URL for the next 50 results are in Next field.  If
// page is -1, all results are returned.
func (f *FilesService) Search(ctx context.Context, query string, page int) (Search, error) {
	if page == 0 || page < -1 {
		return Search{}, fmt.Errorf("invalid page number")
	}
	if query == "" {
		return Search{}, fmt.Errorf("no query given")
	}

	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/search/"+query+"/page/"+strconv.Itoa(page), nil)
	if err != nil {
		return Search{}, err
	}

	var r Search
	_, err = f.client.Do(req, &r)
	if err != nil {
		return Search{}, err
	}

	return r, nil
}

// FIXME: is it worth to export this method?
func (f *FilesService) convert(ctx context.Context, id int) error {
	if id < 0 {
		return errNegativeID
	}

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/"+strconv.Itoa(id)+"/mp4", nil)
	if err != nil {
		return err
	}

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}

	return nil
}

// Share shares given files with given friends. Friends are list of usernames.
// If no friends is given, files are shared with all of your friends.
func (f *FilesService) share(ctx context.Context, files []int, friends ...string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files given")
	}

	var ids []string
	for _, file := range files {
		if file < 0 {
			return errNegativeID
		}
		ids = append(ids, strconv.Itoa(file))
	}

	var friendsParam string
	if len(friends) == 0 {
		friendsParam = "everyone"
	} else {
		friendsParam = strings.Join(friends, ",")
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))
	params.Set("friends", friendsParam)

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/share", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}

	return nil
}

// Shared returns list of shared files and share information.
func (f *FilesService) shared(ctx context.Context) ([]share, error) {
	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/shared", nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Shared []share
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Shared, nil
}

// SharedWith returns list of users the given file is shared with.
func (f *FilesService) sharedWith(ctx context.Context, id int) ([]share, error) {
	if id < 0 {
		return nil, errNegativeID
	}

	// FIXME: shared-with returns different json structure than /shared/
	// endpoint. so it's not an exported method until a common structure is
	// decided
	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id)+"/shared-with", nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Shared []share `json:"shared-with"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Shared, nil
}

// Subtitles lists available subtitles for the given file for user's prefered
// subtitle language.
func (f *FilesService) Subtitles(ctx context.Context, id int) ([]Subtitle, error) {
	if id < 0 {
		return nil, errNegativeID
	}

	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id)+"/subtitles", nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Subtitles []Subtitle
		Default   string
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Subtitles, nil
}

// DownloadSubtitle sends the contents of the subtitle file. If the key is empty string,
// `default` key is used. This key is used to search for a subtitle in the
// following order and returns the first match:
// - A subtitle file that has identical parent folder and name with the video.
// - Subtitle file extracted from video if the format is MKV.
// - First match from OpenSubtitles.org.
func (f *FilesService) DownloadSubtitle(ctx context.Context, id int, key string, format string) (io.ReadCloser, error) {
	if id < 0 {
		return nil, errNegativeID
	}

	if key == "" {
		key = "default"
	}
	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id)+"/subtitles/"+key, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// HLSPlaylist serves a HLS playlist for a video file. Use “all” as
// subtitleKey to get available subtitles for user’s preferred languages.
func (f *FilesService) HLSPlaylist(ctx context.Context, id int, subtitleKey string) (io.ReadCloser, error) {
	if id < 0 {
		return nil, errNegativeID
	}

	if subtitleKey == "" {
		return nil, fmt.Errorf("empty subtitle key is given")
	}

	req, err := f.client.NewRequest(ctx, "GET", "/v2/files/"+strconv.Itoa(id)+"/hls/media.m3u8?subtitle_key"+subtitleKey, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// SetVideoPosition sets default video position for a video file.
func (f *FilesService) SetVideoPosition(ctx context.Context, id int, t int) error {
	if id < 0 {
		return errNegativeID
	}

	if t < 0 {
		return fmt.Errorf("time cannot be negative")
	}

	params := url.Values{}
	params.Set("time", strconv.Itoa(t))

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/"+strconv.Itoa(id)+"/start-from", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}

	return nil
}

// DeleteVideoPosition deletes video position for a video file.
func (f *FilesService) DeleteVideoPosition(ctx context.Context, id int) error {
	if id < 0 {
		return errNegativeID
	}

	req, err := f.client.NewRequest(ctx, "POST", "/v2/files/"+strconv.Itoa(id)+"/start-from/delete", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = f.client.Do(req, &struct{}{})
	if err != nil {
		return err
	}

	return nil
}
