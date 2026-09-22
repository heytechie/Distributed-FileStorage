package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	addr := ":4000"
	tr := NewTCPTransport(addr)

	assert.Equal(t, tr.listenAddress, addr)

	//server
	// tr.listner.Accept()

	assert.Nil(t, tr.ListenAndAccept())
}
