package main

type Film struct {
	Title    string   `json:"title"`
	ImageURL string   `json:"poster_href"`
	Rating   float64  `json:"rating"`
	Genres   []string `json:"genres"`
}

type FilmsResponse struct {
	Page           uint64 `json:"current_page"`
	PageSize       uint64 `json:"page_size"`
	CollectionName string `json:"collection_name"`
	Total          uint64 `json:"total"`
	Films          []Film `json:"films"`
}
