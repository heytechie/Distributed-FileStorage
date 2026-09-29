package p2p

import (
	"io"
	"net"
)

type Peer interface {
	io.ReadWriter
	Close() error
	RemoteAddr() net.Addr
	Send([]byte) error
	SendMessage([]byte) error
}

// Transport is an interface that defines the methods for a transport layer in a peer-to-peer network.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
