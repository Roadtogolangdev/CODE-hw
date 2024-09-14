package notes

import (
	"context"
	"net/http"
)

type Handler interface {
	AddNote(ctx context.Context, w http.ResponseWriter, r *http.Request)
	GetNotes(ctx context.Context, w http.ResponseWriter, r *http.Request)
}
