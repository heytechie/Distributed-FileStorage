package p2p

import "net"

const (
	IncomingMessage byte = 0x1
	IncomingStream  byte = 0x2
)

type RPC struct {
	From    net.Addr
	Payload []byte
	Stream  bool	
}
