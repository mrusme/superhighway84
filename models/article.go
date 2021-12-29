package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Article struct {
  ID           string    `mmapstructure:"id" json:"id" validate:"uuid_rfc4122"`
  InReplyToID  string    `mmapstructure:"in-reply-to-id" json:"in-reply-to-id" validate:"omitempty,uuid_rfc4122"`
  From         string    `mmapstructure:"from" json:"from" validate:"required,printascii"`
  Newsgroup    string    `mmapstructure:"newsgroup" json:"newsgroup" validate:"required,min=2,max=80,printascii,lowercase"`
  Subject      string    `mmapstructure:"subject" json:"subject" validate:"required,min=2,max=128,printascii"`
  Date         int64     `mmapstructure:"date" json:"date" validate:"required,number"`
  Organization string    `mmapstructure:"organization" json:"organization" validate:"printascii"`
  Body         string    `mmapstructure:"body" json:"body" validate:"required,min=3,max=524288"`

  Replies      []*Article `mmapstructure:"-" json:"-" validate:"-"`
}

func NewArticle() (*Article) {
  article := new(Article)

  id, _ := uuid.NewUUID()
  article.ID = id.String()

  article.Date = time.Now().UnixNano() / int64(time.Millisecond)

  return article
}

func (article *Article) IsValid() (bool, error) {
  validate := validator.New()
  errs := validate.Struct(article)
  if errs != nil {
    // validationErrors := errs.(validator.ValidationErrors)
    return false, errs
  }

  return true, nil
}

