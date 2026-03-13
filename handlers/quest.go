package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// HandleQuestGenerate proxies POST /quest_v2/generate_outline_and_quest.
func HandleQuestGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	stageNum, _ := strconv.Atoi(r.FormValue("stage_number"))
	payload := map[string]any{
		"world_id":    r.FormValue("world_id"),
		"npc_id":      r.FormValue("character_id"),
		"index":       stageNum,
	}
	if v := r.FormValue("current_world_state"); v != "" {
		payload["current_world_state"] = v
	}
	data, status, err := doJSON(http.MethodPost, "/quest_v2/generate_outline_and_quest", payload)
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleQuestResult proxies GET /quest_v2/get_quest_result/{npc_quests_id}/{stage_number}.
func HandleQuestResult(w http.ResponseWriter, r *http.Request) {
	npcQuestsID := r.URL.Query().Get("npc_quests_id")
	stageNumber := r.URL.Query().Get("stage_number")
	if npcQuestsID == "" || stageNumber == "" {
		errorJSON(w, http.StatusBadRequest, "npc_quests_id 和 stage_number 不能为空")
		return
	}
	data, status, err := doRequest(http.MethodGet,
		"/quest_v2/get_quest_result/"+npcQuestsID+"/"+stageNumber, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleQuestAllStages proxies GET /quest_v2/get_all_quest_stages/{npc_quests_id}.
func HandleQuestAllStages(w http.ResponseWriter, r *http.Request) {
	npcQuestsID := r.URL.Query().Get("npc_quests_id")
	if npcQuestsID == "" {
		errorJSON(w, http.StatusBadRequest, "npc_quests_id 不能为空")
		return
	}
	data, status, err := doRequest(http.MethodGet,
		"/quest_v2/get_all_quest_stages/"+npcQuestsID, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleQuestPoll polls /quest_v2/get_quest_result until status=completed or timeout.
func HandleQuestPoll(w http.ResponseWriter, r *http.Request) {
	npcQuestsID := r.URL.Query().Get("npc_quests_id")
	stageNumber := r.URL.Query().Get("stage_number")
	if npcQuestsID == "" || stageNumber == "" {
		errorJSON(w, http.StatusBadRequest, "npc_quests_id 和 stage_number 不能为空")
		return
	}
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		data, status, err := doRequest(http.MethodGet,
			"/quest_v2/get_quest_result/"+npcQuestsID+"/"+stageNumber, nil, "")
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
		var resp map[string]any
		if json.Unmarshal(data, &resp) == nil {
			s, _ := resp["status"].(string)
			if s == "completed" || s == "failed" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(data)
				return
			}
		}
		time.Sleep(3 * time.Second)
	}
	errorJSON(w, http.StatusGatewayTimeout, "等待任务生成超时")
}

// HandleQuestItem proxies GET /quest_v2/get_generated_item/{item_id}.
func HandleQuestItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("item_id")
	if itemID == "" {
		errorJSON(w, http.StatusBadRequest, "item_id 不能为空")
		return
	}
	data, status, err := doRequest(http.MethodGet, "/quest_v2/get_generated_item/"+itemID, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleItemGenerate proxies POST /itemGen/generate.
func HandleItemGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := map[string]any{
		"item_prompt": r.FormValue("item_prompt"),
	}
	if v := r.FormValue("item_sort"); v != "" {
		payload["item_sort"] = v
	}
	if v := r.FormValue("context"); v != "" {
		payload["context"] = v
	}
	data, status, err := doJSON(http.MethodPost, "/itemGen/generate", payload)
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// HandleItemStatus proxies GET /itemGen/task_status/{task_id}.
func HandleItemStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		errorJSON(w, http.StatusBadRequest, "task_id 不能为空")
		return
	}
	data, status, err := doRequest(http.MethodGet, "/itemGen/task_status/"+taskID, nil, "")
	if err != nil {
		errorJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}
