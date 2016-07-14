package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestZips_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/zips/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "zips":[{"id": null,"created_at":"NULL"}]}`)
	})

	_, err := client.Zips.List()
	if err != nil {
		t.Error("zips.List() returned error: %v", err)
	}
}

func TestZips_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/zips/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "zip_id": "1"}`)
	})

	id, err := client.Zips.Create(666)
	if err != nil {
		t.Error("zips.Create() returned error: %v", err)
	}
	if id != 1 {
		t.Errorf("got: %v, want 1", id)
	}
}

func TestZips_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/zips/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "size": 123456, "url": "https://put.io"}`)
	})

	_, err := client.Zips.Get(1)
	if err != nil {
		t.Error("zips.Get() returned error: %v", err)
	}
}
