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

type handler struct {
	storage storage.Storage
	speller yandexSpeller.Speller
}

func NewHandler(storage storage.Storage, speller yandexSpeller.Speller) Handler {
	return &handler{storage: storage, speller: speller}
}
func (h *handler) AddNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req AddNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	correctedResults, err := h.speller.CheckSpelling(req.Text)
	if err != nil {
		http.Error(w, "Error checking spelling: "+err.Error(), http.StatusInternalServerError)
		return
	}

	correctedText := req.Text
	if len(correctedResults) > 0 {
		correctedText = correctedResults[0].Suggestions[0]
	}

	note, err := h.storage.AddNote(ctx, correctedText, r.Context().Value("userID").(int))
	if err != nil {
		http.Error(w, "Error adding note: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}
func (h *handler) GetNotes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	notes, err := h.storage.GetNotes(ctx, ctx.Value("userID").(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
