package models

type Post struct {
	ID           string   `json:"id"`
	Author       string   `json:"author"`
	Content      string   `json:"content"`
	CreateDate   string   `json:"createDate"`
	IsDraft      bool     `json:"isDraft"`
	Layout       string   `json:"layout"`
	ModifiedDate string   `json:"modifiedDate"`
	Summary      string   `json:"summary"`
	Tags         []string `json:"tags"`
	Title        string   `json:"title"`
	Previous     Info     `json:"previous"`
	Next         Info     `json:"next"`
}

type Info struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}
