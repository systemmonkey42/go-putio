package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestZips_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/zips/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "size": 123456, "url": "https://put.io"}`)
	})

	_, err := client.Zips.Get(1)
	if err != nil {
		t.Error(err)
	}
}

func TestZips_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"zips": [
		{
			"created_at": "2016-07-15T10:42:12",
			"id": 4177262
		}
	]
}
`
	mux.HandleFunc("/v2/zips/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	zips, err := client.Zips.List()
	if err != nil {
		t.Error(err)
	}

	if len(zips) != 1 {
		t.Errorf("got: %v, want: 1", len(zips))
	}

	if zips[0].ID != 4177262 {
		t.Errorf("got: %v, want: 4177262", zips[0].ID)
	}
}

func TestZips_Create(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"zip_id": 4177264
}
`
	mux.HandleFunc("/v2/zips/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, fixture)
	})

	id, err := client.Zips.Create(666)
	if err != nil {
		t.Error("zips.Create() returned error: %v", err)
	}

	if id != 4177264 {
		t.Errorf("got: %v, want 4177264", id)
	}
}
