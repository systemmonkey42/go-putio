package putio

import "encoding/json"

// EventsService is the service to gather information about user's events.
type EventsService struct {
	client *Client
}

// List gets list of dashboard events. It includes downloads and share events.
// FIXME: events list returns inconsistent data structures.
func (e *EventsService) List() ([]Event, error) {
	req, err := e.client.NewRequest("GET", "/v2/events/list", nil)
	if err != nil {
		return nil, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
		Events []Event
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Events, nil

}

// Delete Clears all all dashboard events.
func (e *EventsService) Delete() error {
	req, err := e.client.NewRequest("POST", "/v2/events/delete", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}

	// FIXME: return original error
	return json.NewDecoder(resp.Body).Decode(&r)
}
