package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOption{
		ListenAddress: "localhost:3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       &DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddress, "localhost:3000")

	//server
	// tr.listner.Accept()

	assert.Nil(t, tr.ListenAndAccept())
}
