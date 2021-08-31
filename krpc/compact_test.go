package krpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalSlice(t *testing.T) {
	var data CompactIPv4NodeInfo
	ipAndPort1 := "\x01\x02\x03\x04\x00\x06"
	ipAndPort2 := "\x02\x03\x04\x05\x00\x07"
	ia := "\x00\x13\xff\xaa\x00\x01\x0e\x4b" // 19-ffaa:1:e4b
	nodeId := "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	err := data.UnmarshalBencode([]byte("68:" + nodeId + ia + ipAndPort1 + nodeId + ia + ipAndPort2))
	require.NoError(t, err)
	require.Len(t, data, 2)

	assert.Equal(t, "1.2.3.4", data[0].Addr.IP.String())
	assert.Equal(t, "2.3.4.5", data[1].Addr.IP.String())

	assert.Equal(t, 6, data[0].Addr.Port)
	assert.Equal(t, 7, data[1].Addr.Port)

	assert.Equal(t, 19, int(data[0].Addr.IA.I))
	assert.Equal(t, 281105609592395, int(data[0].Addr.IA.A))

	assert.Equal(t, 19, int(data[1].Addr.IA.I))
	assert.Equal(t, 281105609592395, int(data[1].Addr.IA.A))
}
