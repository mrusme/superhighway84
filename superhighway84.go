package main

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"log"

	"github.com/mrusme/superhighway84/cache"
	"github.com/mrusme/superhighway84/config"
	"github.com/mrusme/superhighway84/database"
	"github.com/mrusme/superhighway84/models"
	"github.com/mrusme/superhighway84/tui"
	"go.uber.org/zap"
)

//go:embed superhighway84.jpeg
var EMBEDFS embed.FS
var version = "v0.0.0"

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

  cch, err := cache.NewCache()
  if err != nil {
    log.Panicln(err)
  }
  defer cch.Close()

  var articles []*models.Article
  var articlesRoots []*models.Article

  TUI := tui.Init(&EMBEDFS, cfg, cch, logger)
  TUI.SetVersion(version, getLatestVersion())

  TUI.ArticlesDatasource = &articles
  TUI.ArticlesRoots = &articlesRoots


  db, err := database.NewDatabase(ctx, cfg.ConnectionString, cfg.CachePath, cch, logger)
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

  TUI.Launch()
}


func getLatestVersion() (string) {
  var client = &http.Client{Timeout: 10 * time.Second}
  r, err := client.Get(
    "https://api.github.com/repos/mrusme/superhighway84/releases/latest",
  )
  if err != nil {
      return version
  }
  defer r.Body.Close()
  var result map[string]interface{}
  json.NewDecoder(r.Body).Decode(&result)

  if val, exist := result["tag_name"]; exist == true {
    return val.(string)
  }
  return version
}

