package putio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
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
}

func TestFiles_CreateFolder(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/create-folder", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "file":{"id": 1,"name":"foo", "parent": 0}}`)
	})

	_, err := client.Files.CreateFolder("foo", 0)
	if err != nil {
		t.Error(err)
	}
}

func TestFiles_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Files.Delete(1, 2, 3)
	if err != nil {
		t.Error(err)
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
}

func TestFiles_Download(t *testing.T) {
	setup()
	defer teardown()

	fileContent := "this is the body of a file"
	mux.HandleFunc("/v2/files/1/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fileContent)
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
}
