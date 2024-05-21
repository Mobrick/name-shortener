package models

// BatchRequestURL модель использующаяся при парсинге запроса на множественное сокращение.
type BatchRequestURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponseURL модель использующаяся формировании ответа после множественного сокращения.
type BatchResponseURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Request модель использующаяся при парсинге запроса на сокращение из формата JSON.
type Request struct {
	URL string `json:"url"`
}

// Response модель использующаяся при формировании запроса после  сокращения в формате JSON.
type Response struct {
	Result string `json:"result"`
}

// URLRecord запись использующася для взаимодействия с хранилищем.
type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// SimpleURLRecord запись использующася для упрощенного предствления данных из хранилища.
type SimpleURLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ConfigFromFile модель использующаяся при парсинге конфигурационного файла при старте сервера
type ConfigFromFile struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}
