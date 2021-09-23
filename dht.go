package dht

import (
	"crypto"
	_ "crypto/sha1"
	"math/rand"
	"net"
	"time"

	peer_store "github.com/netsys-lab/dht/v2/peer-store"

	"github.com/anacrolix/log"
	"github.com/anacrolix/missinggo/v2"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/metainfo"

	"github.com/netsys-lab/dht/v2/krpc"
	"github.com/scionproto/scion/go/lib/snet"
)

func defaultQueryResendDelay() time.Duration {
	return jitterDuration(5*time.Second, time.Second)
}

// Uniquely identifies a transaction to us.
type transactionKey struct {
	RemoteAddr string // host:port
	T          string // The KRPC transaction ID.
}

type StartingNodesGetter func() ([]Addr, error)

// ServerConfig allows to set up a  configuration of the `Server` instance
// to be created with NewServer
type ServerConfig struct {
	// Set NodeId Manually. Caller must ensure that if NodeId does not conform
	// to DHT Security Extensions, that NoSecurity is also set.
	NodeId krpc.ID
	Conn   *snet.Conn
	// Don't respond to queries from other nodes.
	Passive       bool
	StartingNodes StartingNodesGetter
	// Disable the DHT security extension: http://www.libtorrent.org/dht_sec.html.
	NoSecurity bool
	// Initial IP blocklist to use. Applied before serving and bootstrapping
	// begins.
	IPBlocklist iplist.Ranger
	// Used to secure the server's ID. Defaults to the Conn's LocalAddr(). Set to the IP that remote
	// nodes will see, as that IP is what they'll use to validate our ID.
	PublicIP net.IP

	// Hook received queries. Return false if you don't want to propagate to the default handlers.
	OnQuery func(query *krpc.Msg, source snet.UDPAddr) (propagate bool)
	// Called when a peer successfully announces to us.
	OnAnnouncePeer func(infoHash metainfo.Hash, ip net.IP, port int, portOk bool)
	// How long to wait before resending queries that haven't received a response. Defaults to a
	// random value between 4.5 and 5.5s.
	QueryResendDelay func() time.Duration
	// TODO: Expose Peers, to return NodeInfo for received get_peers queries.
	PeerStore peer_store.Interface

	// If no Logger is provided, log.Default is used and log.Debug messages are filtered out. Note
	// that all messages without a log.Level, have log.Debug added to them before being passed to
	// this Logger.
	Logger log.Logger

	DefaultWant []krpc.Want
}

// ServerStats instance is returned by Server.Stats() and stores Server metrics
type ServerStats struct {
	// Count of nodes in the node table that responded to our last query or
	// haven't yet been queried.
	GoodNodes int
	// Count of nodes in the node table.
	Nodes int
	// Transactions awaiting a response.
	OutstandingTransactions int
	// Individual announce_peer requests that got a success response.
	SuccessfulOutboundAnnouncePeerQueries int64
	// Nodes that have been blocked.
	BadNodes                 uint
	OutboundQueriesAttempted int64
}

func jitterDuration(average time.Duration, plusMinus time.Duration) time.Duration {
	return average - plusMinus/2 + time.Duration(rand.Int63n(int64(plusMinus)))
}

type Peer = krpc.NodeAddr

func GlobalBootstrapAddrs(network string) (addrs []Addr, err error) {
	initDnsResolver()
	return
}

// Deprecated: Use function from krpc.
func RandomNodeID() (id krpc.ID) {
	return krpc.RandomNodeID()
}

func MakeDeterministicNodeID(public net.Addr) (id krpc.ID) {
	h := crypto.SHA1.New()
	h.Write([]byte(public.String()))
	h.Sum(id[:0:20])
	SecureNodeId(&id, missinggo.AddrIP(public))
	return
}
