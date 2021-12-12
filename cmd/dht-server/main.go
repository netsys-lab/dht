package main

import (
	"context"
	peer_store "github.com/netsys-lab/dht/peer-store"
	stdLog "log"
	"net"
	"os"
	"os/signal"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/log"
	"github.com/anacrolix/tagflag"

	"github.com/netsec-ethz/scion-apps/pkg/appnet"
	"github.com/netsys-lab/dht"
)

var (
	flags = struct {
		TableFile   string `help:"name of file for storing node info"`
		NoBootstrap bool
		NoSecurity  bool
	}{
		NoSecurity: true,
	}
	s *dht.Server
)

func loadTable() (err error) {
	added, err := s.AddNodesFromFile(flags.TableFile)
	log.Printf("loaded %d nodes from table file", added)
	return
}

func saveTable() error {
	return dht.WriteNodesToFile(s.Nodes(), flags.TableFile)
}

func main() {
	stdLog.SetFlags(stdLog.LstdFlags | stdLog.Lshortfile)
	err := mainErr()
	if err != nil {
		log.Printf("error in main: %v", err)
		os.Exit(1)
	}
}

func mainErr() error {
	tagflag.Parse(&flags)
	addr, _ := net.ResolveUDPAddr("udp", os.Getenv("DHT_LISTEN_ADDRESS"))
	conn, err := appnet.Listen(addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	cfg := dht.NewDefaultServerConfig()
	cfg.Conn = conn
	cfg.Logger = log.Default.FilterLevel(log.Info)
	cfg.NoSecurity = flags.NoSecurity
	cfg.PeerStore = &peer_store.InMemory{}
	s, err = dht.NewServer(cfg)
	if err != nil {
		return err
	}

	if flags.TableFile != "" {
		err = loadTable()
		if err != nil {
			return err
		}
	}
	log.Printf("dht server on %s, ID is %x", s.Addr(), s.ID())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		log.Printf("got signal: %v", <-ch)
		cancel()
	}()
	if !flags.NoBootstrap {
		go func() {
			if tried, err := s.Bootstrap(); err != nil {
				log.Printf("error bootstrapping: %s", err)
			} else {
				log.Printf("finished bootstrapping: %#v", tried)
			}
		}()
	}
	<-ctx.Done()
	s.Close()

	if flags.TableFile != "" {
		if err := saveTable(); err != nil {
			log.Printf("error saving node table: %s", err)
		}
	}
	return nil
}
