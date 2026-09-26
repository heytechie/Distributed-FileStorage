package p2p

type Peer interface {
	Close() error
}

// Transport is an interface that defines the methods for a transport layer in a peer-to-peer network.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
