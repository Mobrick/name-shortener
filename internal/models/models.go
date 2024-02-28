package models

type BatchRequestURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

type SimpleURLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
