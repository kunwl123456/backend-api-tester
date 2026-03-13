package handlers

import (
	"net/http"
)

// HandleAnnouncement proxies GET /info/announcement.
func HandleAnnouncement(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("language")
	path := "/info/announcement"
	if lang != "" {
		path += "?language=" + lang
	}
	data, status, err := doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleAgreement proxies GET /info/agreement.
func HandleAgreement(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("language")
	path := "/info/agreement"
	if lang != "" {
		path += "?language=" + lang
	}
	data, status, err := doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleUserInfo proxies GET /user/get (current user).
func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	data, status, err := doRequest(http.MethodGet, "/user/get", nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleActivityStats proxies GET /activity/my/stats.
func HandleActivityStats(w http.ResponseWriter, r *http.Request) {
	days := r.URL.Query().Get("days")
	if days == "" {
		days = "30"
	}
	data, status, err := doRequest(http.MethodGet, "/activity/my/stats?days="+days, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}
