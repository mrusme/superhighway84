package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/iface"
	"github.com/google/uuid"
	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/zap"
)

func setupPlugins(path string) error {
	// Load plugins. This will skip the repo if not available.
	plugins, err := loader.NewPluginLoader(filepath.Join(path, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo: repo,
    ExtraOpts: map[string]bool{
      "pubsub": true,
    },
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	// peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	// for _, addrStr := range peers {
	// 	addr, err := ma.NewMultiaddr(addrStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	pii, err := peer.AddrInfoFromP2pAddr(addr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	pi, ok := peerInfos[pii.ID]
	// 	if !ok {
	// 		pi = &peer.AddrInfo{ID: pii.ID}
	// 		peerInfos[pi.ID] = pi
	// 	}
	// 	pi.Addrs = append(pi.Addrs, pii.Addrs...)
	// }

  peerInfos, err := config.DefaultBootstrapPeers()
  if err != nil {
    return err
  }

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *peerInfo)
			if err != nil {
				log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
			} else {
        log.Printf("connected to %s!", peerInfo.ID)
      }
		}(&peerInfo)
	}
	wg.Wait()
	return nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func main() {
  var err error
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  log.Println("press w for writer or r for reader")
  var testtype string
  var testid string
  fmt.Scanln(&testtype)

  if testtype == "r" {
    log.Println("enter the id")
    fmt.Scanln(&testid)
  }

  defaultPath, err := config.PathRoot()
  if err != nil {
    log.Println(err)
    return
  }
  log.Println(defaultPath)

  if err := setupPlugins(defaultPath); err != nil {
    log.Println(err)
		return
	}

  ipfs, err := createNode(ctx, defaultPath)
  if err != nil {
    log.Println(err)
    return
  }

	bootstrapNodes := []string{
		// // IPFS Bootstrapper nodes.
		// "/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		// "/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		// "/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		// "/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
    //
		// // IPFS Cluster Pinning nodes
		// "/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		// "/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		// "/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		// "/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		// "/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		// "/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		// "/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
		// "/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
    //
		// // You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
		// // "/ip4/127.0.0.1/tcp/4010/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
		// // "/ip4/127.0.0.1/udp/4010/quic/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
	}

	go func() {
		err := connectToPeers(ctx, ipfs, bootstrapNodes)
		if err != nil {
			log.Printf("failed connect to peers: %s", err)
    } else {
      log.Println("connected to node!")
    }
	}()


  logger, err := zap.NewDevelopment()

  orbitDir := os.Getenv("SUPERHIGHWAY84_DB")
  var orbitdb1 orbitdb.OrbitDB
  var db1 orbitdb.DocumentStore
  if testtype == "w" {
    log.Println("Opening OrbitDB as writer ...")
    orbitdb1, err = orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
      Directory: &orbitDir,
      Logger: logger,
    })
    ac := &accesscontroller.CreateAccessControllerOptions{
      Access: map[string][]string{
        "write": {
          "*",
        },
      },
    }

    if err != nil {
      log.Println(err)
      return
    }
    defer orbitdb1.Close()

    log.Println(orbitdb1.Identity().ID)
    addr, err := orbitdb1.DetermineAddress(ctx, "sync-test", "docstore", &orbitdb.DetermineAddressOptions{})
    if err != nil {
      log.Println(err)
      return
    }
    log.Println(addr.String())

    db1, err = orbitdb1.Docs(ctx, "sync-test", &orbitdb.CreateDBOptions{
      AccessController: ac,
    })
    if err != nil {
      log.Println(err)
      return
    }
  } else {
    log.Println("Opening OrbitDB as reader ...")
    orbitdb1, err = orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
      Directory: &orbitDir,
    })
    if err != nil {
      log.Println(err)
      return
    }
    log.Println("NewOrbitDB succeeded")
    create := false
    storetype := "docstore"
    dbstore, err := orbitdb1.Open(ctx, testid, &orbitdb.CreateDBOptions{Create: &create, StoreType: &storetype})
    if err != nil {
      log.Println(err)
      return
    }
    log.Println("Test")
    db1 = dbstore.(orbitdb.DocumentStore)
  }
  log.Println("opened!")


	// someDirectory, err := getUnixfsNode(orbitDir)
	// if err != nil {
	// 	panic(fmt.Errorf("Could not get File: %s", err))
	// }
  // cidDirectory, err := ipfs.Unixfs().Add(ctx, someDirectory)
	// if err != nil {
	// 	panic(fmt.Errorf("Could not add Directory: %s", err))
	// }
  //
	// fmt.Printf("Added directory to IPFS with CID %s\n", cidDirectory.String())

  if testtype == "w" {
  } else {
    err = db1.Load(ctx, -1)
    if err != nil {
      log.Println(err)
      return
    }
  }

  log.Println(db1.ReplicationStatus().GetBuffered())
  log.Println(db1.ReplicationStatus().GetQueued())
  log.Println(db1.ReplicationStatus().GetProgress())

  log.Println("Running ...")

  var input string
  for {
    fmt.Scanln(&input)

    switch input {
    case "q":
      return
    case "g":
      fmt.Scanln(&input)
      docs, err := db1.Get(ctx, input, &iface.DocumentStoreGetOptions{CaseInsensitive: false})
      if err != nil {
        log.Println(err)
      } else {
        log.Println(docs)
      }
    case "p":
      id, _ := uuid.NewUUID()
      _, err = db1.Put(ctx, map[string]interface{}{"_id": id.String(), "hello": "world"})
      if err != nil {
        log.Println(err)
      } else {
        log.Println(id)
      }
    case "l":
      docs, err := db1.Query(ctx, func(e interface{})(bool, error) {
        return true, nil
      })
      if err != nil {
        log.Println(err)
      } else {
        log.Println(docs)
      }
    }

  }

}
