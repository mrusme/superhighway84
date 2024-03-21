package rss

import (
  "fmt"
  "io"
  "github.com/mrusme/superhighway84/models"
)

// ☠️  This is Proof of concept code!
// For example, it's currently composing the feed for each request, which doesn't
// make sense given how (not) frequently the articles update.

type FeedOptions struct {
  MaxArticles    int
}

type Feed struct {
  Options     FeedOptions
  Articles    []*models.Article
}

func NewFeedOptions() (FeedOptions) {
  return FeedOptions{10}
}

// Apparently there's no `min` for ints in golang
// https://stackoverflow.com/questions/27516387/what-is-the-correct-way-to-find-the-min-between-two-integers-in-go
func minInt(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func NewFeed(articles []*models.Article, feedOptions FeedOptions) (Feed) {

  numArticles := minInt(len(articles), feedOptions.MaxArticles)
  articlesSlice := articles[0:numArticles]

  return Feed{
    Options: feedOptions,
    Articles: articlesSlice,
  }
}

const header = `
<rss version="2.0">
  <channel>
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <language>%s</language>
    <pubDate>%s</pubDate>
    <lastBuildDate>%s</lastBuildDate>
    <docs>%s</docs>
    <generator>%s</generator>
`

const footer = `
  </channel>
</rss>
`

const item = `
    <item>
      <title>%s</title>
      <link>%s</link>
      <description>%s</description>
      <pubDate>%s</pubDate>
      <guid>%s</guid>
    </item>
`

func (feed *Feed) Write(w io.Writer) (error) {

  title := "Superhighway84"
  link := "https://xn--gckvb8fzb.com/superhighway84/"
  description := "USENET-INSPIRED DECENTRALIZED INTERNET DISCUSSION SYSTEM"
  language := "en-us"
  pubDate := "Tue, 10 Jun 2003 04:00:00 GMT"
  buildDate := "Tue, 10 Jun 2003 09:41:01 GMT"
  docsLink := "http://blogs.law.harvard.edu/tech/rss"
  generator := "Superhighway84 RSS Generator"
  fmt.Fprintf(w, header, title, link, description, language, pubDate, buildDate, docsLink, generator)

  for _, article := range feed.Articles {
    title := article.Subject
    link := fmt.Sprintf("superhighway84://%s", article.ID)
    description := fmt.Sprintf("Posted by %s in %s<br>%d replies<br>%s", article.From, article.Newsgroup, len(article.Replies), article.Body)
    pubDate := "Tue, 03 Jun 2003 09:39:21 GMT"
    guid := article.ID
    fmt.Fprintf(w, item, title, link, description, pubDate, guid)
  }

  fmt.Fprintf(w, footer)

  return nil
}
