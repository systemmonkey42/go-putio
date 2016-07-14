package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAccount_Info(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/account/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"OK", "info":{"user_id": 1}}`)
	})

	_, err := client.Account.Info()
	if err != nil {
		t.Error("account.Info() returned error: %v", err)
	}
}

func TestAccount_Settings(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/account/settings", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"OK", "settings":{"routing": "Istanbul"}}`)
	})

	_, err := client.Account.Settings()
	if err != nil {
		t.Error("account.Settings() returned error: %v", err)
	}
}
