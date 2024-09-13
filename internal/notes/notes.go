package notes

import (
	"code-hw/internal/storage"
	"code-hw/internal/yandexSpeller"
	"context"
	"encoding/json"
	"net/http"
)

type AddNoteRequest struct {
	Text string `json:"text"`
}

func AddNoteHandler(w http.ResponseWriter, r *http.Request, store *storage.Storage, ctx context.Context) {
	var req AddNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	correctedResults, err := yandexSpeller.CheckSpelling(req.Text)
	if err != nil {
		http.Error(w, "Error checking spelling", http.StatusInternalServerError)
		return
	}

	correctedText := req.Text
	if len(correctedResults) > 0 {
		correctedText = correctedResults[0].Suggestions[0] // Берём первое исправление
	}

	note, err := store.AddNote(ctx, correctedText)
	if err != nil {
		http.Error(w, "Error adding note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func GetNotesHandler(w http.ResponseWriter, r *http.Request, db *storage.Storage, ctx context.Context) {
	notes, err := db.GetNotes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
