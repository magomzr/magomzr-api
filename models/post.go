package models

import (
	"strings"

	"github.com/google/uuid"
)

type Post struct {
	Card
	Author       string `json:"author"`
	Content      string `json:"content"`
	IsDraft      bool   `json:"isDraft"`
	Layout       string `json:"layout"`
	ModifiedDate string `json:"modifiedDate"`
	Previous     Info   `json:"previous"`
	Next         Info   `json:"next"`
}

type Card struct {
	ID         string   `json:"id"`
	CreateDate string   `json:"createDate"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Tags       []string `json:"tags"`
}

type Info struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

type Tags map[string]int

func (p *Post) GenerateId() {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	if len(id) > 24 {
		id = id[:24]
	}
	p.ID = id
}

func (p *Post) TagsToLower() {
	for i, tag := range p.Tags {
		p.Tags[i] = strings.ToLower(tag)
	}
}
