package database

import (
	"context"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/documentstore"
	config "github.com/ipfs/go-ipfs-config"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/mrusme/superhighway84/models"
)

type Database struct {
  ctx                 context.Context
  URI                 string
  Cache               string

  Logger              *zap.Logger
  IPFSNode            icore.CoreAPI
  OrbitDB             orbitdb.OrbitDB
  Store               orbitdb.DocumentStore
}

func (db *Database)init() (error) {
  var err error

  db.OrbitDB, err = orbitdb.NewOrbitDB(db.ctx, db.IPFSNode, &orbitdb.NewOrbitDBOptions{
    Directory: &db.Cache,
    Logger: db.Logger,
  })
  if err != nil {
    return err
  }

  ac := &accesscontroller.CreateAccessControllerOptions{
    Access: map[string][]string{
      "write": {
        "*",
      },
    },
  }

  if err != nil {
    return err
  }

  addr, err := db.OrbitDB.DetermineAddress(db.ctx, "sync-test", "docstore", &orbitdb.DetermineAddressOptions{})
  if err != nil {
    return err
  }
  db.URI = addr.String()

  db.Store, err = db.OrbitDB.Docs(db.ctx, "sync-test", &orbitdb.CreateDBOptions{
    AccessController: ac,
    StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
  })
  if err != nil {
    return err
  }

  return nil
}

func (db *Database)open() (error) {
  var err error

  db.OrbitDB, err = orbitdb.NewOrbitDB(db.ctx, db.IPFSNode, &orbitdb.NewOrbitDBOptions{
    Directory: &db.Cache,
  })
  if err != nil {
    return err
  }

  create := false
  storetype := "docstore"
  dbstore, err := db.OrbitDB.Open(db.ctx, db.URI, &orbitdb.CreateDBOptions{
    Create: &create,
    StoreType: &storetype,
    StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
  })
  if err != nil {
    return err
  }

  db.Store = dbstore.(orbitdb.DocumentStore)

  return nil
}

func NewDatabase(
  ctx context.Context,
  dbURI string,
  dbCache string,
  dbInit bool,
  logger *zap.Logger,
) (*Database, error) {
  var err error

  db := new(Database)
  db.ctx = ctx
  db.URI = dbURI
  db.Cache = dbCache
  db.Logger = logger


  defaultPath, err := config.PathRoot()
  if err != nil {
    return nil, err
  }

  if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

  db.IPFSNode, err = createNode(ctx, defaultPath)
  if err != nil {
    return nil, err
  }


  if dbInit {
    err = db.init()
    if err != nil {
      return nil, err
    }
  } else {
    err = db.open()
    if err != nil {
      return nil, err
    }
  }


	// someDirectory, err := getUnixfsNode(dbCache)
	// if err != nil {
	// 	panic(fmt.Errorf("Could not get File: %s", err))
	// }
  // cidDirectory, err := ipfs.Unixfs().Add(ctx, someDirectory)
	// if err != nil {
	// 	panic(fmt.Errorf("Could not add Directory: %s", err))
	// }
  //
	// fmt.Printf("Added directory to IPFS with CID %s\n", cidDirectory.String())

  err = db.Store.Load(ctx, -1)
  if err != nil {
    // TODO: clean up
    return nil, err
  }

  // log.Println(db.Store.ReplicationStatus().GetBuffered())
  // log.Println(db.Store.ReplicationStatus().GetQueued())
  // log.Println(db.Store.ReplicationStatus().GetProgress())

  db.Logger.Info("running ...")

  return db, nil
}

func (db *Database) Connect() {
	go func() {
		err := connectToPeers(db.ctx, db.IPFSNode)
		if err != nil {
      db.Logger.Debug("failed to connect: %s", zap.Error(err))
    } else {
      db.Logger.Debug("connected to peer!")
    }
	}()
}

func (db *Database) Disconnect() {
  db.OrbitDB.Close()
}

func (db *Database) SubmitArticle(article *models.Article) (error) {
  entity := structToMap(&article)
  entity["type"] = "article"

  _, err := db.Store.Put(db.ctx, entity)
  return err
}

func (db *Database) GetArticleByID(id string) (models.Article, error) {
  entity, err := db.Store.Get(db.ctx, id, &iface.DocumentStoreGetOptions{CaseInsensitive: false})
  if err != nil {
    return models.Article{}, err
  }

  var article models.Article
  err = mapstructure.Decode(entity[0], &article)
  if err != nil {
    return models.Article{}, err
  }

  return article, nil
}

func (db *Database) ListArticles() ([]models.Article, error) {
  var articles []models.Article

  entities, err := db.Store.Query(db.ctx, func(e interface{})(bool, error) {
    entity := e.(map[string]interface{})
    if entity["type"] == "article" {
      return true, nil
    }
    return false, nil
  })
  if err != nil {
    return articles, err
  }

  for _, entity := range entities {
    var article models.Article
    err = mapstructure.Decode(entity, &article)
    if err != nil {
      return articles, err
    }
    articles = append(articles, article)
  }

  return articles, nil
}

