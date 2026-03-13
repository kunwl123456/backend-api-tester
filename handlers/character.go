package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/kunwl123456/backend-api-tester/config"
)

// HandleCharacterCreate proxies POST /character/create_npc_charac (multipart).
func HandleCharacterCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	for _, field := range []string{"name", "gender", "backgroundstory", "information"} {
		if v := r.FormValue(field); v != "" {
			_ = mw.WriteField(field, v)
		}
	}
	mw.Close()

	data, status, err := doRequest(http.MethodPost, "/character/create_npc_charac", buf, mw.FormDataContentType())
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleCharacterGenerate proxies POST /character/gen_npc_charac.
func HandleCharacterGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	characterID := r.FormValue("character_id")
	if characterID == "" {
		errorJSON(w, http.StatusBadRequest, "character_id 不能为空")
		return
	}

	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	_ = mw.WriteField("character_id", characterID)
	mw.Close()

	data, status, err := doRequest(http.MethodPost, "/character/gen_npc_charac", buf, mw.FormDataContentType())
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleCharacterInfo proxies GET /character/get_character_info/{id}.
func HandleCharacterInfo(w http.ResponseWriter, r *http.Request) {
	characterID := r.URL.Query().Get("character_id")
	if characterID == "" {
		errorJSON(w, http.StatusBadRequest, "character_id 不能为空")
		return
	}
	data, status, err := doRequest(http.MethodGet, "/character/get_character_info/"+characterID, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleCharacterPoll polls character info until stage matches expected or timeout.
func HandleCharacterPoll(w http.ResponseWriter, r *http.Request) {
	characterID := r.URL.Query().Get("character_id")
	expectStage := r.URL.Query().Get("expect_stage")
	if characterID == "" {
		errorJSON(w, http.StatusBadRequest, "character_id 不能为空")
		return
	}

	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		data, status, err := doRequest(http.MethodGet, "/character/get_character_info/"+characterID, nil, "")
		if err != nil {
			errorJSON(w, http.StatusBadGateway, err.Error())
			return
		}
		if status != http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write(data)
			return
		}
		var info map[string]any
		if json.Unmarshal(data, &info) == nil {
			stage, _ := info["stage"].(string)
			if stage == expectStage || expectStage == "" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(data)
				return
			}
			// Still processing, wait and retry.
		}
		time.Sleep(3 * time.Second)
	}
	errorJSON(w, http.StatusGatewayTimeout, fmt.Sprintf("等待角色状态 %s 超时", expectStage))
}

// HandleCharacterList proxies GET /user/characters.
func HandleCharacterList(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	if cfg.BaseURL == "" {
		errorJSON(w, http.StatusBadRequest, "后端地址未配置")
		return
	}
	data, status, err := doRequest(http.MethodGet, "/user/characters?include_details=true", nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleCharacterCheckName proxies POST /character/check_name.
func HandleCharacterCheckName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	data, status, err := doJSON(http.MethodPost, "/character/check_name", map[string]string{"name": name})
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// forwardBody is a helper for streaming body forwarding (unused but kept for reference).
var _ = io.Discard
