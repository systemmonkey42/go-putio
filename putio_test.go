package putio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("got: %v, want: %v", r.Method, want)
	}
}

func TestNewClient(t *testing.T) {
	client, _ = NewClient(nil)
	if client.BaseURL.String() != defaultBaseURL {
		t.Errorf("got: %v, want: %v", client.BaseURL.String(), defaultBaseURL)
	}
}

func TestFiles_Get(t *testing.T) {
	setup()
	defer teardown()

	var (
		fileStr = `{
"file":
	{
		"content_type": "text/plain",
		"crc32": "66a1512f",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": "2013-09-07T21:32:03",
		"icon": "https://put.io/images/file_types/text.png",
		"id": 6546533,
		"is_mp4_available": false,
		"name": "MyFile.txt",
		"opensubtitles_hash": "",
		"parent_id": 123,
		"screenshot": "",
		"size": 92
	},
"status": "OK"
}`

		fileJSON = File{
			ContentType:       "text/plain",
			CRC32:             "66a1512f",
			CreatedAt:         "2013-09-07T21:32:03",
			FirstAccessedAt:   "2013-09-07T21:32:03",
			Icon:              "https://put.io/images/file_types/text.png",
			ID:                6546533,
			IsMP4Available:    false,
			Filename:          "MyFile.txt",
			OpensubtitlesHash: "",
			ParentID:          123,
			Screenshot:        "",
			Filesize:          92,
		}
	)

	mux.HandleFunc("/v2/files/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fileStr)
	})

	file, err := client.Get(0)
	if err != nil {
		t.Errorf("get returned error: %v", err)
	}
	if !reflect.DeepEqual(file, fileJSON) {
		t.Errorf("got: %v, want: %v", file, fileJSON)
	}
}

func TestFiles_List(t *testing.T) {
	setup()
	defer teardown()

	var (
		filelistStr = `
{
    "files": [
        {
            "content_type": "text/plain",
            "crc32": "66a1512f",
            "created_at": "2013-09-07T21:32:03",
            "first_accessed_at": "2013-09-07T21:32:03",
            "icon": "https://put.io/images/file_types/text.png",
            "id": 6546533,
            "is_mp4_available": false,
            "name": "MyFile.txt",
            "parent_id": 123,
            "size": 92
        },
        {
            "content_type": "video/x-matroska",
            "crc32": "cb97ba70",
            "created_at": "2013-09-07T21:32:03",
            "first_accessed_at": "2013-09-07T21:32:03",
            "icon": "https://put.io/thumbnails/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
            "id": 7645645,
            "is_mp4_available": false,
            "name": "MyVideo.mkv",
            "parent_id": 123,
            "size": 1155197659
        }
    ],
    "parent": {
        "content_type": "application/x-directory",
        "crc32": "foo",
        "created_at": "2013-09-07T21:32:03",
        "first_accessed_at": "2013-09-07T21:32:03",
        "icon": "https://put.io/images/file_types/folder.png",
        "id": 123,
        "is_mp4_available": false,
        "name": "MyFolder",
        "parent_id": 0,
        "size": 1155197751
        },
    "status": "OK"
}
`

		filelistJSON = FileList{
			Files: []File{
				{
					ContentType:     "text/plain",
					CRC32:           "66a1512f",
					CreatedAt:       "2013-09-07T21:32:03",
					FirstAccessedAt: "2013-09-07T21:32:03",
					Icon:            "https://put.io/images/file_types/text.png",
					ID:              6546533,
					IsMP4Available:  false,
					Filename:        "MyFile.txt",
					ParentID:        123,
					Filesize:        92,
				},
				{
					ContentType:     "video/x-matroska",
					CRC32:           "cb97ba70",
					CreatedAt:       "2013-09-07T21:32:03",
					FirstAccessedAt: "2013-09-07T21:32:03",
					Icon:            "https://put.io/thumbnails/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
					ID:              7645645,
					IsMP4Available:  false,
					Filename:        "MyVideo.mkv",
					ParentID:        123,
					Filesize:        1155197659,
				},
			},
			Parent: File{
				ContentType:     "application/x-directory",
				CRC32:           "foo",
				CreatedAt:       "2013-09-07T21:32:03",
				FirstAccessedAt: "2013-09-07T21:32:03",
				Icon:            "https://put.io/images/file_types/folder.png",
				ID:              123,
				IsMP4Available:  false,
				Filename:        "MyFolder",
				ParentID:        0,
				Filesize:        1155197751,
			},
		}
	)

	mux.HandleFunc("/v2/files/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, filelistStr)
	})

	filelist, err := client.List(0)
	if err != nil {
		t.Errorf("list returned error: %v", err)
	}
	if !reflect.DeepEqual(filelist, filelistJSON) {
		t.Errorf("got: %v, want: %v", filelist, filelistJSON)
	}
}
