package putio

import (
	"encoding/json"
)

// AccountService is the service to gather information about user account.
type AccountService struct {
	client *Client
}

// Info retrieves user account information.
func (a *AccountService) Info() (Info, error) {
	req, err := a.client.NewRequest("GET", "/v2/account/info", nil)
	if err != nil {
		return Info{}, nil
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()

	var r struct {
		Info   Info
		Status string
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Info{}, err
	}

	return r.Info, nil
}

// Settings retrieves user preferences.
func (a *AccountService) Settings() (Settings, error) {
	req, err := a.client.NewRequest("GET", "/v2/account/settings", nil)
	if err != nil {
		return Settings{}, nil
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return Settings{}, err
	}
	defer resp.Body.Close()

	var r struct {
		Settings Settings
		Status   string
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Settings{}, err
	}

	return r.Settings, nil
}

// FIXME: fill
func (a *AccountService) UpdateSettings() error {
	panic("not implemented yet")
}
