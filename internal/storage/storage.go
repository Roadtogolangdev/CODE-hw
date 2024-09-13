package storage

import (
	"code-hw/internal/auth"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib" // драйвер для работы с PostgreSQL
	_ "github.com/lib/pq"
	"time"
)

type Storage struct {
	DB *sql.DB
}

// Note - структура заметки
type Note struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// NewStorage создает новый экземпляр Storage с уже установленным подключением
func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

// GetUserByUsername получает пользователя из базы данных по имени
func (s *Storage) GetUserByUsername(username string) (*auth.User, error) {
	var user auth.User
	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := s.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}

	return &user, nil
}

// AddNote - добавление новой заметки в базу данных
func (s *Storage) AddNote(ctx context.Context, text string) (Note, error) {
	var note Note
	query := `INSERT INTO notes (text) VALUES ($1) RETURNING id, text, created_at`

	err := s.DB.QueryRowContext(ctx, query, text).Scan(&note.ID, &note.Text, &note.CreatedAt)
	if err != nil {
		return Note{}, fmt.Errorf("error inserting note: %w", err)
	}

	return note, nil
}

// GetNotes - получение всех заметок из базы данных
func (s *Storage) GetNotes(ctx context.Context) ([]Note, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, text, created_at FROM notes`)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Text, &note.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetNoteByID - получение заметки по ID
func (s *Storage) GetNoteByID(ctx context.Context, id int) (Note, error) {
	var note Note
	query := `SELECT id, text, created_at FROM notes WHERE id = $1`
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Text, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, fmt.Errorf("note not found")
		}
		return Note{}, fmt.Errorf("error retrieving note: %w", err)
	}

	return note, nil
}
