package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
  ID           string    `mmapstructure:"id" json:"id"`
  InReplyToID  string    `mmapstructure:"in-reply-to-id" json:"in-reply-to-id"`
  From         string    `mmapstructure:"from" json:"from"`
  Newsgroup    string    `mmapstructure:"newsgroup" json:"newsgroup"`
  Subject      string    `mmapstructure:"subject" json:"subject"`
  Date         int64     `mmapstructure:"date" json:"date"`
  Organization string    `mmapstructure:"organization" json:"organization"`
  Body         string    `mmapstructure:"body" json:"body"`
}

func NewArticle() (*Article) {
  article := new(Article)

  id, _ := uuid.NewUUID()
  article.ID = id.String()

  article.Date = time.Now().UnixNano() / int64(time.Millisecond)

  return article
}

