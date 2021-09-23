package peer_store

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/netsys-lab/dht/v2/krpc"
)

type InfoHash = metainfo.Hash

type Interface interface {
	AddPeer(InfoHash, krpc.NodeAddr)
	GetPeers(InfoHash) []krpc.NodeAddr
}
