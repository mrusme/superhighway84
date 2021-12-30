package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Article struct {
  ID           string    `mmapstructure:"id" json:"-" validate:"uuid_rfc4122"`
  InReplyToID  string    `mmapstructure:"in-reply-to-id" json:"-" validate:"omitempty,uuid_rfc4122"`
  From         string    `mmapstructure:"from" json:"-" validate:"required,printascii"`
  Newsgroup    string    `mmapstructure:"newsgroup" json:"-" validate:"required,min=2,max=80,printascii,lowercase"`
  Subject      string    `mmapstructure:"subject" json:"-" validate:"required,min=2,max=128,printascii"`
  Date         int64     `mmapstructure:"date" json:"-" validate:"required,number"`
  Organization string    `mmapstructure:"organization" json:"-" validate:"printascii"`
  Body         string    `mmapstructure:"body" json:"-" validate:"required,min=3,max=524288"`

  Replies      []*Article `mmapstructure:"-" json:"-" validate:"-"`
  Read         bool      `mmapstructure:"-" json:"read" validate:"-"`
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

