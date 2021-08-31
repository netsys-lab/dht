package krpc

import (
	"bytes"
	"encoding/binary"
	"github.com/scionproto/scion/go/lib/snet"
	"net"
	"strconv"

	"github.com/anacrolix/torrent/bencode"
	"github.com/scionproto/scion/go/lib/addr"
)

type NodeAddr struct {
	IP   net.IP
	Port int
	IA   addr.IA
}

// A zero Port is taken to mean no port provided, per BEP 7.
func (me NodeAddr) String() string {
	if me.Port == 0 {
		return me.IA.String() + ",[" + me.IP.String() + "]"
	}
	return me.IA.String() + ",[" + me.IP.String() + "]:" + strconv.FormatInt(int64(me.Port), 10)
}

// [0-13 = IA][IPv4 or v6][last 2 = Port]
func (me *NodeAddr) UnmarshalBinary(b []byte) error {
	me.IA = addr.IAFromRaw(b[0:8])
	me.IP = make(net.IP, len(b)-2-8)
	copy(me.IP, b[8:len(b)-2])
	me.Port = int(binary.BigEndian.Uint16(b[len(b)-2:]))
	return nil
}

func (me *NodeAddr) UnmarshalBencode(b []byte) (err error) {
	var _b []byte
	err = bencode.Unmarshal(b, &_b)
	if err != nil {
		return
	}
	return me.UnmarshalBinary(_b)
}

func (me NodeAddr) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, me.IA.IAInt())
	b.Write(me.IP)
	binary.Write(&b, binary.BigEndian, uint16(me.Port))
	return b.Bytes(), nil
}

func (me NodeAddr) MarshalBencode() ([]byte, error) {
	return bencodeBytesResult(me.MarshalBinary())
}

func (me NodeAddr) UDP() snet.UDPAddr {
	addr, _ := snet.ParseUDPAddr(me.String())
	return *addr
}

func (me *NodeAddr) FromUDPAddr(ua *net.UDPAddr) {
	me.IP = ua.IP
	me.Port = ua.Port
}
