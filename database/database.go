package database

import (
	"context"
	"log"
	"sync"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/events"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/documentstore"
	config "github.com/ipfs/go-ipfs-config"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/mrusme/superhighway84/models"
)

type Database struct {
  ctx                 context.Context
  Name                string
  Init                bool
  URI                 string
  Cache               string

  Logger              *zap.Logger
  IPFSNode            icore.CoreAPI
  OrbitDB             orbitdb.OrbitDB
  Store               orbitdb.DocumentStore
  StoreEventChan      <-chan events.Event
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

  addr, err := db.OrbitDB.DetermineAddress(db.ctx, db.Name, "docstore", &orbitdb.DetermineAddressOptions{})
  if err != nil {
    return err
  }
  db.URI = addr.String()

  db.Store, err = db.OrbitDB.Docs(db.ctx, db.Name, &orbitdb.CreateDBOptions{
    AccessController: ac,
    StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
  })
  if err != nil {
    return err
  }

  db.StoreEventChan = db.Store.Subscribe(db.ctx)
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

  db.StoreEventChan = db.Store.Subscribe(db.ctx)
  return nil
}

func(db *Database) connectToPeers() error {
	var wg sync.WaitGroup

  peerInfos, err := config.DefaultBootstrapPeers()
  if err != nil {
    return err
  }

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := db.IPFSNode.Swarm().Connect(db.ctx, *peerInfo)
			if err != nil {
        db.Logger.Debug("failed to connect", zap.String("peerID", peerInfo.ID.String()), zap.Error(err))
			} else {
        db.Logger.Debug("connected!", zap.String("peerID", peerInfo.ID.String()))
      }
		}(&peerInfo)
	}
	wg.Wait()
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
  db.Name = "sync-test"
  db.Init = dbInit
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

  return db, nil
}

func (db *Database) Connect(onReady func()) (error) {
  var err error

  if db.Init {
    err = db.init()
    if err != nil {
      return err
    }
  } else {
    err = db.open()
    if err != nil {
      return err
    }
  }

  // log.Println(db.Store.ReplicationStatus().GetBuffered())
  // log.Println(db.Store.ReplicationStatus().GetQueued())
  // log.Println(db.Store.ReplicationStatus().GetProgress())

  db.Logger.Info("running ...")

	// go func() {
		err = db.connectToPeers()
		if err != nil {
      db.Logger.Debug("failed to connect: %s", zap.Error(err))
    } else {
      db.Logger.Debug("connected to peer!")
    }
	// }()

  go func() {
    for {
      for ev := range db.StoreEventChan {
        log.Printf("GOT EVENT %+v\n", ev)
        switch ev.(type) {
        case *stores.EventReady:
          onReady()
        }
      }
    }
  }()

  err = db.Store.Load(db.ctx, -1)
  if err != nil {
    // TODO: clean up
    return err
  }

  return nil
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

