package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// データベース接続の初期化
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("データベース接続エラー:", err)
	}

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// 静的ファイルの提供設定
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	r.Handle("/*", http.FileServer(filesDir))

	// APIのエンドポイントの設定
	r.Route("/api", func(r chi.Router) {
		r.Post("/route", handleRoute)
	})

	port := ":8080"
	log.Printf("Server starting on %s...\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Server error:", err)
	}
}
