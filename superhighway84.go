package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mrusme/superhighway84/database"
	"github.com/mrusme/superhighway84/models"
	"go.uber.org/zap"
)


func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  dbInit := false
  dbInitValue := os.Getenv("SUPERHIGHWAY84_DB_INIT")
  if dbInitValue == "1" {
    dbInit = true
  }

  dbURI := os.Getenv("SUPERHIGHWAY84_DB_URI")
  if dbInit == false && dbURI == "" {
    log.Panicln("SUPERHIGHWAY84_DB_URI missing!")
  }

  dbCache := os.Getenv("SUPERHIGHWAY84_DB_CACHE")
  if dbCache == "" {
    log.Panicln("SUPERHIGHWAY84_DB_CACHE missing!")
  }

  logger, err := zap.NewDevelopment()
  if err != nil {
    log.Panicln(err)
  }

  db, err := database.NewDatabase(ctx, dbURI, dbCache, dbInit, logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()
  db.Connect()

  var input string
  for {
    fmt.Scanln(&input)

    switch input {
    case "q":
      return
    case "g":
      fmt.Scanln(&input)
      article, err := db.GetArticleByID(input)
      if err != nil {
        log.Println(err)
      } else {
        log.Println(article)
      }
    case "p":
      article := models.NewArticle()
      article.From = "test@example.com"
      article.Newsgroup = "comp.test"
      article.Subject = "This is a test!"
      article.Body = "Hey there, this is a test!"

      err = db.SubmitArticle(article)
      if err != nil {
        log.Println(err)
      } else {
        log.Println(article)
      }
    case "l":
      articles, err := db.ListArticles()
      if err != nil {
        log.Println(err)
      } else {
        log.Println(articles)
      }
    }

  }
}

