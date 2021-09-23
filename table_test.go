package dht

import (
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
	"net"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/netsys-lab/dht/int160"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tbl := table{k: 8}
	var maxFar int160.T
	maxFar.SetMax()
	assert.Equal(t, 0, tbl.bucketIndex(maxFar))
	assert.Panics(t, func() { tbl.bucketIndex(tbl.rootID) })

	assert.Error(t, tbl.addNode(&node{}))
	assert.Equal(t, 0, tbl.buckets[0].Len())

	id0 := int160.FromByteString("\x2f\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	id1 := int160.FromByteString("\x2e\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	addr0 := snet.UDPAddr{IA: appnet.DefNetwork().IA, Host: &net.UDPAddr{}}
	addr1 := snet.UDPAddr{IA: appnet.DefNetwork().IA, Host: &net.UDPAddr{}}
	n0 := &node{nodeKey: nodeKey{
		Id:   id0,
		Addr: NewAddr(addr0),
	}}
	n1 := &node{nodeKey: nodeKey{
		Id:   id1,
		Addr: NewAddr(addr1),
	}}

	assert.NoError(t, tbl.addNode(n0))
	assert.Equal(t, 1, tbl.buckets[2].Len())

	assert.Error(t, tbl.addNode(n0))
	assert.Equal(t, 1, tbl.buckets[2].Len())
	assert.Equal(t, 1, tbl.numNodes())

	assert.NoError(t, tbl.addNode(n1))
	assert.Equal(t, 2, tbl.buckets[2].Len())
	assert.Equal(t, 2, tbl.numNodes())

	tbl.dropNode(n0)
	assert.Equal(t, 1, tbl.buckets[2].Len())
	assert.Equal(t, 1, tbl.numNodes())

	tbl.dropNode(n1)
	assert.Equal(t, 0, tbl.buckets[2].Len())
	assert.Equal(t, 0, tbl.numNodes())
}

func TestRandomIdInBucket(t *testing.T) {
	tbl := table{
		rootID: int160.FromByteArray(RandomNodeID()),
	}
	t.Logf("%v: table root id", tbl.rootID)
	for i := range tbl.buckets {
		id := tbl.randomIdForBucket(i)
		t.Logf("%v: random id for bucket index %v", id, i)
		qt.Assert(t, tbl.bucketIndex(id), qt.Equals, i)
	}
}
