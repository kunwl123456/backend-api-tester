package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kunwl123456/backend-api-tester/config"
)

// HandleConfig updates the runtime configuration (Base URL + Token).
func HandleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	cfg := config.RuntimeConfig{
		BaseURL:  r.FormValue("base_url"),
		Token:    r.FormValue("token"),
		UserName: r.FormValue("username"),
	}
	config.Set(cfg)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "配置已保存"})
}

// HandleGetConfig returns current config (hides token partially).
func HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	masked := cfg.Token
	if len(masked) > 8 {
		masked = masked[:4] + "****" + masked[len(masked)-4:]
	}
	jsonResponse(w, http.StatusOK, map[string]string{
		"base_url": cfg.BaseURL,
		"token":    masked,
		"username": cfg.UserName,
	})
}

// HandleLogin proxies a login request and auto-saves the returned token.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := map[string]string{
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
		"provider": "password",
	}
	data, status, err := doJSON(http.MethodPost, "/auth/login", payload)
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	// Auto-save token if login succeeded.
	if status == http.StatusOK {
		var resp struct {
			Session struct {
				AccessToken string `json:"access_token"`
			} `json:"session"`
		}
		if json.Unmarshal(data, &resp) == nil && resp.Session.AccessToken != "" {
			cfg := config.Get()
			cfg.Token = resp.Session.AccessToken
			config.Set(cfg)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}
