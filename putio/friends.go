package putio

import (
	"encoding/json"
	"fmt"
)

// FriendsService is the service to operate on user friends.
type FriendsService struct {
	client *Client
}

// List lists users friends.
func (f *FriendsService) List() ([]Friend, error) {
	req, err := f.client.NewRequest("GET", "/v2/friends/list", nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Friends []Friend
		Status  string
		Total   int
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.Friends, nil
}

// WaitingRequests lists user's pending friend requests.
func (f *FriendsService) WaitingRequests() ([]Friend, error) {
	req, err := f.client.NewRequest("GET", "/v2/friends/waiting-requests", nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Friends []Friend
		Status  string
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.Friends, nil
}

// Request sends a friend request to the given username.
func (f *FriendsService) Request(username string) error {
	if username == "" {
		return fmt.Errorf("empty username")
	}
	req, err := f.client.NewRequest("POST", "/v2/friends/"+username+"/request", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}

	return json.NewDecoder(resp.Body).Decode(&r)
}

// Approve approves a friend request from the given username.
func (f *FriendsService) Approve(username string) error {
	if username == "" {
		return fmt.Errorf("empty username")
	}

	req, err := f.client.NewRequest("POST", "/v2/friends/"+username+"/approve", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}

	return json.NewDecoder(resp.Body).Decode(&r)
}

// Deny denies a friend request from the given username.
func (f *FriendsService) Deny(username string) error {
	if username == "" {
		return fmt.Errorf("empty username")
	}

	req, err := f.client.NewRequest("POST", "/v2/friends/"+username+"/deny", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}

	return json.NewDecoder(resp.Body).Decode(&r)
}

// Unfriend removed friend from user's friend list.
func (f *FriendsService) Unfriend(username string) error {
	if username == "" {
		return fmt.Errorf("empty username")
	}

	req, err := f.client.NewRequest("POST", "/v2/friends/"+username+"/unfriend", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Status string
	}

	return json.NewDecoder(resp.Body).Decode(&r)
}
