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

type TCPTransport struct {
	listenAddress string
	listner       net.Listener
	handshakeFunc HandShakeFunc
	decoder       Decoder
	mu            sync.RWMutex
	peers         map[string]Peer
}

func NewPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
func NewTCPTransport(listenAddress string) *TCPTransport {
	return &TCPTransport{
		handshakeFunc: NOPHandshakeFunc,
		listenAddress: listenAddress,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.listenAddress)
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

	if err := t.handshakeFunc(peer); err != nil {
		fmt.Printf("Handshake error: %s\n", err)
		conn.Close()
		return
	}
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("Error decoding message: %s\n", err)
			continue
		}
	}
}
