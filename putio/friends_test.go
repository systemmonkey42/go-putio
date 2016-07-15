package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFriends_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	_, err := client.Friends.List()
	if err != nil {
		t.Error(err)
	}
}

func TestFriends_WaitingRequests(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/waiting-requests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	_, err := client.Friends.WaitingRequests()
	if err != nil {
		t.Error(err)
	}
}

func TestFriends_Request(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/naber/request", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	err := client.Friends.Request("naber")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Request("")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Approve(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/naber/approve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	err := client.Friends.Approve("naber")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Approve("")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Deny(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/naber/deny", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	err := client.Friends.Deny("naber")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Deny("")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Unfriend(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/naber/unfriend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	err := client.Friends.Unfriend("naber")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Unfriend("")
	if err == nil {
		t.Error("empty username accepted")
	}
}
