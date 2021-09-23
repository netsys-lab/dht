// Pings DHT nodes with the given network addresses.
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/anacrolix/tagflag"

	"github.com/netsys-lab/dht/v2"
	"github.com/scionproto/scion/go/lib/snet"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var args = struct {
		Timeout time.Duration
		tagflag.StartPos
		Nodes []string `help:"nodes to ping e.g. router.bittorrent.com:6881"`
	}{}
	tagflag.Parse(&args)
	s, err := dht.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dht server on %s with id %x", s.Addr(), s.ID())
	var wg sync.WaitGroup
	for _, a := range args.Nodes {
		func(a string) {
			ua, err := snet.ParseUDPAddr(a)
			if err != nil {
				log.Fatal(err)
			}
			started := time.Now()
			wg.Add(1)
			go func() {
				defer wg.Done()
				res := s.Ping(ua)
				err := res.Err
				if err != nil {
					fmt.Printf("%s: %s: %s\n", a, time.Since(started), err)
					return
				}
				id := *res.Reply.SenderID()
				fmt.Printf("%s: %x %c: %s\n", a, id, func() rune {
					if dht.NodeIdSecure(id, ua.Host.IP) {
						return '✔'
					} else {
						return '✘'
					}
				}(), time.Since(started))
			}()
		}(a)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	timeout := make(chan struct{})
	if args.Timeout != 0 {
		go func() {
			time.Sleep(args.Timeout)
			close(timeout)
		}()
	}
	select {
	case <-done:
	case <-timeout:
		log.Print("timed out")
	}
}
