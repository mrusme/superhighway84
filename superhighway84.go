package main

import (
	"context"
	"embed"
	"time"

	"log"

	"github.com/mrusme/superhighway84/config"
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
    filename,
  }
  return cfg.Build()
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  cfg, err := config.LoadConfig()
  if err != nil {
    log.Panicln(err)
  }
  if cfg.WasSetup() == false {
    cfg.Setup()
  }

  logger, err := NewLogger(cfg.Logfile)
  if err != nil {
    log.Panicln(err)
  }

  var articles []models.Article

  TUI := tui.Init(&EMBEDFS, cfg, logger)
  TUI.ArticlesDatasource = &articles

  db, err := database.NewDatabase(ctx, cfg.ConnectionString, cfg.CachePath, logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()

  TUI.CallbackRefreshArticles = func() (error) {
    articles, err = db.ListArticles()
    return err
  }
  TUI.CallbackSubmitArticle = func(article *models.Article) (error) {
    return db.SubmitArticle(article)
  }

  err = db.Connect(func(address string) {
    TUI.Views["mainscreen"].(*tui.Mainscreen).SetFooter(address)
    articles, _ = db.ListArticles()

    time.Sleep(time.Second * 2)
    TUI.SetView("mainscreen", true)

    TUI.RefreshData()
    TUI.Refresh()
    TUI.App.Draw()
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


  TUI.Launch()
}

