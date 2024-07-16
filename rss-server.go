package main

import (
  "context"
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
  "github.com/mrusme/superhighway84/rss"
  "go.uber.org/zap"
)

const LISTEN_ADDR_AND_PORT = ":8080"

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

  log.Println("loading configuration ...")
  cfg, err := config.LoadConfig()
  if err != nil {
    log.Panicln(err)
  }
  if cfg.WasSetup() == false {
    cfg.Setup()
  }

  log.Println("initializing logger ...")
  logger, err := NewLogger(cfg.Logfile)
  if err != nil {
    log.Panicln(err)
  }

  log.Println("initializing cache ...")
  cch, err := cache.NewCache(cfg.ProgramCachePath)
  if err != nil {
    log.Panicln(err)
  }
  defer cch.Close()

  var articles []*models.Article
  var articlesRoots []*models.Article

  log.Println(articles)
  log.Println(articlesRoots)

  log.Println("Creating DB")
  db, err := database.NewDatabase(ctx, cfg.ConnectionString, cfg.DatabaseCachePath, cch, logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()

  log.Println("Connecting")
  err = db.Connect(func(address string) {
    articles, articlesRoots, _ = db.ListArticles()
    log.Printf("Pre-loaded %d articles, %d roots", len(articles), len(articlesRoots))
  })
  if err != nil {
    log.Panicln(err)
  }

  log.Println("Connected")

  // ☠️  This is Proof of concept code.
  // It's ugly, it's buggy and it's very much not thread-safe!

  go func () {
    for {
      select {
      case <-ctx.Done():
        return
      case <-time.After(30 * time.Second):
        log.Println("Refreshing...")
        articles, articlesRoots, _ = db.ListArticles()
        log.Printf("Loaded %d articles, %d roots", len(articles), len(articlesRoots))
      }
    }
  }()

  createResponse := func(w http.ResponseWriter, articles []*models.Article) {

    feedOptions := rss.NewFeedOptions()
    rssFeed := rss.NewFeed(articles, feedOptions)
    err := rssFeed.Write(w)

    if err != nil {
      log.Printf("Failed to write feed in response: %v", err)
    }

  }

  articlesHandler := func(w http.ResponseWriter, r *http.Request) {
    log.Printf("/%s", r.URL.Path[1:])
    createResponse(w, articlesRoots)
  }

  commentsHandler := func(w http.ResponseWriter, r *http.Request) {
    log.Printf("/%s", r.URL.Path[1:])
    createResponse(w, articles)
  }

  http.HandleFunc("/rss/articles", articlesHandler)
  http.HandleFunc("/rss/comments", commentsHandler)

  log.Printf("Listening on %s", LISTEN_ADDR_AND_PORT)
  log.Fatal(http.ListenAndServe(LISTEN_ADDR_AND_PORT, nil))
}
