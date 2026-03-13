package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kunwl123456/backend-api-tester/config"
)

var httpClient = &http.Client{Timeout: 120 * time.Second}

// doRequest sends an HTTP request to the backend and returns the response body.
func doRequest(method, path string, body io.Reader, contentType string) ([]byte, int, error) {
	cfg := config.Get()
	if cfg.BaseURL == "" {
		return nil, 0, fmt.Errorf("后端地址未配置，请先在配置页设置 Base URL")
	}
	url := cfg.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

// jsonResponse writes a JSON response to the client.
func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// errorJSON writes a JSON error message.
func errorJSON(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// doJSON sends a JSON POST request.
func doJSON(method, path string, payload any) ([]byte, int, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, 0, err
	}
	return doRequest(method, path, buf, "application/json")
}
