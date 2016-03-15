package putio

import (
	"golang.org/x/oauth2"
	"net/http"
)

// NewAuthHelper returns an oauth2 enabled http client.
func NewAuthHelper(token string) *http.Client {
	return oauth2.NewClient(oauth2.NoContext, tokenSource(token))
}

// tokenSource implements the oauth2.TokenSource interface.
type tokenSource string

func (t tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: string(t)}, nil
}
