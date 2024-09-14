package storage

import (
	"code-hw/internal/models"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib" // драйвер для работы с PostgreSQL
	_ "github.com/lib/pq"
)

type SqlStorage struct {
	DB *sql.DB
}

// Note - структура заметки

// NewStorage создает новый экземпляр Storage с уже установленным подключением
func NewSqlStorage(db *sql.DB) Storage {
	return &SqlStorage{DB: db}
}

// GetUserByUsername получает пользователя из базы данных по имени
func (s *SqlStorage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
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
func (s *SqlStorage) AddNote(ctx context.Context, text string, userId int) (models.Note, error) {
	var note models.Note
	query := `INSERT INTO notes (text, user_id) VALUES ($1, $2) RETURNING id, text, user_id, created_at`

	err := s.DB.QueryRowContext(ctx, query, text, userId).Scan(&note.ID, &note.Text, &note.UserID, &note.CreatedAt)
	if err != nil {
		return models.Note{}, fmt.Errorf("error inserting note: %w", err)
	}

	return note, nil
}

// GetNotes - получение всех заметок из базы данных
func (s *SqlStorage) GetNotes(ctx context.Context, userId int) ([]models.Note, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, text, created_at FROM notes where user_id = $1`, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Text, &note.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetNoteByID - получение заметки по ID
func (s *SqlStorage) GetNoteByID(ctx context.Context, id int) (models.Note, error) {
	var note models.Note
	query := `SELECT id, text, created_at FROM notes WHERE id = $1`
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Text, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Note{}, fmt.Errorf("note not found")
		}
		return models.Note{}, fmt.Errorf("error retrieving note: %w", err)
	}

	return note, nil
}
