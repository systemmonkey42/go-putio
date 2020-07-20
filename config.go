package putio

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ConfigService struct {
	client *Client
}

func (f *ConfigService) GetAll(ctx context.Context) (map[string]interface{}, error) {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config", nil)
	if err != nil {
		return nil, err
	}
	var r = struct {
		Config map[string]interface{} `json:"config"`
	}{
		Config: make(map[string]interface{}),
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, err
	}
	return r.Config, nil
}

func (f *ConfigService) Get(ctx context.Context, key string) (interface{}, error) {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config/"+key, nil)
	if err != nil {
		return nil, err
	}
	var r struct {
		Value interface{} `json:"value"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return nil, err
	}
	return r.Value, nil
}

func (f *ConfigService) SetAll(ctx context.Context, config map[string]interface{}) error {
	v := struct {
		Config interface{} `json:"config"`
	}{
		Config: config,
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	req, err := f.client.NewRequest(ctx, http.MethodPut, "/v2/config", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	_, err = f.client.Do(req, nil)
	return err
}

func (f *ConfigService) Set(ctx context.Context, key string, value interface{}) error {
	v := struct {
		Value interface{} `json:"value"`
	}{
		Value: value,
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	req, err := f.client.NewRequest(ctx, http.MethodPut, "/v2/config/"+key, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	_, err = f.client.Do(req, nil)
	return err
}

func (f *ConfigService) Del(ctx context.Context, key string) error {
	req, err := f.client.NewRequest(ctx, http.MethodDelete, "/v2/config/"+key, nil)
	if err != nil {
		return err
	}
	_, err = f.client.Do(req, nil)
	return err
}
