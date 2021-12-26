package tui

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)

func OpenArticle(app *tview.Application, article *models.Article) (models.Article, error) {
  tmpFile, err := ioutil.TempFile(os.TempDir(), "article-*.txt")
  if err != nil {
    return *article, err
  }

  defer os.Remove(tmpFile.Name())

  tmpContent := []byte(article.Body)
  if _, err = tmpFile.Write(tmpContent); err != nil {
    return *article, err
  }

  if err := tmpFile.Close(); err != nil {
    return *article, err
  }

  wasSuspended := app.Suspend(func() {
    cmd := exec.Command(os.Getenv("EDITOR"), tmpFile.Name())
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err := cmd.Run()
    if err != nil {
      log.Println(err)
    }
    return
  })

  if wasSuspended == false {
    return *article, err
  }

  tmpContent, err = os.ReadFile(tmpFile.Name())
  if err != nil {
    return *article, err
  }

  newArticle := *article
  newArticle.Body = string(tmpContent)

  return newArticle, nil
}

