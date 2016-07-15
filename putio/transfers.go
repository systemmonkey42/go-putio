package putio

import (
	"encoding/json"
	"fmt"
	urlpkg "net/url"
	"strconv"
	"strings"
)

// TransfersService is the service to operate on torrent transfers, such as
// adding a torrent or magnet link, retrying a current one etc.
type TransfersService struct {
	client *Client
}

// List lists all active transfers. If a transfer is completed, it will not be
// available in response.
func (t *TransfersService) List() ([]Transfer, error) {
	req, err := t.client.NewRequest("GET", "/v2/transfers/list", nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Transfers []Transfer
		Status    string
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.Transfers, nil
}

// Add creates a new transfer.
// Add expects a valid torrent/magnet URL. Other parameters are optional.
//
// FIXME: change signateure?
func (t *TransfersService) Add(url string, parent int, extract bool, callbackURL string) (Transfer, error) {
	if parent < 0 {
		return Transfer{}, errNegativeID
	}

	if url == "" {
		return Transfer{}, fmt.Errorf("empty URL")
	}

	params := urlpkg.Values{}
	params.Set("url", url)
	params.Set("parent", strconv.Itoa(parent))
	params.Set("extract", strconv.FormatBool(extract))
	params.Set("callback_url", callbackURL)

	req, err := t.client.NewRequest("POST", "/v2/transfers/add", strings.NewReader(params.Encode()))
	if err != nil {
		return Transfer{}, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return Transfer{}, err
	}
	defer resp.Body.Close()

	var r struct {
		Transfer Transfer
		Status   string
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Transfer{}, err
	}
	return r.Transfer, nil
}

// Get returns the given transfer's properties.
func (t *TransfersService) Get(id int) (Transfer, error) {
	if id < 0 {
		return Transfer{}, errNegativeID
	}

	req, err := t.client.NewRequest("GET", "/v2/transfers/"+strconv.Itoa(id), nil)
	if err != nil {
		return Transfer{}, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return Transfer{}, err
	}
	defer resp.Body.Close()

	var r struct {
		Transfer Transfer
		Status   string
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Transfer{}, err
	}
	return r.Transfer, nil
}

// FIXME: fill
func (t *TransfersService) Retry(id int) {
}

// Cancel deletes given transfers.
func (t *TransfersService) Cancel(ids ...int) error {
	if len(ids) == 0 {
		return fmt.Errorf("no id given")
	}

	var transfers []string
	for _, id := range ids {
		if id < 0 {
			return errNegativeID
		}
		transfers = append(transfers, strconv.Itoa(id))
	}

	params := urlpkg.Values{}
	params.Set("transfer_ids", strings.Join(transfers, ","))
	req, err := t.client.NewRequest("POST", "/v2/transfers/cancel", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}
	if r.Status == "OK" {
		return nil
	}

	// FIXME: send the actual error
	return fmt.Errorf("err")
}

// Clean removes completed transfers from the transfer list.
func (t *TransfersService) Clean() error {
	req, err := t.client.NewRequest("POST", "/v2/transfers/clean", nil)
	if err != nil {
		return err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}
	if r.Status == "OK" {
		return nil
	}

	// FIXME: send the actual error
	return fmt.Errorf("err")
}
