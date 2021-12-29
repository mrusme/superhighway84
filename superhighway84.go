package main

import (
	"context"
	"embed"
	"net/url"
	"os"
	"runtime"
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
  if runtime.GOOS == "windows" {
    zap.RegisterSink("winfile", func(u *url.URL) (zap.Sink, error) {
      return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    })
  }

  cfg := zap.NewProductionConfig()
  if runtime.GOOS == "windows" {
    cfg.OutputPaths = []string{
      "stdout",
      "winfile:///" + filename,
    }
  } else {
    cfg.OutputPaths = []string{
      filename,
    }
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

  var articles []*models.Article
  var articlesRoots []*models.Article

  TUI := tui.Init(&EMBEDFS, cfg, logger)
  TUI.ArticlesDatasource = &articles
  TUI.ArticlesRoots = &articlesRoots

  db, err := database.NewDatabase(ctx, cfg.ConnectionString, cfg.CachePath, logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()

  TUI.CallbackRefreshArticles = func() (error) {
    articles, articlesRoots, err = db.ListArticles()
    return err
  }
  TUI.CallbackSubmitArticle = func(article *models.Article) (error) {
    return db.SubmitArticle(article)
  }

  err = db.Connect(func(address string) {
    TUI.Views["mainscreen"].(*tui.Mainscreen).SetFooter(address)
    articles, articlesRoots, _ = db.ListArticles()

    time.Sleep(time.Second * 2)
    TUI.SetView("mainscreen", true)

    TUI.RefreshData()
    TUI.Refresh()
    TUI.App.Draw()
  })
  if err != nil {
    log.Panicln(err)
  }


  go func() {
    peers := 0
    for {
      bw := db.IPFSNode.Reporter.GetBandwidthTotals()
      connections, err := db.IPFSCoreAPI.Swarm().Peers(context.Background())
      if err == nil {
        peers = len(connections)
      }
      TUI.SetStats(int64(peers), int64(bw.RateIn), int64(bw.RateOut), bw.TotalIn , bw.TotalOut)
      time.Sleep(time.Second * 5)
    }
  }()

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

