package models

type Image struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
}
