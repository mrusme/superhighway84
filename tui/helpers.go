package tui

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)

func OpenArticle(app *tview.Application, article *models.Article) (models.Article, error) {
  tmpFile, err := ioutil.TempFile(os.TempDir(), "article-*.txt")
  if err != nil {
    return *article, err
  }

  defer os.Remove(tmpFile.Name())

  tmpContent := []byte(fmt.Sprintf("Subject: %s\n= = = = = =\n%s", article.Subject, article.Body))
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

  content := strings.Split(string(tmpContent), "\n= = = = = =\n")
  if len(content) != 2 {
    return *article, errors.New("Document malformatted")
  }

  headerPart := strings.TrimSpace(content[0])
  subject := strings.TrimPrefix(headerPart, "Subject: ")
  if len(subject) <= 1 {
    return *article, errors.New("Invalid subject")
  }

  body := strings.TrimSpace(content[1])
  // TODO: Perform more validations
  if len(body) <= 1 {
    return *article, errors.New("Invalid body")
  }

  newArticle := *article
  newArticle.Subject = subject
  newArticle.Body = body

  return newArticle, nil
}

