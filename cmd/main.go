package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"code-hw/internal/auth"
	"code-hw/internal/logging"
	"code-hw/internal/notes"
	"code-hw/internal/storage"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://user:password@db:5432/basa?sslmode=disable"

	db, err := connectToDB(connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	store := storage.NewStorage(db)

	// Выполнение миграций
	err = storage.RunMigrations(store.DB)
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL запроса: %v", err)
	}

	router := mux.NewRouter()

	router.Use(logging.LogRequest)
	router.Use(auth.BasicAuth(store))

	router.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		notes.GetNotesHandler(w, r, store, ctx)
	}).Methods("GET")

	router.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		notes.AddNoteHandler(w, r, store, ctx)
	}).Methods("POST")

	log.Println("Сервер запущен :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func connectToDB(connStr string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Ошибка подключения к БД фаза проверки: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if err = db.Ping(); err == nil {
			return db, nil
		}
		log.Printf("Сервер не отвечает на запросы: %v", err)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}
