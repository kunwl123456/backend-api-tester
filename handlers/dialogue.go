package handlers

import (
	"net/http"

	"github.com/kunwl123456/backend-api-tester/config"
)

// HandleDialogueTalk proxies POST /gameagent/talk.
func HandleDialogueTalk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	cfg := config.Get()
	userName := r.FormValue("user_name")
	if userName == "" {
		userName = cfg.UserName
	}
	payload := map[string]any{
		"timestamp":       float64(0),
		"world_id":        r.FormValue("world_id"),
		"character_id":    r.FormValue("character_id"),
		"user_name":       userName,
		"text":            r.FormValue("text"),
		"observation_log": r.FormValue("observation_log"),
		"environment":     r.FormValue("environment"),
		"skills":          map[string]int{},
		"is_audio":        false,
		"user_pronoun":    r.FormValue("user_pronoun"),
	}
	data, status, err := doJSON(http.MethodPost, "/gameagent/talk", payload)
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleDialogueHistory proxies GET /internal/talk/history.
func HandleDialogueHistory(w http.ResponseWriter, r *http.Request) {
	worldID := r.URL.Query().Get("world_id")
	characterID := r.URL.Query().Get("character_id")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}
	path := "/internal/talk/history?world_id=" + worldID +
		"&character_id=" + characterID + "&limit=" + limit
	data, status, err := doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleDialogueRewrite proxies POST /gameagent/dialogue_rewrite.
func HandleDialogueRewrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	// Expect raw JSON body from client.
	data, status, err := doRequest(http.MethodPost, "/gameagent/dialogue_rewrite", r.Body, "application/json")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleDialogueDeleteHistory proxies DELETE /internal/talk/history.
func HandleDialogueDeleteHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := map[string]string{
		"world_id":     r.FormValue("world_id"),
		"character_id": r.FormValue("character_id"),
	}
	data, status, err := doJSON(http.MethodDelete, "/internal/talk/history", payload)
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}
