package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"booking-backend/internal/config"
	"booking-backend/internal/server"
	"booking-backend/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	h := server.NewHandler(st, cfg)

	r := h.Routes()

	staticDir := "frontend/dist"
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs := http.FileServer(http.Dir(staticDir))
			path := filepath.Join(staticDir, r.URL.Path)
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				fs.ServeHTTP(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		}))
	}

	addr := ":" + cfg.Port
	log.Printf("Календарь звонков — бэкенд запущен на %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}
