package main

import (
	"context"
	"embed"
	"fmt"
	"time"

	"log"
	"os"

	"github.com/mrusme/superhighway84/database"
	"github.com/mrusme/superhighway84/models"
	"github.com/mrusme/superhighway84/tui"
	"go.uber.org/zap"
)

//go:embed superhighway84.jpeg
var EMBEDFS embed.FS

func NewLogger(filename string) (*zap.Logger, error) {
  cfg := zap.NewProductionConfig()
  cfg.OutputPaths = []string{
    fmt.Sprintf("%s.log", filename),
  }
  return cfg.Build()
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  dbInit := false
  dbInitValue := os.Getenv("SUPERHIGHWAY84_DB_INIT")
  if dbInitValue == "1" {
    dbInit = true
  }

  dbName := os.Getenv("SUPERHIGHWAY84_DB_NAME")
  if dbInit == true && dbName == "" {
    log.Panicln("SUPERHIGHWAY84_DB_NAME missing!")
  }

  dbURI := os.Getenv("SUPERHIGHWAY84_DB_URI")
  if dbInit == false && dbURI == "" {
    log.Panicln("SUPERHIGHWAY84_DB_URI missing!")
  }

  dbCache := os.Getenv("SUPERHIGHWAY84_DB_CACHE")
  if dbCache == "" {
    log.Panicln("SUPERHIGHWAY84_DB_CACHE missing!")
  }

  logger, err := NewLogger(os.Getenv("LOGFILE"))
  if err != nil {
    log.Panicln(err)
  }

  var articles []models.Article

  TUI := tui.Init(&EMBEDFS)
  TUI.ArticlesDatasource = &articles

  db, err := database.NewDatabase(ctx, dbURI, dbCache, dbInit, dbName, logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()
  err = db.Connect(func(address string) {
    //TUI.App.Stop()
    TUI.Views["mainscreen"].(*tui.Mainscreen).SetFooter(address)
    articles, _ = db.ListArticles()
  })
  if err != nil {
    log.Panicln(err)
  }


  // ======================== TESTING ===============================
  // var articles []models.Article
  // mockGroups := []string{
  //   "comp.test",
  //   "news.conspiracy",
  //   "sci.physics",
  //   "talk.lolz",
  //   "sci.chemistry",
  //   "talk.random",
  //   "alt.anarchism",
  //   "alt.tv.simpsons",
  // }
  //
  // go func() {
  //   var prev models.Article
  //   for i := 0; i < 100; i++ {
  //     grp := mockGroups[(rand.Intn(len(mockGroups) - 1))]
  //
  //     time.Sleep(time.Millisecond * 250)
  //     art1 := *models.NewArticle()
  //     art1.Subject = fmt.Sprintf("A test in %s", grp)
  //     art1.Body = "This is just a test article\nWhat's up there?"
  //     art1.From = "test@example.com"
  //     art1.Newsgroup = grp
  //
  //     if prev.Newsgroup == art1.Newsgroup {
  //       art1.InReplyToID = prev.ID
  //       art1.Subject = fmt.Sprintf("Re: %s", prev.Subject)
  //     }
  //
  //     articles = append(articles, art1)
  //     prev = art1
  //   }
  // }()
  // ======================== /TESTING ==============================

  TUI.CallbackRefreshArticles = func() (error) {
    articles, err = db.ListArticles()
    return err
  }
  TUI.CallbackSubmitArticle = func(article *models.Article) (error) {
    return db.SubmitArticle(article)
    // return nil
  }

  go func() {
    time.Sleep(time.Second * 2)
    TUI.SetView("mainscreen", true)
    TUI.Refresh()
  }()
  TUI.Launch()

  // var input string
  // for {
  //   fmt.Scanln(&input)
  //
  //   switch input {
  //   case "q":
  //     return
  //   case "g":
  //     fmt.Scanln(&input)
  //     article, err := db.GetArticleByID(input)
  //     if err != nil {
  //       log.Println(err)
  //     } else {
  //       log.Println(article)
  //     }
  //   case "p":
  //     article := models.NewArticle()
  //     article.From = "test@example.com"
  //     article.Newsgroup = "comp.test"
  //     article.Subject = "This is a test!"
  //     article.Body = "Hey there, this is a test!"
  //
  //     err = db.SubmitArticle(article)
  //     if err != nil {
  //       log.Println(err)
  //     } else {
  //       log.Println(article)
  //     }
  //   case "l":
  //     articles, err := db.ListArticles()
  //     if err != nil {
  //       log.Println(err)
  //     } else {
  //       log.Println(articles)
  //     }
  //   }
  //
  // }
}

