package dht

import (
	"github.com/scionproto/scion/go/lib/snet"
	"net"
	"strconv"

	"github.com/anacrolix/dht/v2/krpc"
	"github.com/scionproto/scion/go/lib/addr"
)

// Used internally to refer to node network addresses. String() is called a
// lot, and so can be optimized. Network() is not exposed, so that the
// interface does not satisfy net.Addr, as the underlying type must be passed
// to any OS-level function that take net.Addr.
type Addr interface {
	Raw() snet.UDPAddr
	Port() int
	IP() net.IP
	IA() addr.IA
	String() string
	KRPC() krpc.NodeAddr
}

// Speeds up some of the commonly called Addr methods.
type cachedAddr struct {
	raw  snet.UDPAddr
	port int
	ip   net.IP
	ia   addr.IA
	s    string
}

func (ca cachedAddr) String() string {
	return ca.s
}

func (ca cachedAddr) KRPC() krpc.NodeAddr {
	return krpc.NodeAddr{
		IP:   ca.ip,
		Port: ca.port,
		IA:   ca.ia,
	}
}

func (ca cachedAddr) IA() addr.IA {
	return ca.ia
}

func (ca cachedAddr) IP() net.IP {
	return ca.ip
}

func (ca cachedAddr) Port() int {
	return ca.port
}

func (ca cachedAddr) Raw() snet.UDPAddr {
	return ca.raw
}

func NewAddr(raw snet.UDPAddr) Addr {
	ip := raw.Host.IP
	port := raw.Host.Port
	ia := raw.IA
	str := ia.String() + ",[" + ip.String() + "]:" + strconv.FormatInt(int64(port), 10)
	return cachedAddr{
		raw:  raw,
		s:    str,
		ip:   ip,
		port: port,
		ia:   ia,
	}
}
