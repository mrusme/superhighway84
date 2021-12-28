package main

import (
	"context"
	"embed"
	"time"

	"log"

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

  logger, err := zap.NewProduction()
  if err != nil {
    log.Panicln(err)
  }

  var articles []models.Article

  TUI := tui.Init(&EMBEDFS,  logger)
  TUI.ArticlesDatasource = &articles

  log.Println("Starting a new database")
  db, err := database.NewDatabase(ctx, "/orbitdb/bafyreifdpagppa7ve45odxuvudz5snbzcybwyfer777huckl4li4zbc5k4/superhighway84", "D://b", logger)
  if err != nil {
    log.Panicln(err)
  }
  defer db.Disconnect()

  log.Println("Database created")

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

