package models

import (
	"strings"

	"github.com/google/uuid"
)

type Post struct {
	Card
	Author       string `json:"author" dynamodbav:"author"`
	Content      string `json:"content" dynamodbav:"content"`
	IsDraft      bool   `json:"isDraft" dynamodbav:"isDraft"`
	Layout       string `json:"layout" dynamodbav:"layout"`
	ModifiedDate string `json:"modifiedDate" dynamodbav:"modifiedDate"`
	Previous     Info   `json:"previous" dynamodbav:"previous"`
	Next         Info   `json:"next" dynamodbav:"next"`
}

type Card struct {
	ID         string   `json:"id" dynamodbav:"id"`
	CreateDate string   `json:"createDate" dynamodbav:"createDate"`
	Title      string   `json:"title" dynamodbav:"title"`
	Summary    string   `json:"summary" dynamodbav:"summary"`
	Tags       []string `json:"tags" dynamodbav:"tags"`
}

type Info struct {
	Title string `json:"title" dynamodbav:"title"`
	ID    string `json:"id" dynamodbav:"id"`
}

type Tags map[string]int

func (p *Post) GenerateID() {
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
