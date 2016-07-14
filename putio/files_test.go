package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFiles_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "file":{"id": 1,"name":"foo", "size":92}}`)
	})

	_, err := client.Files.Get(1)
	if err != nil {
		t.Error(err)
	}
}

func TestFiles_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "files":[{"id": 1,"name":"foo", "size":92}]}`)
	})

	_, err := client.Files.List(0)
	if err != nil {
		t.Error(err)
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
