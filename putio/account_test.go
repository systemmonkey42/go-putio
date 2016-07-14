package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAccount_Info(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"info": {
	"username": "naber",
        "mail": "naber@iyidir.com",
        "plan_expiration_date": "2014-03-04T06:33:30",
        "subtitle_languages": ["tr", "eng"],
        "default_subtitle_language": "tr",
        "disk": {
            "avail": 20849243836,
            "used": 32837847364,
            "size": 53687091200
        }
    },
    "status": "OK"
}`

	mux.HandleFunc("/v2/account/info", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	info, err := client.Account.Info()
	if err != nil {
		t.Error("account.Info() returned error: %v", err)
	}

	if info.Username != "naber" {
		t.Errorf("got: %v, want: naber", info.Username)
	}

	if info.Mail != "naber@iyidir.com" {
		t.Errorf("got: %v, want: naber@iyidir.com", info.Mail)
	}
}

func TestAccount_Settings(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
    "status": "OK",
    "settings": {
        "default_download_folder": 0,
        "is_invisible": false,
        "subtitle_languages": ["tr", "eng"],
        "default_subtitle_language": "tr"
    }
}
`
	mux.HandleFunc("/v2/account/settings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	settings, err := client.Account.Settings()
	if err != nil {
		t.Error(err)
	}

	if settings.DefaultDownloadFolder != 0 {
		t.Errorf("got: %v, want: 0", settings.DefaultDownloadFolder)
	}

	if settings.DefaultSubtitleLanguage != "tr" {
		t.Errorf("got: %v, want: tr", settings.DefaultSubtitleLanguage)
	}
}
