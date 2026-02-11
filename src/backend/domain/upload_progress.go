package domain

type UploadPhase string

const (
	PhaseReceiving UploadPhase = "receiving"
	PhaseUploading UploadPhase = "uploading"
	PhaseComplete  UploadPhase = "complete"
	PhaseError     UploadPhase = "error"
)

type UploadProgress struct {
	UploadID       string      `json:"upload_id"`
	Phase          UploadPhase `json:"phase"`
	BytesProcessed int64       `json:"bytes_processed"`
	TotalBytes     int64       `json:"total_bytes"`
}

type UploadComplete struct {
	UploadID string `json:"upload_id"`
	Result   any    `json:"result,omitempty"`
}

type UploadError struct {
	UploadID string `json:"upload_id"`
	Error    string `json:"error"`
}

type UploadEvent struct {
	EventType string
	Data      any
}
