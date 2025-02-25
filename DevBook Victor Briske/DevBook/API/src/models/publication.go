package models

import (
	"errors"
	"strings"
	"time"
)

type Publication struct {
	ID         uint64    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorID   uint64    `json:"author_id"`
	AuthorNick string    `json:"author_nick"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      uint64    `json:"likes"`
}

func (pub *Publication) Prepare() error {
	if erro := pub.validate(); erro != nil {
		return erro
	}

	pub.formate()
	return nil
}

func (pub *Publication) validate() error {
	if pub.Title == "" {
		return errors.New("o título não pode estar em branco")
	}
	if pub.Content == "" {
		return errors.New("o conteúdo não pode estar em branco")
	}

	return nil
}

func (pub *Publication) formate() {
	pub.Title = strings.TrimSpace(pub.Title)
	pub.Content = strings.TrimSpace(pub.Content)
}
