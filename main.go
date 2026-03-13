package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/kunwl123456/backend-api-tester/handlers"
)

//go:embed templates/* static/*
var assetFS embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(assetFS, "templates/*.html"))
}

func renderPage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, name, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	staticFS, _ := fs.Sub(assetFS, "static")

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// ── Pages ──────────────────────────────────────────────────────────────
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderPage("index.html")(w, r)
	})
	mux.HandleFunc("/character", renderPage("character.html"))
	mux.HandleFunc("/quest", renderPage("quest.html"))
	mux.HandleFunc("/dialogue", renderPage("dialogue.html"))
	mux.HandleFunc("/announcement", renderPage("announcement.html"))
	mux.HandleFunc("/item", renderPage("item.html"))

	// ── Config API ─────────────────────────────────────────────────────────
	mux.HandleFunc("/api/config", handlers.HandleConfig)
	mux.HandleFunc("/api/config/get", handlers.HandleGetConfig)

	// ── Auth API ───────────────────────────────────────────────────────────
	mux.HandleFunc("/api/auth/login", handlers.HandleLogin)

	// ── Character API ──────────────────────────────────────────────────────
	mux.HandleFunc("/api/character/create", handlers.HandleCharacterCreate)
	mux.HandleFunc("/api/character/generate", handlers.HandleCharacterGenerate)
	mux.HandleFunc("/api/character/info", handlers.HandleCharacterInfo)
	mux.HandleFunc("/api/character/poll", handlers.HandleCharacterPoll)
	mux.HandleFunc("/api/character/list", handlers.HandleCharacterList)
	mux.HandleFunc("/api/character/check-name", handlers.HandleCharacterCheckName)

	// ── Quest API ──────────────────────────────────────────────────────────
	mux.HandleFunc("/api/quest/generate", handlers.HandleQuestGenerate)
	mux.HandleFunc("/api/quest/result", handlers.HandleQuestResult)
	mux.HandleFunc("/api/quest/stages", handlers.HandleQuestAllStages)
	mux.HandleFunc("/api/quest/poll", handlers.HandleQuestPoll)
	mux.HandleFunc("/api/quest/item", handlers.HandleQuestItem)

	// ── Item API ───────────────────────────────────────────────────────────
	mux.HandleFunc("/api/item/generate", handlers.HandleItemGenerate)
	mux.HandleFunc("/api/item/status", handlers.HandleItemStatus)

	// ── Dialogue API ───────────────────────────────────────────────────────
	mux.HandleFunc("/api/dialogue/talk", handlers.HandleDialogueTalk)
	mux.HandleFunc("/api/dialogue/history", handlers.HandleDialogueHistory)
	mux.HandleFunc("/api/dialogue/rewrite", handlers.HandleDialogueRewrite)
	mux.HandleFunc("/api/dialogue/delete-history", handlers.HandleDialogueDeleteHistory)

	// ── Announcement / Info API ────────────────────────────────────────────
	mux.HandleFunc("/api/announcement", handlers.HandleAnnouncement)
	mux.HandleFunc("/api/agreement", handlers.HandleAgreement)
	mux.HandleFunc("/api/user/info", handlers.HandleUserInfo)
	mux.HandleFunc("/api/activity/stats", handlers.HandleActivityStats)

	log.Printf("Backend API Tester 已启动 → http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
