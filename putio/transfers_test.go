package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTransfers_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "transfers":[{"file_id": null,"name":"foo"}]}`)
	})

	_, err := client.Transfers.List()
	if err != nil {
		t.Error(err)
	}
}

func TestTransfers_Add(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "transfer":{"file_id": null,"name":"foo"}}`)
	})

	_, err := client.Transfers.Add("filepath", 0, false, "")
	if err != nil {
		t.Error(err)
	}
}

func TestTransfers_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "transfers":[{"file_id": null,"name":"foo"}]}`)
	})

	_, err := client.Transfers.Get(1)
	if err != nil {
		t.Error(err)
	}
}

func TestTransfers_Cancel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Transfers.Cancel(1)
	if err != nil {
		t.Error(err)
	}
}

func TestTransfers_Clean(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/clean", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Transfers.Clean()
	if err != nil {
		t.Error(err)
	}
}
