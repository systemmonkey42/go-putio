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

func (f *ConfigService) GetAll(ctx context.Context, config interface{}) error {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config", nil)
	if err != nil {
		return err
	}
	var r struct {
		Config json.RawMessage `json:"config"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return err
	}
	return json.Unmarshal(r.Config, &config)
}

func (f *ConfigService) Get(ctx context.Context, key string, value interface{}) error {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config/"+key, nil)
	if err != nil {
		return err
	}
	var r struct {
		Value json.RawMessage `json:"value"`
	}
	_, err = f.client.Do(req, &r)
	if err != nil {
		return err
	}
	return json.Unmarshal(r.Value, &value)
}

func (f *ConfigService) SetAll(ctx context.Context, config interface{}) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	v := struct {
		Config json.RawMessage `json:"config"`
	}{
		Config: b,
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
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	v := struct {
		Value json.RawMessage `json:"value"`
	}{
		Value: b,
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
