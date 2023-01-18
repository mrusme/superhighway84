package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Article struct {
	ID           string `mapstructure:"id" json:"-" validate:"uuid_rfc4122"`
	InReplyToID  string `mapstructure:"in-reply-to-id" json:"-" validate:"omitempty,uuid_rfc4122"`
	From         string `mapstructure:"from" json:"-" validate:"required,printascii"`
	Newsgroup    string `mapstructure:"newsgroup" json:"-" validate:"required,min=2,max=80,printascii,lowercase"`
	Subject      string `mapstructure:"subject" json:"-" validate:"required,min=2,max=128,printascii"`
	Date         int64  `mapstructure:"date" json:"-" validate:"required,number"`
	Organization string `mapstructure:"organization" json:"-" validate:"printascii"`
	Body         string `mapstructure:"body" json:"-" validate:"required,min=3,max=524288"`

	Replies     []*Article `mapstructure:"-" json:"-" validate:"-"`
	LatestReply int64      `mapstructure:"-" json:"-" validate:"-"`

	Read bool `mapstructure:"-" json:"read" validate:"-"`
}

func NewArticle() *Article {
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
