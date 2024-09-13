package notes

type Note struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
}
