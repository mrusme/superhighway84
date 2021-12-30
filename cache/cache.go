package cache

import (
	"encoding/json"

	"github.com/mrusme/superhighway84/models"
	"github.com/tidwall/buntdb"
)

type Cache struct {
  db      *buntdb.DB
}

func NewCache() (*Cache, error) {
  var err error

  cache := new(Cache)
  cache.db, err = buntdb.Open(":memory:")
  if err != nil {
    return nil, err
  }

  return cache, nil
}

func(cache *Cache) Close() {
  cache.db.Close()
}

func(cache *Cache) StoreArticle(article *models.Article) (error) {
  modelJson, jsonErr := json.Marshal(article)
  if jsonErr != nil {
    return jsonErr
  }

  err := cache.db.Update(func(tx *buntdb.Tx) error {
    _, _, err := tx.Set(article.ID, string(modelJson), nil)
    return err
  })
  return err
}

func(cache *Cache) LoadArticle(article *models.Article) (error) {
  err := cache.db.View(func(tx *buntdb.Tx) error {
    value, err := tx.Get(article.ID)
    if err != nil{
      return err
    }

    json.Unmarshal([]byte(value), article)
    return nil
  })
  return err
}

