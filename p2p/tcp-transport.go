package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	//basically if we dial and retrive then its a outbound connection, if we accept then its inbound connection
	outbound bool
}

func NewPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOption struct {
	ListenAddress string
	HandshakeFunc HandShakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOption
	listner net.Listener
	mu      sync.RWMutex
	peers   map[string]Peer
}

func NewTCPTransport(opts TCPTransportOption) *TCPTransport {
	return &TCPTransport{
		TCPTransportOption: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {

		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}
		fmt.Printf("Incomming new connection:%+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake error: %s\n", err)
		conn.Close()
		return
	}
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("Error decoding message: %s\n", err)
			continue
		}
		fmt.Printf("Received message: %+v\n", msg)
	}
}
