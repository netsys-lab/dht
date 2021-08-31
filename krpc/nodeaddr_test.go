package krpc

import (
	"github.com/scionproto/scion/go/lib/snet"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalScionNodeAddr(t *testing.T) {
	var na NodeAddr
	ia := []byte("\x00\x13\xff\xaa\x00\x01\x0e\x4b") // 19-ffaa:1:e4b
	ip := []byte("\x01\x02\x03\x04")                 // 1.2.3.4
	port := []byte("\x00\x50")                       // 80
	addr := append(append(ia, ip...), port...)
	require.NoError(t, na.UnmarshalBinary(addr))
	assert.EqualValues(t, "19-ffaa:1:e4b", na.IA.String())
	assert.EqualValues(t, "1.2.3.4", na.IP.String())
	assert.EqualValues(t, 80, na.Port)
}

func TestMarshalUnmarshalScionNodeAddr(t *testing.T) {
	addr, _ := snet.ParseUDPAddr("19-ffaa:1:e4b,[1.2.3.4]:80")
	originalNA := NodeAddr{IA: addr.IA, IP: addr.Host.IP, Port: addr.Host.Port}
	marshBytes, _ := originalNA.MarshalBinary()
	var newNA NodeAddr

	require.NoError(t, newNA.UnmarshalBinary(marshBytes))
	assert.EqualValues(t, "19-ffaa:1:e4b", newNA.IA.String())
	assert.EqualValues(t, "1.2.3.4", newNA.IP.String())
	assert.EqualValues(t, 80, newNA.Port)
}
