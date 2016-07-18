package putio

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

	client = NewClient(nil)
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

func testHeader(t *testing.T, r *http.Request, key, value string) {
	if r.Header.Get(key) != value {
		t.Errorf("missing header. want: %q: %q", key, value)
	}
}

func TestNewClient(t *testing.T) {
	client = NewClient(nil)
	if client.BaseURL.String() != defaultBaseURL {
		t.Errorf("got: %v, want: %v", client.BaseURL.String(), defaultBaseURL)
	}
}
