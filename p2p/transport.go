package p2p

import "net"

type Peer interface {
	Close() error
	RemoteAddr() net.Addr
	Send([]byte) error
}

// Transport is an interface that defines the methods for a transport layer in a peer-to-peer network.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
