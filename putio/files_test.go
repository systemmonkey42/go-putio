package putio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestFiles_Get(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"file": {
		"content_type": "text/plain",
		"crc32": "66a1512f",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": null,
		"icon": "https://put.io/images/file_types/text.png",
		"id": 6546533,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyFile.txt",
		"opensubtitles_hash": null,
		"parent_id": 123,
		"screenshot": null,
		"size": 92
	},
    "status": "OK"
}
`

	mux.HandleFunc("/v2/files/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	file, err := client.Files.Get(1)
	if err != nil {
		t.Error(err)
	}

	if file.Filesize != 92 {
		t.Errorf("got: %v, want: 92", file.Filesize)
	}

	// negative id
	_, err = client.Files.Get(-1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"files": [
	{
		"content_type": "text/plain",
		"crc32": "66a1512f",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": null,
		"icon": "https://put.io/images/file_types/text.png",
		"id": 6546533,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyFile.txt",
		"opensubtitles_hash": null,
		"parent_id": 123,
		"screenshot": null,
		"size": 92
	},
	{
		"content_type": "video/x-matroska",
		"crc32": "cb97ba70",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": "2013-09-07T21:32:13",
		"icon": "https://put.io/thumbnails/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
		"id": 7645645,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyVideo.mkv",
		"opensubtitles_hash": "acc2785ffa573c69",
		"parent_id": 123,
		"screenshot": "https://put.io/screenshots/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
		"size": 1155197659
	}
],
"parent": {
	"content_type": "application/x-directory",
	"crc32": null,
	"created_at": "2013-09-07T21:32:03",
	"first_accessed_at": null,
	"icon": "https://put.io/images/file_types/folder.png",
	"id": 123,
	"is_mp4_available": false,
	"is_shared": false,
	"name": "MyFolder",
	"opensubtitles_hash": null,
	"parent_id": 0,
	"screenshot": null,
	"size": 1155197751
},
"status": "OK"
}
`
	mux.HandleFunc("/v2/files/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	files, parent, err := client.Files.List(0)
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Errorf("got: %v, want: 2", len(files))
	}
	if parent.ID != 123 {
		t.Errorf("got: %v, want: 123", parent.ID)
	}

	// negative id
	_, _, err = client.Files.List(-1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_CreateFolder(t *testing.T) {
	setup()
	defer teardown()

	fixture := `

{
	"file": {
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2016-07-15T09:21:03",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://api.put.io/images/file_types/folder.png",
		"id": 415105276,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "foobar",
		"opensubtitles_hash": null,
		"parent_id": 0,
		"screenshot": null,
		"size": 0
	},
	"status": "OK"
}
`
	mux.HandleFunc("/v2/files/create-folder", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, fixture)
	})

	file, err := client.Files.CreateFolder("foobar", 0)
	if err != nil {
		t.Error(err)
	}

	if file.Filename != "foobar" {
		t.Errorf("got: %v, want: foobar", file.Filename)
	}

	// empty folder name
	_, err = client.Files.CreateFolder("", 0)
	if err == nil {
		t.Errorf("empty folder name accepted")
	}

	// negative id
	_, err = client.Files.CreateFolder("foobar", -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status": "OK"}`)
	})

	err := client.Files.Delete(1, 2, 3)
	if err != nil {
		t.Error(err)
	}

	// empty params
	err = client.Files.Delete()
	if err == nil {
		t.Errorf("empty parameters accepted")
	}

	err = client.Files.Delete(1, 2, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_Rename(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/rename", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Files.Rename(1, "bar")
	if err != nil {
		t.Error(err)
	}

	// negative id
	err = client.Files.Rename(-1, "bar")
	if err == nil {
		t.Errorf("negative file ID accepted")
	}

	// empty name
	err = client.Files.Rename(1, "")
	if err == nil {
		t.Errorf("empty filename accepted")
	}
}

func TestFiles_Download(t *testing.T) {
	setup()
	defer teardown()

	fileContent := "0123456789"
	mux.HandleFunc("/v2/files/1/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		buf := strings.NewReader(fileContent)
		http.ServeContent(w, r, "testfile", time.Now().UTC(), buf)
	})

	rc, err := client.Files.Download(1, false, nil)
	if err != nil {
		t.Error(err)
	}
	defer rc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}

	if buf.String() != fileContent {
		t.Errorf("got: %q, want: %q", buf.String(), fileContent)
	}

	// negative id
	_, err = client.Files.Download(-1, false, nil)
	if err == nil {
		t.Errorf("negative id accepted")
	}

	// range request
	rangeHeader := http.Header{}
	rangeHeader.Set("Range", fmt.Sprintf("bytes=%v-%v", 0, 3))
	rc, err = client.Files.Download(1, false, rangeHeader)
	if err != nil {
		t.Error(err)
	}
	defer rc.Close()

	buf.Reset()
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}

	response := buf.String()
	if response != "0123" {
		t.Errorf("got: %v, want: 0123", response)
	}
}

func TestFiles_Search(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"files": [
	{
		"content_type": "video/x-msvideo",
		"crc32": "812ed74d",
		"created_at": "2013-04-30T21:40:04",
		"extension": "avi",
		"file_type": "VIDEO",
		"first_accessed_at": "2013-12-24T09:18:58",
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 79905833,
		"is_hidden": false,
		"is_mp4_available": true,
		"is_shared": false,
		"name": "some-file.mkv",
		"opensubtitles_hash": "fb5414fd9b9e1e38",
		"parent_id": 79905827,
		"screenshot": "https://some-valid-screenhost-url.com",
		"sender_name": "hafifuyku",
		"size": 738705408,
		"start_from": 0
	},
	{
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2013-04-30T21:40:03",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 79905827,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "Movie 43",
		"opensubtitles_hash": null,
		"parent_id": 2197,
		"screenshot": null,
		"sender_name": "hafifuyku",
		"size": 738831202
	},
	{
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2010-05-19T22:24:21",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 5659875,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MOVIE",
		"opensubtitles_hash": null,
		"parent_id": 0,
		"screenshot": null,
		"sender_name": "emsel",
		"size": 0
	}
],
"next": null,
"status": "OK",
"total": 3
}
`
	mux.HandleFunc("/v2/files/search/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	s, err := client.Files.Search("naber", 1)
	if err != nil {
		t.Error(err)
	}

	if len(s.Files) != 3 {
		t.Errorf("got: %v, want: 3", len(s.Files))
	}

	if s.Files[0].Filename != "some-file.mkv" {
		t.Errorf("got: %v, want: some-file.mkv", s.Files[0].Filename)
	}

	// invalid page number
	_, err = client.Files.Search("naber", 0)
	if err == nil {
		t.Errorf("invalid page number accepted")
	}

	// empty query
	_, err = client.Files.Search("", 1)
	if err == nil {
		t.Errorf("empty query accepted")
	}
}
