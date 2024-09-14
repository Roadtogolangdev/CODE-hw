package storage

import (
	"code-hw/internal/models"
	"context"
)

type Storage interface {
	GetUserByUsername(username string) (*models.User, error)
	AddNote(ctx context.Context, text string, userId int) (models.Note, error)
	GetNotes(ctx context.Context, userId int) ([]models.Note, error)
	GetNoteByID(ctx context.Context, id int) (models.Note, error)
}
