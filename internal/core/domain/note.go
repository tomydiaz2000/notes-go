package domain

import (
	"fmt"
	"time"
)

type Note struct {
	ID           string
	Title        string
	Content      string
	ValidUntilAt time.Time
	CreateAt     time.Time
	UpdateAt     time.Time
	DeleteAt     time.Time
}

func NewNote(title string, content string, validUntilAt time.Time) *Note {
	now := time.Now()
	return &Note{
		Title:        title,
		Content:      content,
		ValidUntilAt: validUntilAt,
		CreateAt:     now,
		UpdateAt:     now,
	}
}

func (n *Note) Update(title string, content string, validUntilAt time.Time) error {
	if time.Now().After(n.ValidUntilAt) {
		return fmt.Errorf("cannot update expired note")
	}
	n.Title = title
	n.Content = content
	n.ValidUntilAt = validUntilAt
	n.UpdateAt = time.Now()
	return nil
}
