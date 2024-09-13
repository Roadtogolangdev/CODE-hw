package storage

import (
	"database/sql"
	"log"
)

func RunMigrations(DB *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);`

	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы пользователей: %v", err)
		return err
	}

	_, err = DB.Exec(notesTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы заметок: %v", err)
		return err
	}

	log.Println("Успешное завершение SQL-запроса")
	return nil
}
