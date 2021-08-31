package dht

import (
	"github.com/scionproto/scion/go/lib/snet"
	"net"
	"os"

	"github.com/anacrolix/dht/v2/int160"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/anacrolix/dht/v2/types"
	"github.com/anacrolix/missinggo/v2/iter"
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
)

func addrInOwnIsAs(addr string) (*snet.UDPAddr, error) {
	isia := os.Getenv("SCION_ISIA")
	return snet.ParseUDPAddr(isia + "," + addr)
}

func mustListen(addr string) *snet.Conn {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	ret, err := appnet.Listen(udpAddr)
	if err != nil {
		panic(err)
	}
	return ret
}

func addrResolver(addr string) func() ([]Addr, error) {
	add, err := net.ResolveUDPAddr("udp", addr)
	return func() ([]Addr, error) {
		return []Addr{NewAddr(snet.UDPAddr{IA: appnet.DefNetwork().IA, Host: add})}, err
	}
}

type addrMaybeId = types.AddrMaybeId

func randomIdInBucket(rootId int160.T, bucketIndex int) int160.T {
	id := int160.FromByteArray(krpc.RandomNodeID())
	for i := range iter.N(bucketIndex) {
		id.SetBit(i, rootId.GetBit(i))
	}
	id.SetBit(bucketIndex, !rootId.GetBit(bucketIndex))
	return id
}
