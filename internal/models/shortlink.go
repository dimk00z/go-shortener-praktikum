package models

type BatchURLs []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"-"`
}

type BatchShortURLs []struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURLs []struct {
	ShortURL string
	URL      string
}
