package main

import (
	"code-hw/internal/yandexSpeller"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"

	"code-hw/internal/auth"
	"code-hw/internal/logging"
	"code-hw/internal/notes"
	"code-hw/internal/storage"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalln("Не задан параметр DATABASE_URL")
		return
	}

	db, err := connectToDB(connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	store := storage.NewSqlStorage(db)

	router := mux.NewRouter()

	router.Use(logging.LogRequest)
	router.Use(auth.BasicAuth(store))

	sp := yandexSpeller.NewSpeller()
	h := notes.NewHandler(store, sp)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Пример логирования

	router.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {

		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Добавление новой заметки")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		h.GetNotes(ctx, w, r)
	}).Methods("GET")

	router.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		h.AddNote(ctx, w, r)
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
