package database

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/documentstore"
	config "github.com/ipfs/go-ipfs-config"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/mrusme/superhighway84/cache"
	"github.com/mrusme/superhighway84/models"
)

type Database struct {
	ctx              context.Context
	ConnectionString string
	URI              string
	CachePath        string
	Cache            *cache.Cache

	Logger *zap.Logger

	IPFSNode    *core.IpfsNode
	IPFSCoreAPI icore.CoreAPI

	OrbitDB orbitdb.OrbitDB
	Store   orbitdb.DocumentStore
	Events  event.Subscription
}

func (db *Database) init() error {
	var err error

	ctx := context.Background()

	db.Logger.Debug("initializing NewOrbitDB ...")
	db.OrbitDB, err = orbitdb.NewOrbitDB(ctx, db.IPFSCoreAPI, &orbitdb.NewOrbitDBOptions{
		Directory: &db.CachePath,
		Logger:    db.Logger,
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

	// addr, err := db.OrbitDB.DetermineAddress(db.ctx, db.Name, "docstore", &orbitdb.DetermineAddressOptions{})
	// if err != nil {
	//   return err
	// }
	// db.URI = addr.String()

	storetype := "docstore"
	db.Logger.Debug("initializing OrbitDB.Docs ...")
	db.Store, err = db.OrbitDB.Docs(ctx, db.ConnectionString, &orbitdb.CreateDBOptions{
		AccessController:  ac,
		StoreType:         &storetype,
		StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
		Timeout:           time.Second * 600,
	})
	if err != nil {
		return err
	}

	db.Logger.Debug("subscribing to EventBus ...")
	db.Events, err = db.Store.EventBus().Subscribe(new(stores.EventReady))
	return nil
}

func (db *Database) GetOwnID() string {
	return db.OrbitDB.Identity().ID
}

func (db *Database) GetOwnPubKey() crypto.PubKey {
	pubKey, err := db.OrbitDB.Identity().GetPublicKey()
	if err != nil {
		return nil
	}

	return pubKey
}

func (db *Database) connectToPeers() error {
	var wg sync.WaitGroup

	peerInfos, err := config.DefaultBootstrapPeers()
	if err != nil {
		return err
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := db.IPFSCoreAPI.Swarm().Connect(db.ctx, *peerInfo)
			if err != nil {
				db.Logger.Error("failed to connect", zap.String("peerID", peerInfo.ID.String()), zap.Error(err))
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
	dbConnectionString string,
	dbCache string,
	cch *cache.Cache,
	logger *zap.Logger,
) (*Database, error) {
	var err error

	db := new(Database)
	db.ctx = ctx
	db.ConnectionString = dbConnectionString
	db.CachePath = dbCache
	db.Cache = cch
	db.Logger = logger

	db.Logger.Debug("getting config root path ...")
	defaultPath, err := config.PathRoot()
	if err != nil {
		return nil, err
	}

	db.Logger.Debug("setting up plugins ...")
	if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

	db.Logger.Debug("creating IPFS node ...")
	db.IPFSNode, db.IPFSCoreAPI, err = createNode(ctx, defaultPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) Connect(onReady func(address string)) error {
	var err error

	db.Logger.Info("connecting to peers ...")
	// go func() {
	err = db.connectToPeers()
	if err != nil {
		db.Logger.Error("failed to connect: %s", zap.Error(err))
	} else {
		db.Logger.Debug("connected to peer!")
	}
	// }()

	db.Logger.Info("initializing database connection ...")
	// if db.Init {
	err = db.init()
	if err != nil {
		db.Logger.Error("%s", zap.Error(err))
		return err
	}
	// } else {
	//   err = db.open()
	//   if err != nil {
	//     return err
	//   }
	// }

	db.Logger.Info("running ...")

	go func() {
		for {
			for ev := range db.Events.Out() {
				db.Logger.Debug("got event", zap.Any("event", ev))
				switch ev.(type) {
				case stores.EventReady:
					db.URI = db.Store.Address().String()
					onReady(db.URI)
					continue
				}
			}
		}
	}()

	err = db.Store.Load(db.ctx, -1)
	if err != nil {
		db.Logger.Error("%s", zap.Error(err))
		// TODO: clean up
		return err
	}

	db.Logger.Debug("connect done")
	return nil
}

func (db *Database) Disconnect() {
	err := db.Events.Close()
	if err != nil {
		log.Println(err)
	}
	err = db.Store.Close()
	if err != nil {
		log.Println(err)
	}
	err = db.OrbitDB.Close()
	if err != nil {
		log.Println(err)
	}
}

func (db *Database) SubmitArticle(article *models.Article) error {
	entity, err := structToMap(*article)
	if err != nil {
		return err
	}
	entity["type"] = "article"

	_, err = db.Store.Put(db.ctx, entity)
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

func (db *Database) ListArticles() ([]*models.Article, []*models.Article, error) {
	var articles []*models.Article
	var articlesMap map[string]*models.Article

	articlesMap = make(map[string]*models.Article)

	_, err := db.Store.Query(db.ctx, func(e interface{}) (bool, error) {
		entity := e.(map[string]interface{})
		if entity["type"] == "article" {
			var article models.Article
			err := mapstructure.Decode(entity, &article)
			if err == nil {
				// TODO: Not sure why mapstructure won't convert this field and simply
				//       leave it ""
				if entity["in-reply-to-id"] != nil {
					article.InReplyToID = entity["in-reply-to-id"].(string)
				}
				err := db.Cache.LoadArticle(&article)
				if err != nil {
					return false, err
				}
				articles = append(articles, &article)
				articlesMap[article.ID] = articles[(len(articles) - 1)]
			}
			return true, err
		}
		return false, nil
	})
	if err != nil {
		return articles, nil, err
	}

	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})

	var articlesRoots []*models.Article
	for i := 0; i < len(articles); i++ {
		if articles[i].InReplyToID != "" {
			inReplyTo := articles[i].InReplyToID
			if _, exist := articlesMap[inReplyTo]; exist == true {

				(*articlesMap[inReplyTo]).Replies =
					append((*articlesMap[inReplyTo]).Replies, articles[i])
				(*articlesMap[inReplyTo]).LatestReply = articles[i].Date
				continue
			}
		}
		articlesRoots = append(articlesRoots, articles[i])
	}

	sort.SliceStable(articlesRoots, func(i, j int) bool {
		iLatest := articlesRoots[i].LatestReply
		if iLatest <= 0 {
			iLatest = articlesRoots[i].Date
		}

		jLatest := articlesRoots[j].LatestReply
		if jLatest <= 0 {
			jLatest = articlesRoots[j].Date
		}

		return iLatest > jLatest
	})

	return articles, articlesRoots, nil
}
